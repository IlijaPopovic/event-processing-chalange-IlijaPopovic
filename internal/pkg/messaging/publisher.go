package messaging

import "github.com/Bitstarz-eng/event-processing-challenge/internal/casino"

type Publisher interface {
    Publish(event casino.Event) error
    Close() error
}