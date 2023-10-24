package main

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
)

func panicIfErr(err error) {
	if err != nil {
		fmt.Println("err :", err)
	}
}

type User struct {
	Name string
	Age  int
}

func main() {
	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true
	connectionString := []string{
		"localhost:29092",
	}

	conn, err := sarama.NewClient(connectionString, conf)

	producer, err := sarama.NewSyncProducerFromClient(conn)
	panicIfErr(err)
	user := User{"Gopher", 7}
	jsonBody, err := json.Marshal(user)

	msg := &sarama.ProducerMessage{
		Topic:     "ttt",
		Value:     sarama.ByteEncoder(jsonBody),
		Partition: 3,
	}
	_, _, err = producer.SendMessage(msg)
}
