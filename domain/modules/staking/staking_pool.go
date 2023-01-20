package staking

import (
	"context"
	"encoding/json"

	"spacebox-writer/adapter/broker"
	"spacebox-writer/adapter/clickhouse"
	"spacebox-writer/internal/configs"

	"github.com/rs/zerolog"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type stakingPool struct {
	db     *clickhouse.Clickhouse // interface with needed methods
	cancel context.CancelFunc
	ch     chan any
	log    *zerolog.Logger
}

func (v *stakingPool) handle(ctx context.Context) {
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

		val := model.StakingPool{}
		if err := json.Unmarshal(bytes, &val); err != nil {
			v.log.Error().Err(err).Msg("unmarshall error")
			continue
		}

		var (
			count int64
			db    = v.db.GetGormDB(ctx)
		)

		if db.Table("staking_pool").
			Where("height = ?", val.Height).
			Count(&count); count != 0 {
			continue

		}

		if err := db.Table("staking_pool").Create(val).Error; err != nil {
			v.log.Error().Err(err).Msg("error of create")
			continue
		}

	}
}

func (v *stakingPool) subscribe(cfg configs.Config, db *clickhouse.Clickhouse, log *zerolog.Logger) error {
	log.Info().Str("consumer", "staking_pool").Msg("start consumer")

	v.log = log
	v.ch = make(chan any, 10)
	v.db = db

	b := broker.New(cfg, v.ch)
	ctx, cancel := context.WithCancel(context.Background())
	v.cancel = cancel

	if err := b.Subscribe(ctx, "staking_pool"); err != nil {
		return err
	}

	go v.handle(ctx)

	return nil
}

func (v *stakingPool) stop() {
	close(v.ch)
	v.cancel()
}
