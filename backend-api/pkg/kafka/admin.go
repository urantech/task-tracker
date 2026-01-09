package kafka

import (
	"log"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

func CreateTopics(brokers []string, topics []string) error {
	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		return err
	}

	defer func() {
		if err = conn.Close(); err != nil {
			log.Printf("Error close connect: %v", err)
		}
	}()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}

	defer func() {
		if err := controllerConn.Close(); err != nil {
			log.Printf("Error close connect: %v", err)
		}
	}()

	topicConfigs := make([]kafka.TopicConfig, len(topics))
	for i, topic := range topics {
		topicConfigs[i] = kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}
	}

	return controllerConn.CreateTopics(topicConfigs...)
}
