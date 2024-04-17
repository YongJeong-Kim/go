package main

import (
	"context"
	"github.com/gin-gonic/gin"
	accountv1 "github.com/yongjeong-kim/go/goapigw/accountpb/pb/account/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient("localhost:19090", opts...)
	if err != nil {
		log.Fatal("connect account client error", err)
	}

	accountClient := accountv1.NewAccountServiceClient(conn)

	r := gin.New()
	r.GET("/service-b/api/v1/111", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "ok1",
		})
		return
	})
	r.GET("/service-b/api/v1/111/111", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "bbb111",
		})
		return
	})
	r.POST("/service-b/api/v1/account", func(c *gin.Context) {
		resp, err := accountClient.CreateAccount(context.Background(), &accountv1.CreateAccountRequest{
			Account: &accountv1.Account{
				AccountId: "12",
			},
		})
		if err != nil {
			s := status.Convert(err)
			log.Println("create account error", err)
			log.Println(s)
			if s.Code() == 14 {
				c.Status(http.StatusServiceUnavailable)
				return
			}
		}
		c.JSON(http.StatusOK, resp)
	})

	r.Run(":8081")
}
