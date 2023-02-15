package broker

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/bro-n-bro/spacebox-writer/adapter/mongo/model"
	"github.com/bro-n-bro/spacebox-writer/internal/rep"
)

const (
	keyMsg           = "msg"
	keyPartition     = "partition"
	keyOffset        = "offset"
	keyTopic         = "topic"
	keyRetry         = "retry"
	keyMessageID     = "message_id"
	keyMessagesCount = "messages_count"

	msgStopReadingMessagesFromTopic = "stop reading messages from topic"
	msgEmptyMessageID               = "empty message id. Generate new id"
	msgCreateBrokerMessageError     = "CreateBrokerMessage error"
	msgUpdateBrokerMessageError     = "UpdateBrokerMessage error"
	msgDeleteBrokerMessageError     = "DeleteBrokerMessage error. But handle message successful"
	msgProduceMessageError          = "produce message error"
	msgReadMessageError             = "read message error"
	msgCommitMessageError           = "commit message error"
	msgRetryLimitExceeded           = "retry limit exceeded!"
	msgHandleMessageError           = "handle message error"
	msgHandleMessageSuccess         = "handle message successful. delete errors in storage"

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
		// "max.poll.interval.ms":     600000, // default 300000
		// "session.timeout.ms":       90000,  // default 45000
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
// Writes to storage any info about the error if message retries from broker <= .env MAX_RETRIES
//
// Do nothing if the message handling func does not return an error and this is the first message from the broker
// with the unique message_id
func (b *Broker) handleError(ctx context.Context, messageHandlerError error, msgs []*kafka.Message) {
	if len(msgs) == 0 {
		return
	}

	if messageHandlerError == nil {
		b.removeProcessedMessages(ctx, msgs)
		return // handle message successful
	}

	// send to broker each message to process it again
	for _, msg := range msgs {
		headers := msg.Headers
		topic := *msg.TopicPartition.Topic

		// find unique message_id from kafka header
		var (
			isFirstAttempt bool
			messageID      = string(findValueFromHeaders(keyMessageID, headers))
		)

		if messageID == "" {
			b.log.Debug().Str(keyTopic, topic).Msg(msgEmptyMessageID)
			messageID = uuid.New().String()
			isFirstAttempt = true
		}

		// handle occurred error
		if b.cfg.MetricsEnabled {
			b.metrics.errorsCounter.With(prometheus.Labels{keyTopic: topic}).Inc()
		}

		// find how much we already tried to handle this message
		retry := bytesToInt(findValueFromHeaders(keyRetry, headers))
		retry++

		// got error of handling message from the broker
		b.log.Error().Err(messageHandlerError).
			Str(keyTopic, topic).
			Int64(keyOffset, int64(msg.TopicPartition.Offset)).
			Int64(keyPartition, int64(msg.TopicPartition.Partition)).
			Int(keyRetry, retry).
			Str(keyMsg, string(msg.Value)).
			Msg(msgHandleMessageError)

		if retry > b.cfg.MaxRetries { // retry limit exceeded
			// TODO: any notifications?
			b.log.Error().
				Str(keyTopic, topic).
				Str(keyMessageID, messageID).
				Str(keyMsg, string(msg.Value)).
				Int(keyRetry, retry).
				Msg(msgRetryLimitExceeded)

			if b.cfg.MetricsEnabled {
				b.metrics.limitExceededCounter.With(prometheus.Labels{keyTopic: topic}).Inc()
			}

			continue
		}

		// add to the end of the queue again with new retry value
		// FIXME: what about another headers?
		msg.Headers = []kafka.Header{
			{
				Key:   keyRetry,
				Value: []byte(strconv.Itoa(retry)),
			},
			{
				Key:   keyMessageID,
				Value: []byte(messageID),
			},
		}

		// produce the message at the end of the broker`s queue
		if err := b.produce(topic, msg.Value, msg.Headers); err != nil {
			// FIXME: what need to do?
			b.log.Error().Err(err).
				Str(keyTopic, topic).
				Str(keyMessageID, messageID).
				Int64(keyOffset, int64(msg.TopicPartition.Offset)).
				Int64(keyPartition, int64(msg.TopicPartition.Partition)).
				Str(keyMsg, string(msg.Value)).
				Int(keyRetry, retry).
				Msg(msgProduceMessageError)

			continue
		}

		if !isFirstAttempt { // error message exists in mongo. just increase an attempts
			if err := b.m.UpdateBrokerMessage(ctx, &model.BrokerMessage{
				ID:               messageID,
				LastErrorMessage: messageHandlerError.Error(),
				Topic:            topic,
				Attempts:         retry,
				Data:             string(msg.Value),
			}); err != nil {
				b.log.Error().Err(err).
					Str(keyTopic, topic).
					Str(keyMessageID, messageID).
					Str(keyMsg, string(msg.Value)).
					Int(keyRetry, retry).
					Msg(msgUpdateBrokerMessageError)
			}
		} else {
			if err := b.m.CreateBrokerMessage(ctx, &model.BrokerMessage{
				ID:               messageID,
				LastErrorMessage: messageHandlerError.Error(),
				Topic:            topic,
				Data:             string(msg.Value),
				Attempts:         retry,
				Created:          time.Now(),
			}); err != nil {
				b.log.Error().Err(err).
					Str(keyTopic, topic).
					Str(keyMessageID, messageID).
					Str(keyMsg, string(msg.Value)).
					Int(keyMsg, retry).
					Msg(msgCreateBrokerMessageError)
			}
		}
	}
}

func (b *Broker) removeProcessedMessages(ctx context.Context, msgs []*kafka.Message) {
	ids := make([]string, 0)

	for _, msg := range msgs {
		if messageID := string(findValueFromHeaders(keyMessageID, msg.Headers)); messageID != "" {
			b.log.Debug().
				Str(keyTopic, *msg.TopicPartition.Topic).
				Str(keyMessageID, messageID).
				Msg(msgHandleMessageSuccess)

			ids = append(ids, messageID)
		}
	}

	if len(ids) > 0 {
		if err := b.m.DeleteBrokerMessages(ctx, ids); err != nil {
			b.log.Warn().Err(err).
				Int(keyMessagesCount, len(ids)).
				Msg(msgDeleteBrokerMessageError)
		}
	}
}
