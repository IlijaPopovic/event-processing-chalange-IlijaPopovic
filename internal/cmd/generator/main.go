package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/pkg/currency"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/pkg/generator"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/pkg/messaging/kafka"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/pkg/repository"
	_ "github.com/lib/pq"
)

func main() {
	// Getting API key from environment
	apiKey := os.Getenv("EXCHANGE_RATE_API_KEY")
	if apiKey == "" {
		log.Fatal("EXCHANGE_RATE_API_KEY environment variable not set")
	}

	// Initializing database connection
	db, err := sql.Open("postgres", "postgres://casino:casino@database:5432/casino?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initializing dependencies
	playerRepo := repository.NewPlayerRepository(db)
	conv := currency.NewConverter(apiKey)

	// Initializing Kafka producer
	brokers := []string{"kafka:9092"}
	topic := "casino-events"
	producer, err := kafka.NewProducer(brokers, topic)
	if err != nil {
		log.Fatal("Failed to initialize Kafka producer:", err)
	}
	defer producer.Close()

	// Creating generator service
	genService := generator.NewService(producer, playerRepo, conv)

	// Setup context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Starting generating events
	events := genService.GenerateEvents(ctx)

	// Processing generated events
	for event := range events {
		log.Printf("Generated event: %+v", event)
	}

	log.Println("Event generation completed")
}