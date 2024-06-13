package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nats-io/nats.go"
	"gounread/api"
	"gounread/embedded"
	"gounread/service"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	// 05f84f46-d4ad-42db-af4f-da63cffcb721
	// febba554-152e-496a-add5-31d0672fdc2a
	userID := flag.String("user", "05f84f46-d4ad-42db-af4f-da63cffcb721", "input user id")
	flag.Parse()

	nc, err := nats.Connect(
		strings.Join(embedded.Servers, ","),
		nats.ConnectHandler(func(conn *nats.Conn) {
			log.Println("connected to server")
		}),
		nats.UserInfo("aaa", "1234"),
	)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/connect", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("user", *userID)
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

	var result []*service.GetRoomsByUserIDResult
	json.Unmarshal(data, &result)
	log.Println(result)
	for _, r := range result {
		joinRoomSubject := fmt.Sprintf("join.room.%s", r.RoomID)
		nc.Subscribe(joinRoomSubject, func(msg *nats.Msg) {

		})

		roomSubject := fmt.Sprintf("%s.%s", "room", r.RoomID)
		nc.Subscribe(roomSubject, func(msg *nats.Msg) {
			err = requestReadMessage(r, *userID)
			if err != nil {
				log.Fatal(err)
			}

			var payload *api.Payload
			json.Unmarshal(msg.Data, &payload)
			log.Println("--------------------------------------------------------")
			log.Println("room: ", msg.Subject)
			log.Println("receiver: ", *userID)
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
			log.Println("receiver: ", *userID)
			log.Println("sender: ", payload.Sender)
			log.Println("msg: ", payload.Message)
			log.Println("--------------------------------------------------------")
		})
	}

	c := make(chan struct{})
	<-c
}

func requestReadMessage(r *service.GetRoomsByUserIDResult, userID string) error {
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
