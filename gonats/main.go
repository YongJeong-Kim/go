package main

import (
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
	"time"
)

//var servers = []string{"nats://localhost:14222", "nats://localhost:24222", "nats://localhost:34222"}

// var servers = []string{"nats://localhost:14222", "nats://localhost:24222", "nats://localhost:34222"}
var servers = []string{"nats://localhost:14222"}

func main() {
	opts := nats.Options{
		Servers:              servers,
		RetryOnFailedConnect: true,
		ReconnectedCB: func(conn *nats.Conn) {
			log.Println("reconnect success")
		},
		MaxReconnect:  5,
		ReconnectWait: 2 * time.Second, // default 2 seconds
		CustomReconnectDelayCB: func(attempts int) time.Duration {
			log.Println("attempts: ", attempts)
			return 2 * time.Second
		},
		DisconnectedErrCB: func(conn *nats.Conn, err error) {
			log.Println("disconnect error.", err)
			conn.ReconnectHandler()
		},
		ConnectedCB: func(conn *nats.Conn) {
			log.Println("connect success")
			_, err := conn.Subscribe("a.b.c", func(msg *nats.Msg) {
				log.Println("sub msg: ", string(msg.Data))
			})
			if err != nil {
				log.Fatal(err)
			}
			//defer sub.Unsubscribe()

			//time.Sleep(15 * time.Second)
			//conn.Close()
		},
		ClosedCB: func(conn *nats.Conn) {
			log.Println("close conn")
		},
		AsyncErrorCB: func(conn *nats.Conn, subscription *nats.Subscription, err error) {
			log.Println("async error.", err)
		},
		AllowReconnect: true,
		User:           "aaa",
		Password:       "1234",
	}

	errCnt := 0
	for range 1030 {
		go func(i int) {
			_, err := opts.Connect()
			if err != nil {
				errCnt++
				log.Fatal("connect error. ", err)
			}
		}(errCnt)
	}

	http.ListenAndServe("localhost:6060", nil)
}
