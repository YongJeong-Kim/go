package main

import (
	"encoding/json"
	"github.com/hibiken/asynq"
	"log"
	"time"
)

type EmailTaskPayload struct {
	// ID for the email recipient.
	UserID int
}

func main() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: "localhost:16379"})

	// Create a task with typename and payload.
	payload, err := json.Marshal(EmailTaskPayload{UserID: 42})
	if err != nil {
		log.Fatal(err)
	}
	t1 := asynq.NewTask("email:welcome", payload)

	t2 := asynq.NewTask("email:reminder", payload)

	// Process the task immediately.
	info, err := client.Enqueue(t1)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(" [*] Successfully enqueued task: %+v", info)

	// Process the task 24 hours later.
	info, err = client.Enqueue(t2, asynq.ProcessIn(24*time.Hour))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(" [*] Successfully enqueued task: %+v", info)
}
