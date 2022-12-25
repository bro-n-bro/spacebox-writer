package stacking

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog"
	"spacebox-writer/internal/configs"

	"github.com/hexy-dev/spacebox/broker/model"
	"spacebox-writer/adapter/broker"
	"spacebox-writer/adapter/clickhouse"
)

type stakingParams struct {
	db     *clickhouse.Clickhouse // interface with needed methods
	cancel context.CancelFunc
	ch     chan any
	log    *zerolog.Logger
}

type StakingParamsStorage struct {
	Params string `json:"params"`
	Height int64  `json:"height"`
}

func (v *stakingParams) handle(ctx context.Context) {
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

		val := model.StakingParams{}
		if err := json.Unmarshal(bytes, &val); err != nil {
			v.log.Error().Err(err).Msg("unmarshall error")
			continue
		}

		bytes, err := json.Marshal(val.Params)
		if err != nil {
			v.log.Error().Err(err).Msg("marshall error")
		}

		val2 := StakingParamsStorage{
			Params: string(bytes),
			Height: val.Height,
		}

		v.db.GetGormDB(ctx).Table("staking_params").Create(val2)

		// v.db.SaveValidator() // interface implementation in adapter
	}
}

func (v *stakingParams) subscribe(cfg configs.Config, db *clickhouse.Clickhouse, log *zerolog.Logger) error {
	log.Info().Str("consumer", "staking_params").Msg("start consumer")

	v.log = log
	v.ch = make(chan any, 10)
	v.db = db

	b := broker.New(cfg, v.ch)
	ctx, cancel := context.WithCancel(context.Background())
	v.cancel = cancel

	if err := b.Subscribe(ctx, "staking_params"); err != nil {
		return err
	}

	go v.handle(ctx)

	return nil
}

func (v *stakingParams) stop() {
	close(v.ch)
	v.cancel()
}
