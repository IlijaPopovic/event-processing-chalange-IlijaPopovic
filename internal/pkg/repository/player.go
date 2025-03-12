package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
)

type PlayerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

func (r *PlayerRepository) GetPlayerByID(ctx context.Context, id int) (casino.Player, error) {
	var player casino.Player
	query := `SELECT id, email, last_signed_in_at FROM players WHERE id = $1`
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&player.ID,
		&player.Email,
		&player.LastSignedInAt,
	)
	
	if errors.Is(err, sql.ErrNoRows) {
		return casino.Player{}, fmt.Errorf("player not found")
	}
	if err != nil {
		return casino.Player{}, err
	}
	return player, nil
}