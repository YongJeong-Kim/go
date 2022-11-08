package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"gokinesis"
	"log"
)

type Consumer struct {
	stream          string
	region          string
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	sessionToken    string
}

func NewConsumer(cfg gokinesis.Config) *Consumer {
	return &Consumer{
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
	consumer := NewConsumer(cfg)
	// connect to aws-kinesis
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(consumer.region),
		//Endpoint:    aws.String(consumer.endpoint),
		Credentials: credentials.NewStaticCredentials(consumer.accessKeyID, consumer.secretAccessKey, consumer.sessionToken),
	})
	if err != nil {
		log.Fatal("consumer aws session failed. ", err)
	}
	kc := kinesis.New(s)
	streamName := aws.String(consumer.stream)
	streams, err := kc.DescribeStream(&kinesis.DescribeStreamInput{StreamName: streamName})
	if err != nil {
		log.Panic(err)
	}

	// retrieve iterator
	iteratorOutput, err := kc.GetShardIterator(&kinesis.GetShardIteratorInput{
		// Shard Id is provided when making put record(s) request.
		ShardId:           aws.String(*streams.StreamDescription.Shards[3].ShardId),
		ShardIteratorType: aws.String("TRIM_HORIZON"),
		//ShardIteratorType: aws.String("AT_SEQUENCE_NUMBER"),
		//ShardIteratorType: aws.String("LATEST"),
		StreamName: streamName,
		//StartingSequenceNumber: aws.String("49633056657994085916131401762705608327794518669689618482"),
	})
	if err != nil {
		log.Panic(err)
	}

	shardIterator := iteratorOutput.ShardIterator
	var a *string

	// get data using infinity looping
	// we will attempt to consume data every 1 secons, if no data, nothing will be happen
	for {
		// get records use shard iterator for making request
		records, err := kc.GetRecords(&kinesis.GetRecordsInput{
			ShardIterator: shardIterator,
		})

		// if error, wait until 1 seconds and continue the looping process
		if err != nil {
			//time.Sleep(1000 * time.Millisecond)
			continue
		}

		// process the data
		if len(records.Records) > 0 {
			for _, d := range records.Records {
				m := make(map[string]interface{})
				err := json.Unmarshal(d.Data, &m)
				if err != nil {
					log.Println(err)
					continue
				}
				log.Printf("GetRecords Data: %v\n", m)
			}
		} else if records.NextShardIterator == a || shardIterator == records.NextShardIterator || err != nil {
			log.Printf("GetRecords ERROR: %v\n", err)
			break
		}
		shardIterator = records.NextShardIterator
		//time.Sleep(1000 * time.Millisecond)
	}
}
