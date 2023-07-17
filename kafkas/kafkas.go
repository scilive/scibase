package kafkas

import (
	"net"
	"strconv"
	"strings"

	"github.com/daqiancode/env"
	"github.com/segmentio/kafka-go"
)

func NewWriter() *kafka.Writer {
	broker := strings.Split(env.Get("KAFKA_BROKERS"), ",")
	return &kafka.Writer{
		Addr: kafka.TCP(broker...),
	}
}

func CreateTopic(topic string, partitions, replicationFactor int) error {
	broker := strings.Split(env.Get("KAFKA_BROKERS"), ",")
	conn, err := kafka.Dial("tcp", broker[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     partitions,
			ReplicationFactor: replicationFactor,
		},
	}

	return controllerConn.CreateTopics(topicConfigs...)
}
