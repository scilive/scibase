package kafkas

import (
	"context"

	"github.com/scilive/scibase/logs"
	"github.com/segmentio/kafka-go"
)

type KafkaReader struct {
	Reader  *kafka.Reader
	Handler func(m kafka.Message) error
}

func NewReader(config kafka.ReaderConfig) *KafkaReader {
	return &KafkaReader{
		Reader: kafka.NewReader(config),
	}
}

func (s *KafkaReader) Run() {
	defer func() {
		if err := recover(); err != nil {
			logs.Log.Error().Err(err.(error)).Msg("panic in kafka reader")
		}
	}()
	s.Reader.SetOffset(kafka.LastOffset)
	for {
		m, err := s.Reader.ReadMessage(context.Background())
		if err != nil {
			logs.Log.Error().Err(err).Msg("failed to read message")
			continue
		}
		if s.Handler != nil {
			err := s.Handler(m)
			if err != nil {
				logs.Log.Error().Err(err).Str("topic", m.Topic).Msg("failed to handle message")
			}
		} else {
			logs.Log.Warn().Str("topic", m.Topic).Msg("no handler for topic")
		}

	}
}
