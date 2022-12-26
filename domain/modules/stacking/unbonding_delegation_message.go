package stacking

import (
	"context"
	"encoding/json"
	"github.com/hexy-dev/spacebox/broker/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"spacebox-writer/adapter/broker"
	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"
	"spacebox-writer/internal/configs"
)

type unbondingDelegationMessage struct {
	db     *clickhouse.Clickhouse // interface with needed methods
	cancel context.CancelFunc
	ch     chan any
	log    *zerolog.Logger
}

func (v *unbondingDelegationMessage) handle(ctx context.Context) {
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

		val := model.UnbondingDelegationMessage{}
		if err := json.Unmarshal(bytes, &val); err != nil {
			v.log.Error().Err(err).Msg("unmarshall error")
			continue
		}

		bytes, err := json.Marshal(val.Coin)
		if err != nil {
			v.log.Error().Err(err).Msg("marshall error")
			continue
		}

		val2 := storageModel.UnbondingDelegationMessage{
			CompletionTimestamp: val.CompletionTimestamp,
			Coin:                string(bytes),
			DelegatorAddress:    val.DelegatorAddress,
			ValidatorAddress:    val.ValidatorAddress,
			Height:              val.Height,
			TxHash:              val.TxHash,
		}

		var (
			getVal storageModel.UnbondingDelegationMessage
			db     = v.db.GetGormDB(ctx)
		)

		if err := db.Table("unbonding_delegation_message").
			Where("validator_address = ? AND delegator_address = ? AND height = ?",
				val2.ValidatorAddress,
				val2.DelegatorAddress,
				val2.Height,
			).First(&getVal).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err = db.Table("unbonding_delegation_message").Create(val2).Error; err != nil {
					v.log.Error().Err(err).Msg("error of create")
					continue
				}
			} else {
				v.log.Error().Err(err).Msg("error of database")
				continue
			}

		}

		// v.db.GetGormDB(ctx).Table("unbonding_delegation_message").Create(val2)

	}
}

func (v *unbondingDelegationMessage) subscribe(cfg configs.Config, db *clickhouse.Clickhouse, log *zerolog.Logger) error {
	log.Info().Str("consumer", "unbonding_delegation_message").Msg("start consumer")

	v.log = log
	v.ch = make(chan any, 10)
	v.db = db

	b := broker.New(cfg, v.ch)
	ctx, cancel := context.WithCancel(context.Background())
	v.cancel = cancel

	if err := b.Subscribe(ctx, "unbonding_delegation_message"); err != nil {
		return err
	}

	go v.handle(ctx)

	return nil
}

func (v *unbondingDelegationMessage) stop() {
	close(v.ch)
	v.cancel()
}
