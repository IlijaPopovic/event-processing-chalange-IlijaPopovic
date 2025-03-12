package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/pkg/messaging"
	"github.com/IBM/sarama"
)

var _ messaging.Consumer = (*Consumer)(nil)

type Consumer struct {
	consumer sarama.Consumer
	topic    string
}

func NewConsumer(brokerList []string, topic string) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{consumer: consumer, topic: topic}, nil
}

func (c *Consumer) Consume(ctx context.Context) (<-chan casino.Event, error) {
	eventCh := make(chan casino.Event)
	
	partitionConsumer, err := c.consumer.ConsumePartition(c.topic, 0, sarama.OffsetOldest)
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(eventCh)
		defer partitionConsumer.Close()

		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping consumer...")
				return
			case msg := <-partitionConsumer.Messages():
				var event casino.Event
				if err := json.Unmarshal(msg.Value, &event); err != nil {
					log.Printf("Error decoding message: %v", err)
					continue
				}
				eventCh <- event
			case err := <-partitionConsumer.Errors():
				log.Printf("Error consuming message: %v", err)
			}
		}
	}()

	return eventCh, nil
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}