package broker

import (
	"context"
	"strconv"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"spacebox-writer/adapter/clickhouse"
)

func (b *Broker) Subscribe(
	ctx context.Context,
	wg *sync.WaitGroup,
	topic string,
	handler func(ctx context.Context, msg []byte, db *clickhouse.Clickhouse) error,
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
				b.log.Info().Str("topic", topic).Msg("stop read messages from topic")
				return
			default:
			}

			msg, err := consumer.ReadMessage(100)
			if msg == nil {
				continue
			}

			if err != nil {
				b.log.
					Fatal().
					Err(err).
					Str("msg", string(msg.Value)).
					Msg("read message error")
				return
			} else {
				b.log.Debug().Msgf("[%v]: %s", msg.String(), msg.Value)
			}

			if err := handler(ctx, msg.Value, b.clh); err != nil {
				headers := msg.Headers
				var retry int
				for _, header := range headers {
					if header.Key == "retry" {
						retry = bytesToInt(header.Value)
						break
					}
				}
				retry++

				b.log.
					Error().
					Err(err).
					Str("topic", topic).
					Int64("offset", int64(msg.TopicPartition.Offset)).
					Int64("partition", int64(msg.TopicPartition.Partition)).
					Int("retry", retry).
					Str("msg", string(msg.Value)).
					Msg("handle message error")

				// TODO: check max retries

				msg.Headers = append(msg.Headers, kafka.Header{
					Key:   "retry",
					Value: []byte(strconv.Itoa(retry)),
				})

				if err = b.produce(topic, msg.Value, msg.Headers); err != nil {
					b.log.
						Error().
						Err(err).
						Str("topic", topic).
						Int64("offset", int64(msg.TopicPartition.Offset)).
						Int64("partition", int64(msg.TopicPartition.Partition)).
						Str("msg", string(msg.Value)).
						Int("retry", retry).
						Msg("produce message error")
				}
			}

			_, err = consumer.CommitMessage(msg)
			if err != nil {
				b.log.
					Error().
					Err(err).
					Str("topic", topic).
					Msg("commit message error")
			}
		}
	}(wg)

	// go func() {
	//	for {
	//		select {
	//		case <-ctx.Done():
	//			b.log.Info().Str("topic", topic).Msg("stop read messages from topic")
	//			return
	//		default:
	//		}
	//
	//		ev := consumer.Poll(100)
	//
	//		if ev == nil {
	//			continue
	//		}
	//
	//		switch e := ev.(type) {
	//		case *kafka.Message:
	//			if err := handler(ctx, e.Value, b.clh); err != nil {
	//				b.log.
	//					Error().
	//					Err(err).
	//					Str("topic", topic).
	//					Str("msg", string(e.Value)).
	//					Int64("offset", int64(e.TopicPartition.Offset)).
	//					Msg("handle message error")
	//				continue
	//			}
	//			if _, err := consumer.StoreMessage(e); err != nil {
	//				b.log.
	//					Error().
	//					Err(err).
	//					Str("topic", topic).
	//					Str("msg", string(e.Value)).
	//					Msg("store message error")
	//				continue
	//			}
	//		}
	//	}
	// }()

	return nil
}
