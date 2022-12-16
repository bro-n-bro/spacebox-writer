package broker

import (
	"context"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog"
)

type Broker struct {
	log *zerolog.Logger
	con *kafka.Consumer
	cfg Config
	ch  chan any
}

func New(cfg Config, ch chan any) *Broker {
	lg := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("cmp", "broker").Logger()
	return &Broker{
		log: &lg,
		cfg: cfg,
		ch:  ch,
	}
}

func (brk *Broker) Start(ctx context.Context) (err error) {
	brk.con, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brk.cfg.Address,
		"group.id":          brk.cfg.GroupID,
		"auto.offset.reset": brk.cfg.AutoOffsetReset,
	})

	if err != nil {
		return err
	}

	return nil
}

func (brk *Broker) Subscribe(ctx context.Context, topic string) (err error) {
	brk.con, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brk.cfg.Address,
		"group.id":          brk.cfg.GroupID,
		"auto.offset.reset": brk.cfg.AutoOffsetReset,
	})

	if err != nil {
		return err
	}

	if err = brk.con.Subscribe(topic, nil); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			msg, err := brk.con.ReadMessage(-1)
			if err == nil {
				brk.log.Info().Msgf("msg: [%v] %s", msg.String(), msg.Value)
				brk.ch <- msg.Value

			} else {
				brk.log.
					Error().
					Err(err).
					Interface("msg", msg).
					Msg("consumer error")
			}

		}
	}()

	return nil
}

func (brk *Broker) Stop(ctx context.Context) (err error) {
	if err = brk.con.Close(); err != nil {
		return err
	}
	return nil
}
