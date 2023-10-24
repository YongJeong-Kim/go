package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

func main() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:29092", "localhost:9093", "localhost:9094"},
		Topic:     "ggg",
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})
	r.SetOffset(-1)

	for {
		m, err := r.FetchMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
