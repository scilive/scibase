package kafkas

import "github.com/segmentio/kafka-go"

type KafkaReaderHub struct {
	brokers       []string
	readers       map[string]*KafkaReader
	defaultConfig kafka.ReaderConfig
}

func NewReaderHub(brokers []string, defaultConfig kafka.ReaderConfig) *KafkaReaderHub {
	return &KafkaReaderHub{
		brokers:       brokers,
		readers:       make(map[string]*KafkaReader),
		defaultConfig: defaultConfig,
	}
}

type Handler func(m kafka.Message) error

func (s *KafkaReaderHub) On(topic, groupId string, handler Handler) {
	if s.readers[topic] == nil {
		config := s.defaultConfig
		config.Topic = topic
		config.GroupID = groupId
		s.readers[topic] = NewReader(config)
	}
	s.readers[topic].Handler = handler
}

func (s *KafkaReaderHub) Run() {
	for _, reader := range s.readers {
		go reader.Run()
	}
}
