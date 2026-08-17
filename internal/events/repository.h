#include <sqlite3.h>
#include <stdint.h>

int insert_event(sqlite3 *db, char *eventId, char *recipientIds, int64_t quantity);
