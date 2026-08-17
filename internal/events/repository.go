package events

//#include <repository.h>
import "C"

import (
	sqlite "github.com/tacohirosystems/tacohiro/internal/database"
	"github.com/tacohirosystems/tacohiro/internal/discord"
)

type (
	Repository struct {
		DiscordDBs map[discord.UserID]*sqlite.DB
	}

	CreateEvent struct {
		SenderID     discord.UserID
		RecipientIDs []discord.UserID
		Quantity     int64
		EventID      string
		Source       string
	}
)

func (r *Repository) insertEvent(c CreateEvent) error {
	db := r.DiscordDBs[c.SenderID]
	if db == nil {
		// TODO: Initialize
	}

	db.Write(func(conn *sqlite.WriteConn) error {
		if conn == nil {
			// FIXME
			return nil
		}

		return nil
	})

	return nil
}
