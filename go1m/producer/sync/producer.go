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
	producer := NewSyncProducer(brokers, sarama.DefaultVersion)

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
		producer.SendMessage(&sarama.ProducerMessage{
			Topic: "bbb",
			Value: sarama.StringEncoder(scanner.Text()),
		})
	}
	produceEnd := time.Since(produceStart)
	log.Printf("produce done: %s", produceEnd.String())
}

func NewSyncProducer(brokers []string, version sarama.KafkaVersion) sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Version = version
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	//config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	//config.Producer.Partitioner = sarama.NewHashPartitioner // default
	//config.Producer.Partitioner = sarama.NewManualPartitioner
	//config.Producer.Partitioner = sarama.NewRandomPartitioner
	/*tlsConfig := createTlsConfiguration()
	if tlsConfig != nil {
		config.Net.TLS.Config = tlsConfig
		config.Net.TLS.Enable = true
	}*/

	// On the broker side, you may want to change the following settings to get
	// stronger consistency guarantees:
	// - For your broker, set `unclean.leader.election.enable` to false
	// - For the topic, you could increase `min.insync.replicas`.

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	return producer
}
