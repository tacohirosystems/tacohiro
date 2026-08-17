package events

import (
	"log/slog"

	"github.com/tacohirosystems/tacohiro/internal/discord"
)

type (
	Event struct {
		SenderID     discord.UserID
		RecipientIDs []discord.UserID
		Quantity     int64
	}
)

type (
	Service struct {
		Repository *Repository
		Logger     *slog.Logger
	}
)

func (s *Service) CreateEvent(event Event) error {
	s.Logger.Debug("Creating event...")
	return nil
}
