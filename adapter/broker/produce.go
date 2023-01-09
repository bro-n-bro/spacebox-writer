package broker

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func (b *Broker) produce(topic string, data []byte, headers []kafka.Header) error {
	err := b.pr.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          data,
		Headers:        headers,
	}, nil)

	if kafkaError, ok := err.(kafka.Error); ok && kafkaError.Code() == kafka.ErrQueueFull {
		b.log.Info().Str("topic", topic).Msg("Kafka local queue full error - Going to Flush then retry...")
		flushedMessages := b.pr.Flush(30000)
		b.log.Info().Str("topic", topic).
			Msgf("Flushed kafka messages. Outstanding events still un-flushed: %d", flushedMessages)
		return b.produce(topic, data, headers)
	}

	return nil
}
