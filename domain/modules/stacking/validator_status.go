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

type validatorStatus struct {
	db     *clickhouse.Clickhouse // interface with needed methods
	cancel context.CancelFunc
	ch     chan any
	log    *zerolog.Logger
}

func (v *validatorStatus) handle(ctx context.Context) {
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

		val := model.ValidatorStatus{}
		if err := json.Unmarshal(bytes, &val); err != nil {
			v.log.Error().Err(err).Msg("unmarshall error")
			continue
		}

		v.db.GetGormDB(ctx).Table("validator_status").Create(val)

		// v.db.SaveValidator() // interface implementation in adapter
	}
}

func (v *validatorStatus) subscribe(cfg configs.Config, db *clickhouse.Clickhouse, log *zerolog.Logger) error {
	log.Info().Str("consumer", "validator_status").Msg("start consumer")

	v.log = log
	v.ch = make(chan any, 10)
	v.db = db

	b := broker.New(cfg, v.ch)
	ctx, cancel := context.WithCancel(context.Background())
	v.cancel = cancel

	if err := b.Subscribe(ctx, "validator_status"); err != nil {
		return err
	}

	go v.handle(ctx)

	return nil
}

func (v *validatorStatus) stop() {
	close(v.ch)
	v.cancel()
}
