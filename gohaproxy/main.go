package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
)

type Server struct {
	redis *redis.Client
}

func NewServer(redis *redis.Client) *Server {
	return &Server{
		redis: redis,
	}
}

func (server *Server) Publish(msg string) {
	err := server.redis.Publish(context.Background(), "channel1", "payload").Err()
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) Subscribe() {
	ctx := context.Background()
	sub := server.redis.Subscribe(ctx, "channel1")
	defer sub.Close()

	for {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("received: %s\n", msg)
	}
}

func main() {
	redis := redis.NewClient(&redis.Options{
		Addr:     "localhost:36379",
		Password: "",
	})

	server := NewServer(redis)

	go server.Subscribe()

	a := os.Getenv("aaa")
	log.Println(a)

	r := gin.New()
	r.GET("/channel1/:msg", func(c *gin.Context) {
		var req struct {
			Msg string `uri:"msg"`
		}
		if err := c.ShouldBindUri(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid message",
			})
			return
		}

		server.Publish(req.Msg)
		c.JSON(http.StatusOK, gin.H{
			"ok": "ok",
		})
	})

	r.Run(":8080")
}
