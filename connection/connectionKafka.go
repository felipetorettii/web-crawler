package connection

import (
	"fmt"

	"github.com/segmentio/kafka-go"
)

// ConnectionKafka create connection with kafka
func ConnectionKafka() *kafka.Reader {
	config := kafka.ReaderConfig{
		Brokers:  []string{"kafka:29092"},
		GroupID:  "consumer-crawling",
		Topic:    "topic-crawling",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	}

	fmt.Printf("Listening topic: %s...", config.Topic)

	return kafka.NewReader(config)
}
