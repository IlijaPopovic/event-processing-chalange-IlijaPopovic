package repository

import (
	"context"
	"database/sql"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) SaveEvent(ctx context.Context, event *casino.Event) error {
	query := `
	INSERT INTO events (
		player_id, game_id, type, amount, currency, 
		has_won, created_at, amount_eur, description
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.ExecContext(ctx, query,
		event.PlayerID,
		event.GameID,
		event.Type,
		event.Amount,
		event.Currency,
		event.HasWon,
		event.CreatedAt,
		event.AmountEUR,
		event.Description,
	)
	return err
}