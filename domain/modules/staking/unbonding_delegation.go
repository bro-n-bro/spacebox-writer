package staking

import (
	"context"
	"encoding/json"

	"spacebox-writer/adapter/broker"
	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"
	"spacebox-writer/internal/configs"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type unbondingDelegation struct {
	db     *clickhouse.Clickhouse // interface with needed methods
	cancel context.CancelFunc
	ch     chan any
	log    *zerolog.Logger
}

func (v *unbondingDelegation) handle(ctx context.Context) {
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

		val := model.UnbondingDelegation{}
		if err := json.Unmarshal(bytes, &val); err != nil {
			v.log.Error().Err(err).Msg("unmarshall error")
			continue
		}

		bytes, err := json.Marshal(val.Coin)
		if err != nil {
			v.log.Error().Err(err).Msg("marshall error")
			continue
		}

		val2 := storageModel.UnbondingDelegation{
			CompletionTimestamp: val.CompletionTimestamp,
			Coin:                string(bytes),
			DelegatorAddress:    val.DelegatorAddress,
			ValidatorAddress:    val.ValidatorAddress,
			Height:              val.Height,
		}

		// if validator_address + delegator_address
		// not exists - create
		// if new height greater than height
		// in DB - update height

		var (
			updates storageModel.UnbondingDelegation
			getVal  storageModel.UnbondingDelegation
			db      = v.db.GetGormDB(ctx)
		)

		if err := db.Table("unbonding_delegation").
			Where("validator_address = ? AND delegator_address = ?",
				val2.ValidatorAddress,
				val2.DelegatorAddress,
			).First(&getVal).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err = db.Table("unbonding_delegation").Create(val2).Error; err != nil {
					v.log.Error().Err(err).Msg("error of create")
					continue
				}
			} else {
				v.log.Error().Err(err).Msg("error of database")
				continue
			}

		} else if val2.Height > getVal.Height {

			if err = copier.Copy(&val2, &updates); err != nil {
				v.log.Error().Err(err).Msg("error of prepare update")
				continue
			}

			if err = db.Table("unbonding_delegation").
				Where("validator_address = ? AND delegator_address = ?",
					val2.ValidatorAddress,
					val2.DelegatorAddress,
				).Updates(&updates).Error; err != nil {
				v.log.Error().Err(err).Msg("error of update")
				continue
			}

		}

		// v.db.GetGormDB(ctx).Table("unbonding_delegation").Create(val2)

	}
}

func (v *unbondingDelegation) subscribe(cfg configs.Config, db *clickhouse.Clickhouse, log *zerolog.Logger) error {
	log.Info().Str("consumer", "unbonding_delegation").Msg("start consumer")

	v.log = log
	v.ch = make(chan any, 10)
	v.db = db

	b := broker.New(cfg, v.ch)
	ctx, cancel := context.WithCancel(context.Background())
	v.cancel = cancel

	if err := b.Subscribe(ctx, "unbonding_delegation"); err != nil {
		return err
	}

	go v.handle(ctx)

	return nil
}

func (v *unbondingDelegation) stop() {
	close(v.ch)
	v.cancel()
}
