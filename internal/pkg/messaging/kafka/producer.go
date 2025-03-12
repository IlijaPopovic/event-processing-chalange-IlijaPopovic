package kafka

import (
	"encoding/json"
	"log"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
	"github.com/IBM/sarama"
)

type Producer struct {
    producer sarama.SyncProducer
    topic    string
}

func NewProducer(brokerList []string, topic string) (*Producer, error) {
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true

    producer, err := sarama.NewSyncProducer(brokerList, config)
    if err != nil {
        return nil, err
    }

    return &Producer{producer: producer, topic: topic}, nil
}

func (p *Producer) Publish(event casino.Event) error {
    eventData, err := json.Marshal(event)
    if err != nil {
        return err
    }

    msg := &sarama.ProducerMessage{
        Topic: p.topic,
        Value: sarama.StringEncoder(eventData),
    }

    _, _, err = p.producer.SendMessage(msg)
    if err != nil {
        return err
    }

    log.Printf("Published event to Kafka: ID %d, Type %s", event.ID, event.Type)
    return nil
}

func (p *Producer) Close() error {
    return p.producer.Close()
}