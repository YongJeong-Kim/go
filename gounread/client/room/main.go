package main

import (
	"flag"
	"github.com/nats-io/nats.go"
	"gounread/embedded"
	"log"
	"strings"
)

const (
	user1 = "05f84f46-d4ad-42db-af4f-da63cffcb721"
	user2 = "febba554-152e-496a-add5-31d0672fdc2a"
)

func main() {
	userID := flag.String("user", user1, "input user id")
	flag.Parse()

	nc := natsConnect()

	//rooms := getConnectUserRooms(*userID)
	//eventRoom(nc, *userID, rooms)

	eventJoinRoom(nc, "01f84cfa-e487-494c-82e5-e75f95ef0573", *userID)

	c := make(chan struct{})
	<-c
}

func natsConnect() *nats.Conn {
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

	return nc
}
