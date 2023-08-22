package kafkas

import (
	"context"
	"encoding/json"

	"github.com/scilive/scibase/logs"
	"github.com/segmentio/kafka-go"
)

type KafkaWriter struct {
	Writer *kafka.Writer
}

// Write writes a message to kafka, write payload as json
func (s *KafkaWriter) Write(context context.Context, topic string, payload any) error {
	var msg kafka.Message
	msg.Topic = topic
	bs, err := json.Marshal(payload)
	if err != nil {
		logs.Log.Error().Err(err).Msg("failed to marshal message")
		return err
	}
	msg.Value = bs
	return s.Writer.WriteMessages(context, msg)
}
