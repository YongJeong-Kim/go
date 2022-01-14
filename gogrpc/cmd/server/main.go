package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {
	port := flag.Int("port", 0, "server port")
	flag.Parse()

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server")
	}

	err = http.Serve(listener, nil)
	if err != nil {
		return
	}
}
