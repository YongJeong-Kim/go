package main

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"gounread/api"
	"gounread/repository"
	"io"
	"log"
	"net/http"
)

func requestReadMessage(r *repository.GetRoomsByUserIDResult, userID string) error {
	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/rooms/"+r.RoomID+"/read", nil)
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("read messag error. ")
	}
	return nil
}

func getConnectUserRooms(userID string) []*repository.GetRoomsByUserIDResult {
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/connect", nil)
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

	var result []*repository.GetRoomsByUserIDResult
	json.Unmarshal(data, &result)
	log.Println(result)

	return result
}

func eventRoom(nc *nats.Conn, userID string, rooms []*repository.GetRoomsByUserIDResult) {
	for _, r := range rooms {
		roomSubject := fmt.Sprintf("%s.%s", "room", r.RoomID)
		nc.Subscribe(roomSubject, func(msg *nats.Msg) {
			err := requestReadMessage(r, userID)
			if err != nil {
				log.Fatal(err)
			}

			var payload *api.Payload
			json.Unmarshal(msg.Data, &payload)
			log.Println("--------------------------------------------------------")
			log.Println("room: ", msg.Subject)
			log.Println("receiver: ", userID)
			log.Println("sender: ", payload.Sender)
			log.Println("msg: ", payload.Message)
			log.Println("--------------------------------------------------------")
		})

		focusedLobbySubject := fmt.Sprintf("focus.%s.%s", "lobby", r.RoomID)
		nc.Subscribe(focusedLobbySubject, func(msg *nats.Msg) {
			var payload *api.Payload
			json.Unmarshal(msg.Data, &payload)
			log.Println("--------------------------------------------------------")
			log.Println("lobby: ", msg.Subject)
			log.Println("receiver: ", userID)
			log.Println("sender: ", payload.Sender)
			log.Println("msg: ", payload.Message)
			log.Println("--------------------------------------------------------")
		})
	}
}
