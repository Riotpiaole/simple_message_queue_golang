package queue

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

// ==========CreateQueue API =============
type CreateQueue int32
type CreateQueueOutput struct {
	queue_name string
	user_id    string
}
type CreateQueueInput struct {
	queue_name string
}

// ==========sendMessage API =============
type SendMessage int32
type SendMessageOutput struct {
	timestamp int64
}
type SendMessageInput struct {
	producer   kafka.Producer
	queue_name string
	content    string
}
