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
		log.Fatal("dial failed:", err)
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := conn.ReadBatch(10e3, 1e6)

	b := make([]byte, 10e3)
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		log.Println(string(b[:n]))
	}

	if err := batch.Close(); err != nil {
		log.Fatal("batch close failed:", err)
	}
	if err := conn.Close(); err != nil {
		log.Fatal("conn close  failed", err)
	}
}
