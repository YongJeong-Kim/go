package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nats-io/nats.go"
	"gounread/api"
	"gounread/embedded"
	"gounread/repository"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	user1 = "05f84f46-d4ad-42db-af4f-da63cffcb721"
	user2 = "febba554-152e-496a-add5-31d0672fdc2a"
	room1 = "01f84cfa-e487-494c-82e5-e75f95ef0573"
	room2 = "65934e56-02d7-40bc-b747-59017e49759b"
)

func main() {
	userID := flag.String("user", user1, "input user id")
	roomID := flag.String("room", room1, "input room id")
	flag.Parse()

	nc := natsConnect(*roomID, *userID)

	loadMessageByRoomID(nc, *roomID, *userID)
	rooms := getConnectUserRooms(*userID)
	eventRoom(nc, *userID, rooms)

	c := make(chan struct{})
	<-c
}

func natsConnect(roomID, userID string) *nats.Conn {
	nc, err := nats.Connect(
		strings.Join(embedded.Servers, ","),
		nats.ConnectHandler(func(conn *nats.Conn) {
			log.Printf("connected to server. room id: %s, user id: %s", roomID, userID)
		}),
		nats.UserInfo("aaa", "1234"),
	)
	if err != nil {
		log.Fatal(err)
	}

	return nc
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

func readMessage(roomID, userID string) []*repository.GetMessagesByRoomIDAndTimeResult {
	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/rooms/"+roomID+"/read", nil)
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
		errBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("read error body error.", err)
		}
		var errMap map[string]string
		json.Unmarshal(errBody, &errMap)
		log.Fatal("request read message error.", errMap["error"])
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result []*repository.GetMessagesByRoomIDAndTimeResult
	json.Unmarshal(data, &result)

	return result
}

func loadMessageByRoomID(nc *nats.Conn, roomID, userID string) {
	messages := readMessage(roomID, userID)
	for _, m := range messages {
		log.Println("-------------------------------------------------------")
		log.Printf("%s, %s\n", m.RoomID, m.Sent)
		if userID == m.Sender {
			log.Printf("%s <- %s, %d\n", m.Msg, m.Sender, len(m.Unread))
		} else {
			log.Printf("%d, %s -> %s\n", len(m.Unread), m.Sender, m.Msg)
		}
		log.Printf("unread users: %s\n", m.Unread)
		log.Println("-------------------------------------------------------")
	}

	joinRoomSubject := fmt.Sprintf("join.room.%s", roomID)
	nc.Publish(joinRoomSubject, []byte(userID))
}
func eventRoom(nc *nats.Conn, userID string, rooms []*repository.GetRoomsByUserIDResult) {
	for _, r := range rooms {
		joinRoomSubject := fmt.Sprintf("join.room.%s", r.RoomID)
		nc.Subscribe(joinRoomSubject, func(msg *nats.Msg) {
			connectUserID := string(msg.Data)
			if connectUserID != userID {
				log.Println("connected user:", connectUserID)
				log.Println("must be previous unread messages count recalculate")
			}
		})

		roomSubject := fmt.Sprintf("room.%s", r.RoomID)
		nc.Subscribe(roomSubject, func(msg *nats.Msg) {
			var payload *api.Payload
			json.Unmarshal(msg.Data, &payload)

			rm := strings.Split(roomSubject, ".")
			messages := readMessage(rm[1], userID)
			for _, m := range messages {
				log.Println("-------------------------------------------------------")
				log.Printf("%s, %s\n", m.RoomID, m.Sent)
				if userID == payload.Sender {
					log.Printf("%s <- %s, %d\n", m.Msg, m.Sender, len(m.Unread))
				} else {
					log.Printf("%d, %s -> %s\n", len(m.Unread), m.Sender, m.Msg)
				}
				log.Printf("unread users: %s\n", m.Unread)
				log.Println("-------------------------------------------------------")
			}
		})

		/*	focusedLobbySubject := fmt.Sprintf("focus.%s.%s", "lobby", r.RoomID)
			nc.Subscribe(focusedLobbySubject, func(msg *nats.Msg) {
				var payload *api.Payload
				json.Unmarshal(msg.Data, &payload)
				log.Println("--------------------------------------------------------")
				log.Println("lobby: ", msg.Subject)
				log.Println("receiver: ", userID)
				log.Println("sender: ", payload.Sender)
				log.Println("msg: ", payload.Message)
				log.Println("--------------------------------------------------------")
			})*/
	}
}
