package main

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"gounread/repository"
	"io"
	"log"
	"net/http"
)

func eventJoinRoom(nc *nats.Conn, roomID, userID string) {
	joinRoomSubject := fmt.Sprintf("join.room.%s", roomID)
	nc.Subscribe(joinRoomSubject, func(msg *nats.Msg) {
		messages := requestJoinRoom(roomID, userID)
		for _, m := range messages {
			if userID == m.Sender {
				log.Println("-------------------------------------------------------")
				log.Printf("%s, %s\n", m.RoomID, m.Sent)
				log.Printf("%s <- %s, %d\n", m.Msg, m.Sender, len(m.Unread))
				log.Println("-------------------------------------------------------")
			} else {
				log.Println("-------------------------------------------------------")
				log.Printf("%s, %s\n", m.RoomID, m.Sent)
				log.Printf("%d, %s -> %s\n", len(m.Unread), m.Sender, m.Msg)
				log.Println("-------------------------------------------------------")
			}

		}
	})
}

func requestJoinRoom(roomID, userID string) []*repository.GetMessagesByRoomIDAndTimeResult {
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/rooms/"+roomID+"/join", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("user", userID)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result []*repository.GetMessagesByRoomIDAndTimeResult
	json.Unmarshal(data, &result)
	log.Println(result)

	return result
}
