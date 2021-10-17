package config

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func ConnectRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:16379",
		//Password: "1234",
		DB: 0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("connection fail redis. %s ", err.Error())
	}

	return client, nil
}
