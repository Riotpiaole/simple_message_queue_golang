package main

import (
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func getHost() string {
	const KAFKA_BROKER_SERV = "KAFKA_BROKER_SERV"
	host := os.Getenv(KAFKA_BROKER_SERV)
	if len(host) == 0 {
		host = "localhost:29092"
	}
	log.Printf("Consumer connecting to Kafka broker at [%v]", host)
	return host
}

func main() {
	// Initialize Kafka consumer
	serv := getHost()
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": serv,
		"group.id":          "my-group",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	defer consumer.Close()

	// Subscribe to the 'test-topic' Kafka topic
	err = consumer.Subscribe("test-topic", nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic: %s", err)
		os.Exit(1)
	}

	// Consume messages from the topic
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Received message: %s\n", string(msg.Value))
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
