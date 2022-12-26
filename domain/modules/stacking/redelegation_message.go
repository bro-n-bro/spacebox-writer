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

type redelegationMessage struct {
	db     *clickhouse.Clickhouse // interface with needed methods
	cancel context.CancelFunc
	ch     chan any
	log    *zerolog.Logger
}

func (v *redelegationMessage) handle(ctx context.Context) {
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

		val := model.RedelegationMessage{}
		if err := json.Unmarshal(bytes, &val); err != nil {
			v.log.Error().Err(err).Msg("unmarshall error")
			continue
		}

		bytes, err := json.Marshal(val.Coin)
		if err != nil {
			v.log.Error().Err(err).Msg("marshall error")
			continue
		}

		val2 := storageModel.RedelegationMessage{
			CompletionTime:      val.CompletionTime,
			Coin:                string(bytes),
			DelegatorAddress:    val.DelegatorAddress,
			SrcValidatorAddress: val.SrcValidatorAddress,
			DstValidatorAddress: val.DstValidatorAddress,
			Height:              val.Height,
			TxHash:              val.TxHash,
		}

		var (
			count int64
			db    = v.db.GetGormDB(ctx)
		)

		if db.Table("redelegation_message").
			Where("height = ? AND src_validator_address = ? AND delegator_address = ? AND dst_validator_address = ?",
				val2.Height,
				val2.SrcValidatorAddress,
				val2.DelegatorAddress,
				val2.DstValidatorAddress,
			).Count(&count); count != 0 {
			continue

		}

		if err = db.Table("redelegation_message").Create(val2).Error; err != nil {
			v.log.Error().Err(err).Msg("error of create")
			continue
		}

	}
}

func (v *redelegationMessage) subscribe(cfg configs.Config, db *clickhouse.Clickhouse, log *zerolog.Logger) error {
	log.Info().Str("consumer", "redelegation_message").Msg("start consumer")

	v.log = log
	v.ch = make(chan any, 10)
	v.db = db

	b := broker.New(cfg, v.ch)
	ctx, cancel := context.WithCancel(context.Background())
	v.cancel = cancel

	if err := b.Subscribe(ctx, "redelegation_message"); err != nil {
		return err
	}

	go v.handle(ctx)

	return nil
}

func (v *redelegationMessage) stop() {
	close(v.ch)
	v.cancel()
}
