package database

import _ "embed"

/*
#cgo LDFLAGS: -lsqlite3
#include <stdlib.h>
#include <sqlite3.h>
*/
import "C"
import (
	"fmt"
	"log/slog"
	"unsafe"
)

type (
	DB struct {
		Path string

		ReadPool  chan struct{}
		WritePool chan struct{}
		Logger    *slog.Logger
	}

	ReadConn struct {
		conn *C.sqlite3
	}

	WriteConn struct {
		conn *C.sqlite3
	}
)

var (
	ErrDatabaseReleaseConn               = fmt.Errorf("failed to release database connection")
	ErrDatabaseFailedToInitializePragmas = fmt.Errorf("failed to initialize the database pragmas")
	ErrDatabaseFailedToAcquireReadConn   = fmt.Errorf("failed to acquire read connection")
	ErrDatabaseFailedToAcquireWriteConn  = fmt.Errorf("failed to acquire write connection")

	//go:embed sql/init.sql
	initQuery string
)

func (db *DB) releaseWrite(conn *WriteConn) error {
	db.Logger.Debug("Releasing write connection...")
	var err error
	if rc := C.sqlite3_close_v2(conn.conn); rc != C.SQLITE_OK {
		db.Logger.Error(ErrDatabaseReleaseConn.Error(), "sqlite_err_code", rc)
		err = ErrDatabaseReleaseConn
	}
	db.WritePool <- struct{}{}
	db.Logger.Debug("Returned to conn write pool")
	return err
}

func (db *DB) releaseRead(conn *ReadConn) error {
	db.Logger.Debug("Releasing read connection...")
	var err error
	if rc := C.sqlite3_close_v2(conn.conn); rc != C.SQLITE_OK {
		db.Logger.Error(ErrDatabaseReleaseConn.Error(), "sqlite_err_code", rc)
		err = ErrDatabaseReleaseConn
	}
	db.ReadPool <- struct{}{}
	db.Logger.Debug("Returned to conn read pool")
	return err
}

func (db *DB) Write(execQuery func(conn *WriteConn) error) error {
	conn, err := db.getWriteConn()
	defer db.releaseWrite(conn)
	if err != nil {
		return err
	}

	if err := execQuery(conn); err != nil {
		return err
	}
	return nil
}

func (db *DB) Read(execQuery func(conn *ReadConn) error) error {
	conn, err := db.getReadConn()
	defer db.releaseRead(conn)
	if err != nil {
		return err
	}

	if err := execQuery(conn); err != nil {
		return err
	}
	return nil
}

func (db *DB) getReadConn() (*ReadConn, error) {
	db.Logger.Debug("Acquiring read connection...")
	conn := &ReadConn{conn: nil}
	if rc := C.sqlite3_open_v2(C.CString(db.Path), &conn.conn, C.int(C.SQLITE_OPEN_READONLY), nil); rc != C.SQLITE_OK {
		db.Logger.Info("failed to acquire a read connection", "sqlite_err_code", rc)
		return nil, ErrDatabaseFailedToAcquireReadConn

	}
	return conn, nil
}

func (db *DB) getWriteConn() (*WriteConn, error) {
	db.Logger.Debug("Acquiring write connection...")
	<-db.WritePool
	conn := &WriteConn{conn: nil}

	if rc := C.sqlite3_open_v2(C.CString(db.Path), &conn.conn, C.int(C.SQLITE_OPEN_READWRITE), nil); rc != C.SQLITE_OK {
		db.Logger.Error("failed to acquire a write connection", "sqlite_err_code", rc)
		return nil, ErrDatabaseFailedToAcquireWriteConn
	}
	return conn, nil
}

func (db *DB) SetMaxReadConnections(max int64) {
	for range max {
		db.ReadPool <- struct{}{}
	}
}

func (db *DB) SetMaxWriteConnections(max int64) {
	if max > 1 {
		db.Logger.Warn("Max write connection set to %d, more than 1")
	}

	for range max {
		db.WritePool <- struct{}{}
	}
}

func (db *DB) SetPragmas() error {
	return db.Write(func(conn *WriteConn) error {
		query := C.CString(initQuery)
		errMsg := C.CString("")
		testInt := new(C.int(0))
		defer C.free(unsafe.Pointer(query))
		defer C.free(unsafe.Pointer(errMsg))

		db.Logger.Debug("Initializing pragmas...")

		if rc := C.sqlite3_exec(conn.conn, query, nil, unsafe.Pointer(testInt), &errMsg); rc != C.SQLITE_OK {
			db.Logger.Error("failed to initialize the database", "sqlite_err_code", rc)
			return ErrDatabaseFailedToInitializePragmas
		}
		return nil
	})
}
