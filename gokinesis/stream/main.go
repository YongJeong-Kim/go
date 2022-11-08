package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"gokinesis"
	"log"
)

var (
	kinesisStream AWSKinesis
)

// AWSKinesis struct , the collection of all field will be needed in kinesis stream
type AWSKinesis struct {
	stream          string
	region          string
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	sessionToken    string
}

// initiate configuration
func init() {
	cfg, err := gokinesis.LoadConfig("../")
	if err != nil {
		log.Fatal(err)
	}
	kinesisStream = AWSKinesis{
		stream:          cfg.StreamName,
		region:          cfg.AWSRegion,
		endpoint:        cfg.Endpoint,
		accessKeyID:     cfg.AccessKeyID,
		secretAccessKey: cfg.SecretAccessKey,
		//sessionToken:    cfg.
	}
}

func main() {
	action := flag.String("action", "create", "choose question `create` or `delete`")
	flag.Parse()

	// connect to aws-kinesis
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String(kinesisStream.region),
		Endpoint:    aws.String(kinesisStream.endpoint),
		Credentials: credentials.NewStaticCredentials(kinesisStream.accessKeyID, kinesisStream.secretAccessKey, kinesisStream.sessionToken),
	})
	if err != nil {
		log.Fatal("aws session failed. ", err)
	}
	kc := kinesis.New(s)
	streamName := aws.String(kinesisStream.stream)

	// create or delete kinesis stream name
	if *action == "create" {
		out, err := kc.CreateStream(&kinesis.CreateStreamInput{
			ShardCount: aws.Int64(1),
			StreamName: streamName,
		})
		if err != nil {
			log.Panic(err)
		}
		fmt.Printf("%v\n", out)

		if err := kc.WaitUntilStreamExists(&kinesis.DescribeStreamInput{StreamName: streamName}); err != nil {
			log.Panic(err)
		}
		log.Println("StreamName successfully created")
	} else if *action == "delete" {
		deleteOutput, err := kc.DeleteStream(&kinesis.DeleteStreamInput{
			StreamName: streamName,
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Delete successfully %v\n", deleteOutput)
	} else {
		fmt.Println("Wrong input")
	}
}
