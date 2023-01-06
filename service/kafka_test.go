package service_test

import (
	"testing"
	"time"

	"github.com/lwinmgmg/outbox/service"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
)

func TestNewProducer(t *testing.T) {
	producer := service.NewProducer([]string{"172.30.1.137:9092"}, "erp.address.1", "admin", "admin-secret", time.Millisecond*15000)
	mapData := map[string]any{
		"Content-Type": "application/json",
	}
	if err := producer.Produce(kafka.Message{Key: []byte("Key"), Value: []byte("Value"), Headers: service.MapToHeaders(mapData)}); err != nil {
		t.Errorf("Error on producing : %v\n", err)
	}
}

func TestMapToHeaders(t *testing.T) {
	mapData := map[string]any{
		"Content-Type": "application/json",
	}
	headerList := service.MapToHeaders(mapData)
	if headerList[0].Key != "Content-Type" || string(headerList[0].Value) != "application/json" {
		t.Errorf("Expecting %v but getting %v\n", protocol.Header{
			Key:   "Content-Type",
			Value: []byte("application/json"),
		},
			service.MapToHeaders(mapData)[0])
	}
}
