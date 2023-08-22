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

// Run `go reader.Run()` in a goroutine
func (s *KafkaReader) Run() {
	defer func() {
		s.Reader.Close()
		if err := recover(); err != nil {
			if e, ok := err.(error); ok {
				logs.Log.Error().Err(e).Msg("panic in kafka reader")
			} else {
				logs.Log.Error().Interface("error", e).Msg("panic in kafka reader")
			}
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
