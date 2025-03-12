package generator

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/pkg/currency"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/pkg/messaging"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/pkg/repository"
)

type Service struct {
	publisher    messaging.Publisher
	playerRepo   *repository.PlayerRepository
	converter    *currency.Converter
}

func NewService(
	publisher messaging.Publisher,
	playerRepo *repository.PlayerRepository,
	converter *currency.Converter,
) *Service {
	return &Service{
		publisher:    publisher,
		playerRepo:   playerRepo,
		converter:    converter,
	}
}

func (s *Service) GenerateEvents(ctx context.Context) <-chan casino.Event {
	eventCh := make(chan casino.Event)
	var id int

	go func() {
		defer close(eventCh)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				event := s.generateEvent(ctx, id)
				eventCh <- event
				
				if err := s.publisher.Publish(event); err != nil {
					log.Printf("Failed to publish event: %v", err)
				}
				
				id++
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			}
		}
	}()

	return eventCh
}

func (s *Service) generateEvent(ctx context.Context, id int) casino.Event {
	amount, currency := randomAmountCurrency()
	
	event := casino.Event{
		ID:        id,
		PlayerID:  10 + rand.Intn(10),
		GameID:    100 + rand.Intn(10),
		Type:      randomType(),
		Amount:    amount,
		Currency:  currency,
		HasWon:    randomHasWon(),
		CreatedAt: time.Now(),
	}

    amountEUR, err := s.converter.ConvertToEUR(ctx, event.Amount, event.Currency)
    if err != nil {
        log.Printf("Currency conversion failed for event %d: %v", id, err)
    } else {
        event.AmountEUR = amountEUR
    }

	player, err := s.playerRepo.GetPlayerByID(ctx, event.PlayerID)
	if err != nil {
		log.Printf("Error fetching player %d: %v", event.PlayerID, err)
	} else {
		event.Player = player
	}

	event.Description = fmt.Sprintf("Generated event of type %s", event.Type)

	return event
}

func randomType() string {
	return casino.EventTypes[rand.Intn(len(casino.EventTypes))]
}

func randomAmountCurrency() (amount int, currency string) {
	currency = casino.Currencies[rand.Intn(len(casino.Currencies))]

	switch currency {
	case "BTC":
		amount = rand.Intn(1e5)
	default:
		amount = rand.Intn(2000)
	}

	return
}

func randomHasWon() bool {
	return rand.Intn(100) < 50
}