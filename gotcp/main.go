package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	server, err := NewServer("localhost:8080")
	if err != nil {
		log.Println("create new server failed. ", err)
		os.Exit(1)
	}
	server.start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("shutdown")
	server.stop()
	log.Println("server stopped")
}

func (s *Server) accept() {
	defer s.wg.Done()

	for {
		select {
		case <-s.shutdown:
		default:
			conn, err := s.listener.Accept()
			if err != nil {

			}
			b := make([]byte, 4096)
			n, err := conn.Read(b)
			if err != nil {

			}
			log.Println("receive: ", string(b[:n]))
			s.conn <- conn
		}
	}
}

func (s *Server) connection() {
	defer s.wg.Done()
	for {
		select {
		case c := <-s.conn:
			go s.handle(c)
		case <-s.shutdown:

		}
	}
}

func (s *Server) handle(conn net.Conn) {
	log.Println("start handle")
	time.Sleep(2 * time.Second)
	log.Println("after 2 secs")
	str := "ok"
	n, err := conn.Write([]byte(str))
	if err != nil {

	}
	log.Println(n)
}

func (s *Server) start() {
	s.wg.Add(2)
	go s.accept()
	go s.connection()
}

func (s *Server) stop() {
	close(s.shutdown)
	s.listener.Close()
}

type Server struct {
	listener net.Listener
	shutdown chan struct{}
	conn     chan net.Conn
	wg       sync.WaitGroup
}

func NewServer(address string) (*Server, error) {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		return nil, fmt.Errorf("create conn failed. %v", err)
	}

	return &Server{
		listener: listener,
		shutdown: make(chan struct{}),
		conn:     make(chan net.Conn),
	}, nil
}
