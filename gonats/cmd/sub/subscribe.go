package main

import (
	"github.com/nats-io/nats.go"
	"log"
	_ "net/http/pprof"
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
		nats.DisconnectErrHandler(func(conn *nats.Conn, err error) {
			log.Println("disconnected from server", err)
		}),
		nats.ClosedHandler(func(conn *nats.Conn) {
			log.Println("connection closed", conn)
		}),
		nats.ReconnectHandler(func(conn *nats.Conn) {
			log.Println("reconnected to server")
		}),
		nats.DiscoveredServersHandler(func(conn *nats.Conn) {
			log.Println("discovered servers")
		}),
		nats.ConnectHandler(func(conn *nats.Conn) {
			log.Println("connected to server")
		}),
		nats.MaxReconnects(60),
		nats.ReconnectWait(500*time.Millisecond),
		nats.CustomReconnectDelay(func(attempts int) time.Duration {
			log.Println("attempts: ", attempts)
			return 0
		}),
		//nats.ReconnectBufSize(5*1024*1024),
		//nats.RetryOnFailedConnect(true),
		//nats.DontRandomize(), // avoiding thundering herd
		//nats.ReconnectJitter(0, 0),
	)
	if err != nil {
		log.Fatal("connection error. ", err)
	}
	defer ns.Close()

	sub, err := ns.Subscribe("a.*", func(msg *nats.Msg) {
		log.Printf("sub: %s, msg: %s\n", msg.Subject, msg.Data)
	})
	if err != nil {
		log.Fatal("subs error. ", err)
	}
	_ = sub

	ns.Subscribe("a.b.>", func(msg *nats.Msg) {
		log.Printf("sub: %s, msg: %s\n", msg.Subject, msg.Data)
	})

	c := make(chan struct{})
	<-c
}
