package metrics

import (
	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
)

type Collector struct {
	storage *Storage
}

func NewCollector(storage *Storage) *Collector {
	return &Collector{storage: storage}
}

func (c *Collector) HandleEvent(event casino.Event) {
	c.storage.RecordEvent(event)
}