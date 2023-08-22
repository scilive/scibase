package kafkas

import "github.com/segmentio/kafka-go"

type KafkaReaderHub struct {
	readers       map[string]*KafkaReader
	defaultConfig kafka.ReaderConfig
}

func NewReaderHub(defaultConfig kafka.ReaderConfig) *KafkaReaderHub {
	return &KafkaReaderHub{
		readers:       make(map[string]*KafkaReader),
		defaultConfig: defaultConfig,
	}
}

type Handler func(m kafka.Message) error

func (s *KafkaReaderHub) On(topic, groupId string, handler Handler) {
	config := s.defaultConfig
	config.Topic = topic
	config.GroupID = groupId
	s.readers[topic] = NewReader(config)
}

func (s *KafkaReaderHub) Run() {
	for _, reader := range s.readers {
		go reader.Run()
	}
}
