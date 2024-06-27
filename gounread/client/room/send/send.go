package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	user1    = "05f84f46-d4ad-42db-af4f-da63cffcb721"
	room1    = "65934e56-02d7-40bc-b747-59017e49759b"
	message1 = "hihi hihi"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	userID := flag.String("user", user1, "input user id")
	roomID := flag.String("room", room1, "input room id")
	message := flag.String("message", message1, "input message")
	flag.Parse()

	send(*roomID, *userID, *message)
}

func send(roomID, userID, message string) {
	url := fmt.Sprintf("http://localhost:8080/rooms/%s/send", roomID)
	body := []byte(fmt.Sprintf(`
		{
			"message": "%s"
		}
	`, message))
	/*	body := map[string]string{
		"message": message,
	}*/
	//b, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		log.Fatal("new request send message failed. ", err)
	}

	req.Header.Set("user", userID)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal("request send message failed. ", err)
	}

	if resp.StatusCode != http.StatusCreated {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("read all body error.", err)
		}
		var errMap map[string]any
		err = json.Unmarshal(b, &errMap)
		if err != nil {
			log.Fatal("unmarshal error.", err)
		}
		log.Fatal("cannot send message. ", errMap["error"])
	}
}
