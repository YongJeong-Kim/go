package main

import (
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {

	}

	text := "welcome"
	n, err := conn.Write([]byte(text))
	log.Println(n)

	b1 := make([]byte, 4096)
	r, err := conn.Read(b1)
	if err != nil {

	}
	log.Println(string(b1[:r]))
}
