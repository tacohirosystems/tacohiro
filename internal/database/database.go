package sqlite

/*
#cgo LDFLAGS: -lsqlite3
#include <database.h>
*/
import "C"
import (
	"fmt"
	"log/slog"
)

type (
	DB struct {
		Path string

		ReadPool chan struct{}
		WritePool chan struct{}
	}

	ReadConn struct {
		conn *C.sqlite3
	}

	WriteConn struct {
		conn *C.sqlite3
	}
)

var (
	ErrDatabaseReleaseConn = fmt.Errorf("failed to release database connection")
)

func (db *DB) releaseWrite(conn *WriteConn) error {
	var err error
	if rc := C.sqlite3_close_v2(conn.conn); rc != C.SQLITE_OK {
		err = fmt.Errorf("%w: %d", ErrDatabaseReleaseConn, rc)
	}
	db.WritePool <- struct{}{}
	return err
}

func (db *DB) releaseRead(conn *ReadConn) error {
	var err error
	if rc := C.sqlite3_close_v2(conn.conn); rc != C.SQLITE_OK {
		err = fmt.Errorf("%w: %d", ErrDatabaseReleaseConn, rc)
	}
	db.ReadPool <- struct{}{}
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
	conn := &ReadConn{conn: nil}
	if rc := C.sqlite3_open_v2(C.CString(db.Path), &conn.conn, C.int(C.SQLITE_OPEN_READONLY), nil); rc != C.SQLITE_OK {
		slog.Info(fmt.Sprintf("failed to acquire a read connection: %d", rc));
		return nil, fmt.Errorf("TODO: IT EXPLODED!")
	}
	return conn, nil
}

func (db *DB) getWriteConn() (*WriteConn, error) {
	<- db.WritePool
	conn := &WriteConn{conn: nil}

	if rc := C.sqlite3_open_v2(C.CString(db.Path), &conn.conn, C.int(C.SQLITE_OPEN_READWRITE), nil); rc != C.SQLITE_OK {
		slog.Info(fmt.Sprintf("failed to acquire a write connection: %d", rc));
		return nil, fmt.Errorf("TODO: IT EXPLODED!")
	}
	return conn, nil
}

func (db *DB) SetMaxReadConnections(max int64) {
	for range max {
		db.ReadPool <- struct{}{}
	}
}

func (db *DB) SetMaxWriteConnections(max int64) {
	for range max {
		db.WritePool <- struct{}{}
	}
}

func (db *DB) SetPragmas() error {
	return db.Write(func (conn *WriteConn) error {
		if rc := C.init(conn.conn); rc != C.SQLITE_OK {
			slog.Info(fmt.Sprintf("failed to initialize the database: %d", rc));
			return fmt.Errorf("failed to initialize the database: %d", rc)
		}
		return nil
	})
}
