package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"gokinesis"
	"log"
	"time"
)

type Producer struct {
	stream          string
	region          string
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	sessionToken    string
}

func NewProducer(cfg gokinesis.Config) *Producer {
	return &Producer{
		stream:          cfg.StreamName,
		region:          cfg.AWSRegion,
		endpoint:        cfg.Endpoint,
		accessKeyID:     cfg.AccessKeyID,
		secretAccessKey: cfg.SecretAccessKey,
		//sessionToken:    os.Getenv("AWS_SESSION_TOKEN"),
	}
}

func main() {
	cfg, err := gokinesis.LoadConfig("../")
	if err != nil {
		log.Fatal("load config failed. ", err)
	}
	producer := NewProducer(cfg)

	s, err := session.NewSession(&aws.Config{
		Region: aws.String(producer.region),
		//Endpoint:    aws.String(producer.endpoint),
		Credentials: credentials.NewStaticCredentials(producer.accessKeyID, producer.secretAccessKey, producer.sessionToken),
	})
	if err != nil {
		log.Fatal("producer aws session failed. ", err)
	}
	kc := kinesis.New(s)
	streamName := aws.String(producer.stream)
	_, err = kc.DescribeStream(&kinesis.DescribeStreamInput{StreamName: streamName})

	//if no stream name in AWS
	if err != nil {
		log.Panic(err)
	}

	tick := time.Tick(33 * time.Millisecond)

	for {
		select {
		case <-tick:
			data := `{
				"id": "%s"
			}`
			d := fmt.Sprintf(data, time.Now().String())
			putOutput, err := kc.PutRecord(&kinesis.PutRecordInput{
				Data:         []byte(d),
				StreamName:   streamName,
				PartitionKey: aws.String("key1"),
			})
			if err != nil {
				panic(err)
			}
			fmt.Printf("%v\n", *putOutput)
			//fmt.Println(time.Now().String())
		}
	}
	//ns := [1000]int{}
	//for i := range ns {
	//	data := `{
	//		"id": "%d"
	//	}`
	//
	//	d := fmt.Sprintf(data, i)
	//	putOutput, err := kc.PutRecord(&kinesis.PutRecordInput{
	//		Data:         []byte(d),
	//		StreamName:   streamName,
	//		PartitionKey: aws.String("key1"),
	//	})
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Printf("%v\n", *putOutput)
	//}
}
