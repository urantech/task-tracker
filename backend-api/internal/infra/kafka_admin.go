package infra

import (
	"backend-api/internal/config"
	"log"
	"net"
	"strconv"

	"github.com/segmentio/kafka-go"
)

func CreateTopics(cfg *config.ProducerConfig) error {
	topics := []string{cfg.WelcomeTopic, cfg.ReportTopic}

	conn, err := kafka.Dial("tcp", cfg.Brokers[0])
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
		if err = controllerConn.Close(); err != nil {
			log.Printf("Error close connect: %v", err)
		}
	}()

	numPart := cfg.NumPartitions
	repFactor := cfg.ReplicationFactor

	topicConfigs := make([]kafka.TopicConfig, len(topics))
	for i, topic := range topics {
		topicConfigs[i] = kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     numPart,
			ReplicationFactor: repFactor,
		}
	}

	return controllerConn.CreateTopics(topicConfigs...)
}
