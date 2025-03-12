package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/pkg/messaging/kafka"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/pkg/metrics"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/pkg/processor"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/pkg/repository"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Initializing database
	db, err := sql.Open("postgres", "postgres://casino:casino@database:5432/casino?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initializing dependencies
	repo := repository.NewEventRepository(db)

	// Initializing Kafka consumer
	brokers := []string{"kafka:9092"}
	topic := "casino-events"
	consumer, err := kafka.NewConsumer(brokers, topic)
	if err != nil {
		log.Fatal("Failed to initialize Kafka consumer:", err)
	}
	defer consumer.Close()

	// Initializing metrics
	metricsStorage := metrics.NewStorage()
	metricsCollector := metrics.NewCollector(metricsStorage)
	metricsServer := metrics.NewServer(metricsStorage)

	// Creatting processor service with metrics
	procService := processor.NewService(consumer, repo, metricsCollector)

	// Starting HTTP server
	router := mux.NewRouter()
	if router == nil {
		log.Fatal("Failed to create router")
	}

	// Registering routes
	metricsServer.RegisterRoutes(router)

	// Starting server with better error handling
	go func() {
		log.Println("Starting metrics server on :8080")
		if err := http.ListenAndServe(":8080", router); err != nil {
			log.Fatalf("Metrics server failed: %v", err)
		}
	}()

	// Setup context with graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("Starting event processing...")
	procService.ProcessEvents(ctx)
	log.Println("Event processing stopped")
}