package stacking

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"spacebox-writer/internal/configs"

	"github.com/hexy-dev/spacebox/broker/model"
	"spacebox-writer/adapter/broker"
	"spacebox-writer/adapter/clickhouse"
)

type validator struct {
	db     *clickhouse.Clickhouse // interface with needed methods
	cancel context.CancelFunc
	ch     chan any
	log    *zerolog.Logger
}

func (v *validator) handle(ctx context.Context) {
	for message := range v.ch {
		select {
		case <-ctx.Done():
			return
		default:
		}

		bytes, ok := message.([]byte)
		if !ok {
			v.log.Error().Bool("converted", ok).Msg("type error")
			continue
		}

		val := model.Validator{}
		if err := json.Unmarshal(bytes, &val); err != nil {
			v.log.Error().Err(err).Msg("unmarshall error")
			continue
		}

		var (
			count int64
			db    = v.db.GetGormDB(ctx)
		)

		if db.Table("validator").
			Where("consensus_address = ?", val.ConsensusAddress).
			Count(&count); count == 0 {
			db.Table("validator").
				Create(val)
		} else {
			fmt.Println(val.ConsensusAddress, "already exists. Skip.")
		}

	}
}

func (v *validator) subscribe(cfg configs.Config, db *clickhouse.Clickhouse, log *zerolog.Logger) error {
	log.Info().Str("consumer", "validator").Msg("start consumer")

	v.log = log
	v.ch = make(chan any, 10)
	v.db = db

	b := broker.New(cfg, v.ch)
	ctx, cancel := context.WithCancel(context.Background())
	v.cancel = cancel

	if err := b.Subscribe(ctx, "validator"); err != nil {
		return err
	}

	go v.handle(ctx)

	return nil
}

func (v *validator) stop() {
	close(v.ch)
	v.cancel()
}
