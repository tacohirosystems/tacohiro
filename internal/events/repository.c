#include <stdio.h>
#include <sqlite3.h>
#include <string.h>
#include <stdint.h>

typedef struct {

} CreateEvent;

int insert_event(sqlite3 *db, char *eventId, char *recipientIds, int64_t quantity) {
  int rc;

  sqlite3_stmt* stmt = NULL;
  char* query =
    "INSERT INTO events(event_id, recipient_user_ids, quantity, source) VALUES (?, ?, ?, 'discord');";

  rc = sqlite3_prepare_v2(db, query, strlen(query) + 1, &stmt, NULL);
  if (rc != SQLITE_OK) {
    printf("Failed to initialize the database: %d\n", rc);
    return rc;
  }
  return 0;
}
