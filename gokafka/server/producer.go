package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

func main() {
	topic := "ttt"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:29092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Fatal("set deadline failed:", err)
	}
	_, err = conn.WriteMessages(
		kafka.Message{
			Key:   []byte("key1"),
			Value: []byte("one!"),
		},
		kafka.Message{
			Key:   []byte("key2"),
			Value: []byte("two!"),
		},
		kafka.Message{
			Key:   []byte("key3"),
			Value: []byte("three!"),
		},
	)
	if err != nil {
		log.Fatal("write messages failed:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("close failed:", err)
	}
}
