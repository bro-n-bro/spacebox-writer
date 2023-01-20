package broker

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"

	"github.com/hexy-dev/spacebox-writer/adapter/mongo/model"
	"github.com/hexy-dev/spacebox-writer/internal/rep"
)

func (b *Broker) Subscribe(
	ctx context.Context,
	wg *sync.WaitGroup,
	topic string,
	handler func(ctx context.Context, msg []byte, db rep.Storage) error,
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

		for {
			select {
			case <-ctx.Done():
				b.log.Info().Str(keyTopic, topic).Msg("stop read messages from topic")
				return
			default:
			}

			msg, err := consumer.ReadMessage(100)
			if msg == nil {
				continue
			}

			if err != nil {
				b.log.Fatal().Err(err).
					Str(keyMsg, string(msg.Value)).
					Msg("read message error")
				return
			} else {
				b.log.Debug().Msgf("[%v]: %s", msg.String(), msg.Value)
			}

			hndlErr := handler(ctx, msg.Value, b.st)

			// call handler and process error if needed
			if err = b.handleError(ctx, hndlErr, msg); err != nil {
				b.log.Error().Err(err).Msg("smth went wrong with handle error")
			}

			if _, err = consumer.CommitMessage(msg); err != nil {
				b.log.Error().Err(err).
					Str(keyTopic, topic).
					Msg("commit message error")
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
func (b *Broker) handleError(ctx context.Context, messageHandlerError error, msg *kafka.Message) error {
	headers := msg.Headers
	topic := *msg.TopicPartition.Topic

	// find unique message_id from kafka header
	messageID := string(findValueFromHeaders(keyMessageID, headers))

	if messageID == "" {
		b.log.Debug().Str(keyTopic, topic).Msg("empty message id. generate new")
		messageID = uuid.New().String()
	}

	// check if error message already exists in mongo
	exists, err := b.m.HasBrokerMessage(ctx, messageID)
	if err != nil {
		return err
	}

	if messageHandlerError == nil {
		if exists {
			b.log.Debug().
				Str(keyTopic, topic).
				Str(keyMessageID, messageID).
				Msg("handle message successful. delete errors in storage")

			if err = b.m.DeleteBrokerMessage(ctx, messageID); err != nil {
				b.log.Warn().Err(err).
					Str(keyMessageID, messageID).
					Msg("DeleteBrokerMessage error. But handle message successful")
			}
		}

		return nil // handle message successful
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
		Msg("handle message error")

	if retry > b.cfg.MaxRetries { // retry limit exceeded
		// TODO: any notifications?
		b.log.Error().
			Str(keyTopic, topic).
			Str(keyMessageID, messageID).
			Str(keyMsg, string(msg.Value)).
			Int(keyRetry, retry).
			Msg("retry limit exceeded!!!")

		return nil
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
	if err = b.produce(topic, msg.Value, msg.Headers); err != nil {
		// FIXME: what need to do?
		b.log.Error().Err(err).
			Str(keyTopic, topic).
			Str(keyMessageID, messageID).
			Int64(keyOffset, int64(msg.TopicPartition.Offset)).
			Int64(keyPartition, int64(msg.TopicPartition.Partition)).
			Str(keyMsg, string(msg.Value)).
			Int(keyRetry, retry).
			Msg("produce message error")

		return nil // log above we dont need an error here
	}

	if exists { // error message exists in mongo. just increase an attempts
		if err = b.m.UpdateBrokerMessage(ctx, &model.BrokerMessage{
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
				Msg("UpdateBrokerMessage error")
		}
	} else {
		if err = b.m.CreateBrokerMessage(ctx, &model.BrokerMessage{
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
				Msg("CreateBrokerMessage error")
		}
	}

	return nil
}
