package stacking

import (
	"context"
	"encoding/json"

	"spacebox-writer/adapter/broker"
	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"
	"spacebox-writer/internal/configs"

	"github.com/hexy-dev/spacebox/broker/model"
	"github.com/rs/zerolog"
)

type delegation struct {
	db     *clickhouse.Clickhouse // interface with needed methods
	cancel context.CancelFunc
	ch     chan any
	log    *zerolog.Logger
}

func (v *delegation) handle(ctx context.Context) {
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

		val := model.Delegation{}
		if err := json.Unmarshal(bytes, &val); err != nil {
			v.log.Error().Err(err).Msg("marshall error")
			continue
		}

		bytes, err := json.Marshal(val.Coin)
		if err != nil {
			v.log.Error().Err(err).Msg("marshall error")
			continue
		}

		val2 := storageModel.Delegation{
			OperatorAddress:  val.OperatorAddress,
			DelegatorAddress: val.DelegatorAddress,
			Coin:             string(bytes),
			Height:           val.Height,
		}

		v.db.GetGormDB(ctx).Table("delegation").Create(val2)

	}
}

func (v *delegation) subscribe(cfg configs.Config, db *clickhouse.Clickhouse, log *zerolog.Logger) error {
	log.Info().Str("consumer", "delegation").Msg("start consumer")

	v.log = log
	v.ch = make(chan any, 10)
	v.db = db

	b := broker.New(cfg, v.ch)
	ctx, cancel := context.WithCancel(context.Background())
	v.cancel = cancel

	if err := b.Subscribe(ctx, "delegation"); err != nil {
		return err
	}

	go v.handle(ctx)

	return nil
}

func (v *delegation) stop() {
	close(v.ch)
	v.cancel()
}