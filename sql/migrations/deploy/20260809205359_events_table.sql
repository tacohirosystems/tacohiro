-- Deploy tacohiro:20260809205359_events_table to sqlite

BEGIN IMMEDIATE;

-- Represents what taco gifts a user sent in Discord
CREATE TABLE events (
  event_id INTEGER NOT NULL,
  -- This is JSON as TEXT since we're dealing with multiple user IDs potentially.
  recipient_user_ids TEXT NOT NULL,
  quantity TEXT NOT NULL,
  source TEXT NOT NULL,
  created_at TEXT DEFAULT current_timestamp NOT NULL,
  updated_at TEXT,
  PRIMARY KEY (event_id)
) STRICT;

COMMIT;
