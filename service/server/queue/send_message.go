package queue

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func (queue *SendMessage) SendMessage(input *SendMessageInput, output *SendMessageOutput) {
	var (
		producer          = input.producer
		destination_queue = input.queue_name
		content           = input.content
	)

	// Delivery report handler for produced messages
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	if err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &destination_queue, Partition: kafka.PartitionAny},
		Value:          []byte(content),
	}, nil); err != nil {
		log.Printf("Failed to produce message: %s to %s", err, destination_queue)
	} else {
		fmt.Printf("Produced message: %s sent to %s\n", content, destination_queue)
	}

	time.Sleep(1 * time.Second)
	producer.Flush(15 * 1000)

	defer producer.Close()
}
