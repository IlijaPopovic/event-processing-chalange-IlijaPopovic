package messaging

import (
	"context"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
)

type Consumer interface {
	Consume(ctx context.Context) (<-chan casino.Event, error)
	Close() error
}