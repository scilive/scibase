package kafkas

import (
	"net"
	"strconv"
	"strings"

	"github.com/daqiancode/env"
	"github.com/segmentio/kafka-go"
)

func GetReaderConfigEnv(groupId string) kafka.ReaderConfig {
	brokers := strings.Split(env.Get("KAFKA_BROKERS"), ",")
	return kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupId,
		MaxBytes: env.GetIntMust("KAFKA_MAX_BYTES", 1*1024*1024),
	}
}
func NewReaderEnv(groupId string) *KafkaReader {
	return &KafkaReader{
		Reader: kafka.NewReader(GetReaderConfigEnv(groupId))}
}

func NewReader(config kafka.ReaderConfig) *KafkaReader {
	return &KafkaReader{
		Reader: kafka.NewReader(config),
	}
}
func NewWriterEnv(autoCreateTopic bool) *kafka.Writer {
	brokers := strings.Split(env.Get("KAFKA_BROKERS"), ",")
	return &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		AllowAutoTopicCreation: autoCreateTopic,
	}
}

func NewWriter(brokers []string, autoCreateTopic bool) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		AllowAutoTopicCreation: autoCreateTopic,
	}
}

func CreateTopic(topic string, partitions, replicationFactor int) error {
	brokers := strings.Split(env.Get("KAFKA_BROKERS"), ",")
	conn, err := kafka.Dial("tcp", brokers[0])
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
