package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/pkg/formatter"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/pkg/messaging"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/pkg/metrics"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/pkg/repository"
)

type Service struct {
	consumer      messaging.Consumer
	repo          *repository.EventRepository
	metrics       *metrics.Collector
}

func NewService(
	consumer messaging.Consumer,
	repo *repository.EventRepository,
	metrics       *metrics.Collector,
) *Service {
	return &Service{
		consumer:     consumer,
		repo:         repo,
		metrics: metrics,
	}
}

func (s *Service) ProcessEvents(ctx context.Context) {
	eventCh, err := s.consumer.Consume(ctx)
	if err != nil {
		log.Fatal("Failed to start consuming events:", err)
	}

	for event := range eventCh {
		if err := s.processEvent(ctx, event); err != nil {
			log.Printf("Error processing event %d: %v", event.ID, err)
		}
	}
}

func (s *Service) processEvent(ctx context.Context, event casino.Event) error {

	event.Description = formatter.FormatEventDescription(event)

	s.metrics.HandleEvent(event)

	if err := logEventJSON(event); err != nil {
		return fmt.Errorf("failed to log event: %v", err)
	}

	if err := s.repo.SaveEvent(ctx, &event); err != nil {
		return fmt.Errorf("failed to save event: %w", err)
	}

	return nil
}

func logEventJSON(event casino.Event) error {
	jsonData, err := json.Marshal(event)
	if err != nil{
		return fmt.Errorf("error marshaling event to JSON: %w", err)
	}

	log.Println(string(jsonData))

	return nil
}
