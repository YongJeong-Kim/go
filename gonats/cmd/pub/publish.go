package main

import (
	"github.com/nats-io/nats.go"
	"log"
	"strings"
	"time"
)

var servers = []string{
	"nats://localhost:14222",
	"nats://localhost:24222",
	"nats://localhost:34222",
}

func main() {
	ns, err := nats.Connect(
		strings.Join(servers, ","),
		nats.ConnectHandler(func(conn *nats.Conn) {
			log.Println("publish connection success")
		}),
	)

	if err != nil {
		log.Fatal("connection error. ", err)
	}
	defer ns.Drain()

	//for range 10 {
	//	err = ns.Publish("a.b", []byte("to a.b message"))
	//	if err != nil {
	//		log.Println("to a.b send message error.", err)
	//	}
	//
	//	err = ns.Publish("a.b.c", []byte("to a.b.c message"))
	//	if err != nil {
	//		log.Println("to a.b.c send message error.", err)
	//	}
	//	time.Sleep(500 * time.Millisecond)
	//}

	err = ns.Publish("a.b.c", []byte("to a.b.c message"))
	if err != nil {
		log.Println("to a.b.c send message error.", err)
	}
	time.Sleep(500 * time.Millisecond)
}
