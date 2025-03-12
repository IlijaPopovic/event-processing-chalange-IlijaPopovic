package formatter

import (
	"fmt"
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
)

func FormatEventDescription(event casino.Event) string {
	switch event.Type {
	case "game_start":
		return formatGameStart(event)
	case "bet":
		return formatBet(event)
	case "deposit":
		return formatDeposit(event)
	case "game_stop":
		return formatGameStop(event)
	default:
		return fmt.Sprintf("Unknown event type: %s", event.Type)
	}
}

func formatGameStart(event casino.Event) string {
	game := casino.Games[event.GameID]
	return fmt.Sprintf(
		"Player #%d started playing a game \"%s\" on %s.",
		event.PlayerID,
		game.Title,
		formatDateTime(event.CreatedAt),
	)
}

func formatBet(event casino.Event) string {
	game := casino.Games[event.GameID]
	playerEmail := ""
	if event.Player.Email != "" {
		playerEmail = fmt.Sprintf(" (%s)", event.Player.Email)
	}

	return fmt.Sprintf(
		"Player #%d%s placed a bet of %d %s (%.2f EUR) on a game \"%s\" on %s.",
		event.PlayerID,
		playerEmail,
		event.Amount,
		event.Currency,
		float64(event.AmountEUR)/100, // Convert cents to EUR
		game.Title,
		formatDateTime(event.CreatedAt),
	)
}

func formatDeposit(event casino.Event) string {
	return fmt.Sprintf(
		"Player #%d made a deposit of %d %s on %s.",
		event.PlayerID,
		event.Amount,
		event.Currency,
		formatDateTime(event.CreatedAt),
	)
}

func formatGameStop(event casino.Event) string {
	game := casino.Games[event.GameID]
	return fmt.Sprintf(
		"Player #%d stopped playing a game \"%s\" on %s.",
		event.PlayerID,
		game.Title,
		formatDateTime(event.CreatedAt),
	)
}

func formatDateTime(t time.Time) string {
	return t.Format("January 2, 2006 at 15:04 UTC")
}