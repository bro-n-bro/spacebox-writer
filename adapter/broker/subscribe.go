package broker

import (
	"context"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/bro-n-bro/spacebox-writer/adapter/mongo/model"
	"github.com/bro-n-bro/spacebox-writer/internal/rep"
)

const (
	keyMsg       = "msg"
	keyPartition = "partition"
	keyOffset    = "offset"
	keyTopic     = "topic"
	keyRetry     = "retry"

	msgStopReadingMessagesFromTopic = "stop reading messages from topic"
	msgCreateBrokerMessageError     = "CreateBrokerMessage error"
	msgReadMessageError             = "read message error"
	msgCommitMessageError           = "commit message error"
	msgHandleMessageError           = "handle message error"

	logLayout = "[%v]: %s"

	readMessageTimeout = 100
)

// Subscribe subscribes to kafka topic.
func (b *Broker) Subscribe(
	ctx context.Context, //
	wg *sync.WaitGroup, //
	topic string, // name of topic to subscribe
	handler func(ctx context.Context, msg [][]byte, db rep.Storage) error, // handler for processing messages
) error {

	defer wg.Done()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        b.cfg.Address,
		"group.id":                 b.cfg.GroupID,
		"auto.offset.reset":        b.cfg.AutoOffsetReset,
		"allow.auto.create.topics": true,
		"enable.auto.offset.store": false,
	})

	if err != nil {
		return err
	}

	if err = consumer.Subscribe(topic, nil); err != nil {
		return err
	}

	b.consumers = append(b.consumers, consumer)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		var (
			batch  = newBatch(*b.log, topic, b.cfg.BatchBufferSize, handler)
			ticker = time.NewTicker(b.cfg.FlushBufferInterval)
		)

		batch.setErrorHandler(b.handleError)

		if b.cfg.MetricsEnabled {
			batch.setMetrics(b.metrics)
		}

		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				batch.flushBuffer(context.Background(), b.st)
				b.log.Info().Str(keyTopic, topic).Msg(msgStopReadingMessagesFromTopic)
				return
			case <-ticker.C:
				batch.flushBuffer(ctx, b.st)
			default:
			}

			msg, err := consumer.ReadMessage(readMessageTimeout)
			if msg == nil {
				continue
			}

			if err != nil {
				b.log.Fatal().Err(err).
					Str(keyMsg, string(msg.Value)).
					Msg(msgReadMessageError)
				return
			}

			b.log.Debug().Msgf(logLayout, msg.String(), msg.Value)

			batch.insertMessage(ctx, msg, b.st)

			if _, err = consumer.CommitMessage(msg); err != nil {
				b.log.Error().Err(err).
					Str(keyTopic, topic).
					Msg(msgCommitMessageError)
			}
		}
	}(wg)

	return nil
}

// handleError processes an error of handle function for a consumer if needed.
// Writes to mongo info about the error if message handling call times exceeded MAX_RETRIES config setting.
//
// Do nothing if the message handling does not return an error or messages are empty.
func (b *Broker) handleError(ctx context.Context, messageHandlerError error, msgs []*kafka.Message,
	handler func(ctx context.Context, msg [][]byte, db rep.Storage) error) {

	if len(msgs) == 0 || messageHandlerError == nil { // handling error not needed
		return
	}

	// topic the same for all messages
	topic := *msgs[0].TopicPartition.Topic

	if b.cfg.MetricsEnabled {
		b.metrics.errorsCounter.With(prometheus.Labels{keyTopic: topic}).Inc()
	}

	// error occurred. try to handle each message
	for _, msg := range msgs {
		// handle single message cfg.MaxRetries times
		for retry := 0; retry <= b.cfg.MaxRetries; retry++ {
			if err := handler(ctx, [][]byte{msg.Value}, b.st); err == nil { // success handling
				return // bellow steps not needed
			}
		}

		// retry limit exceeded
		if b.cfg.MetricsEnabled {
			b.metrics.limitExceededCounter.With(prometheus.Labels{keyTopic: topic}).Inc()
		}

		// got error of handling message from the broker
		b.log.Error().Err(messageHandlerError).
			Str(keyTopic, topic).
			Int64(keyOffset, int64(msg.TopicPartition.Offset)).
			Int64(keyPartition, int64(msg.TopicPartition.Partition)).
			Int(keyRetry, b.cfg.MaxRetries).
			Str(keyMsg, string(msg.Value)).
			Msg(msgHandleMessageError)

		// write error message to mongo
		if err := b.m.CreateBrokerMessage(ctx, &model.BrokerMessage{
			LastErrorMessage: messageHandlerError.Error(),
			Topic:            topic,
			Data:             string(msg.Value),
			Attempts:         b.cfg.MaxRetries,
			Created:          time.Now(),
		}); err != nil {
			b.log.Error().Err(err).
				Str(keyTopic, topic).
				Str(keyMsg, string(msg.Value)).
				Int(keyRetry, b.cfg.MaxRetries).
				Msg(msgCreateBrokerMessageError)
		}
	}
}
