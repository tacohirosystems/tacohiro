#include <stdio.h>
#include <sqlite3.h>

int init(sqlite3 *db) {
  char *zErrMsg = 0;
  int rc;

  char *initStmt =
    "PRAGMA foreign_keys = ON;"
    "PRAGMA busy_timeout = 5000;"
    "PRAGMA journal_mode = WAL;"
    "PRAGMA synchronous = NORMAL;"
    "PRAGMA cache_size = 1000000000;";

  rc = sqlite3_exec(db, initStmt, NULL, 0, &zErrMsg);
  if (rc != SQLITE_OK) {
    printf("Failed to initialize the database: %d\n", rc);
    sqlite3_free(zErrMsg);
    return rc;
  }
  // printf("Initialized SQLite with:\n %s\n", initStmt);
  // printf("Completed query execution without any issues\n");
  return 0;
}
