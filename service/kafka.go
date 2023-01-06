package service

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
	"github.com/segmentio/kafka-go/sasl/plain"
)

func getProducer(brokers []string, topic string, username string, password string) *kafka.Writer {
	mechanism := plain.Mechanism{
		Username: username,
		Password: password,
	}
	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}
	return kafka.NewWriter(kafka.WriterConfig{
		Dialer:  dialer,
		Brokers: brokers,
		Topic:   topic,
	})
}

type Producer struct {
	Writer  *kafka.Writer
	Timeout time.Duration
}

func (pdr *Producer) Produce(mesg ...kafka.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), pdr.Timeout)
	defer cancel()
	return pdr.Writer.WriteMessages(ctx, mesg...)
}

func NewProducer(brokers []string, topic string, username string, password string, timeout time.Duration) *Producer {
	return &Producer{
		Writer:  getProducer(brokers, topic, username, password),
		Timeout: timeout,
	}
}

func MapToHeaders(mapInput map[string]any) []protocol.Header {
	val := make([]protocol.Header, 0, len(mapInput))
	for k, v := range mapInput {
		val = append(val, protocol.Header{Key: k, Value: []byte(fmt.Sprint(v))})
	}
	return val
}
