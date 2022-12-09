package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog"
	"spacebox-writer/adapter/clickhouse"
)

type Broker struct {
	log *zerolog.Logger
	con *kafka.Consumer
	db  *clickhouse.Clickhouse
	cfg Config
}

func New(cfg Config, db *clickhouse.Clickhouse) *Broker {
	lg := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("cmp", "broker").Logger()
	return &Broker{
		log: &lg,
		cfg: cfg,
		db:  db,
	}
}

type (
	Supply struct {
		Height uint64      `json:"height"`
		Coins  interface{} `json:"coins"`
	}
)

func (brk *Broker) Start(ctx context.Context) (err error) {
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
			return
		}

		if err = brk.con.Subscribe("block-topic", nil); err != nil {
			errCh <- err
			return
		}

		for {
			msg, err := brk.con.ReadMessage(-1)
			if err == nil {
				// TODO:
				brk.log.Info().Msgf("msg: [%v] %s", msg.String(), msg.Value)

				sp := Supply{}

				if err = json.Unmarshal(msg.Value, &sp); err != nil {
					return
				}
				brk.db.GetGormDB(ctx).Table("supply").Create(&sp)

				fmt.Println(sp.Height, sp.Coins)

			} else {
				brk.log.
					Error().
					Err(err).
					Interface("msg", msg).
					Msg("consumer error")
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
