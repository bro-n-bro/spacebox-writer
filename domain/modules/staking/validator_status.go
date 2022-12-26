package staking

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"spacebox-writer/internal/configs"

	"github.com/hexy-dev/spacebox/broker/model"
	"spacebox-writer/adapter/broker"
	"spacebox-writer/adapter/clickhouse"

	"github.com/jinzhu/copier"
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

		var (
			updates model.ValidatorStatus
			getVal  model.ValidatorStatus
			db      = v.db.GetGormDB(ctx)
		)

		if err := db.Table("validator_status").
			Where("validator_address = ?", val.ValidatorAddress).
			First(&getVal).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err = db.Table("validator_status").Create(val).Error; err != nil {
					v.log.Error().Err(err).Msg("error of create")
					continue
				}
			} else {
				v.log.Error().Err(err).Msg("error of database")
				continue
			}

		} else if val.Height > getVal.Height {

			if err = copier.Copy(&val, &updates); err != nil {
				v.log.Error().Err(err).Msg("error of prepare update")
				continue
			}

			if err = db.Table("validator_status").
				Where("validator_address = ?", val.ValidatorAddress).
				Updates(&updates).Error; err != nil {
				v.log.Error().Err(err).Msg("error of update")
				continue
			}

		}

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
