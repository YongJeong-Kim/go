package recv

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"gounread/api"
	"gounread/repository"
	"io"
	"log"
	"net/http"
	"strings"
)

func GetConnectUserRooms(userID string) []*repository.GetRoomsByUserIDResult {
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

func ReadMessage(roomID, userID string) []*repository.GetMessagesByRoomIDAndTimeResult {
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

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result []*repository.GetMessagesByRoomIDAndTimeResult
	json.Unmarshal(data, &result)

	return result
}

func LoadMessageByRoomID(nc *nats.Conn, roomID, userID string) {
	messages := ReadMessage(roomID, userID)
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

	joinRoomSubject := fmt.Sprintf("join.room.%s", roomID)
	nc.Publish(joinRoomSubject, []byte(userID))
}
func EventRoom(nc *nats.Conn, userID string, rooms []*repository.GetRoomsByUserIDResult) {
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
			messages := ReadMessage(rm[1], userID)
			for _, m := range messages {
				if userID == payload.Sender {
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
