package main

import (
	"bufio"
	"github.com/IBM/sarama"
	"log"
	"os"
	"time"
)

func main() {
	brokers := []string{"localhost:19092"}
	producer := NewAsyncProducer(brokers, sarama.DefaultVersion)

	readStart := time.Now().UTC()
	file, err := os.Open("../../MOCK_DATA.csv")
	if err != nil {
		log.Fatal(err)
	}
	readEnd := time.Since(readStart)
	log.Printf("read done: %s", readEnd.String())
	defer file.Close()
	defer producer.Close()

	produceStart := time.Now().UTC()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//log.Println(scanner.Text())
		producer.Input() <- &sarama.ProducerMessage{
			Topic: "bbb",
			Value: sarama.StringEncoder(scanner.Text()),
		}
	}
	produceEnd := time.Since(produceStart)
	log.Printf("produce done: %s", produceEnd.String())
}

func NewAsyncProducer(brokers []string, version sarama.KafkaVersion) sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Version = version
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	/*tlsConfig := createTlsConfiguration()
	if tlsConfig != nil {
		config.Net.TLS.Config = tlsConfig
		config.Net.TLS.Enable = true
	}*/

	// On the broker side, you may want to change the following settings to get
	// stronger consistency guarantees:
	// - For your broker, set `unclean.leader.election.enable` to false
	// - For the topic, you could increase `min.insync.replicas`.

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	return producer
}
