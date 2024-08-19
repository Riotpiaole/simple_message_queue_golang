package queue

import (
	"log"
	"net"
	"net/rpc"
)

func (s *CreateQueue) CreateQueue(input *CreateQueueInput, output *CreateQueueOutput) error {
	return nil
}

func CreateQueueListener() *net.Listener {
	queue_server := new(CreateQueue)
	rpc.Register(queue_server)
	rpc.HandleHTTP()

	listener, err := net.Listen("tpc", ":3300")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	return &listener
}
