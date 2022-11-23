package broker

import (
	"context"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog"
)

type Broker struct {
	log *zerolog.Logger
	con *kafka.Consumer
	cfg Config
}

func New(cfg Config) *Broker {
	lg := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("cmp", "broker").Logger()
	return &Broker{
		log: &lg,
		cfg: cfg,
	}
}

func (brk *Broker) Start(context.Context) (err error) {

	errCh := make(chan error)
	brk.log.Debug().Msgf("start broker consumer")
	go func() {

		brk.con, err = kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers": brk.cfg.Address,
			"group.id":          brk.cfg.GroupID,
			"auto.offset.reset": brk.cfg.AutoOffsetReset,
		})
		if err != nil {
			errCh <- err
		}

		if err = brk.con.SubscribeTopics(brk.cfg.Topics, nil); err != nil {
			errCh <- err
		}

		for {
			msg, err := brk.con.ReadMessage(-1)
			if err == nil {
				brk.log.Info().Msgf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			} else {
				brk.log.Error().Msgf("Consumer error: %v (%v)\n", err, msg)
			}
		}

	}()

	select {
	case err := <-errCh:
		return err
	case <-time.After(brk.cfg.StartTimeout):
		return nil
	}

}

func (brk *Broker) Stop(ctx context.Context) (err error) {
	if err = brk.con.Close(); err != nil {
		return err
	}
	return nil
}
