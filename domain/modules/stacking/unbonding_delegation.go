package stacking

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog"
	"spacebox-writer/internal/configs"
	"time"

	"github.com/hexy-dev/spacebox/broker/model"
	"spacebox-writer/adapter/broker"
	"spacebox-writer/adapter/clickhouse"
)

type unbondingDelegation struct {
	db     *clickhouse.Clickhouse // interface with needed methods
	cancel context.CancelFunc
	ch     chan any
	log    *zerolog.Logger
}

type UnbondingDelegationStorage struct {
	CompletionTimestamp time.Time `json:"completion_timestamp"`
	Coin                string    `json:"coin"`
	DelegatorAddress    string    `json:"delegator_address"`
	ValidatorAddress    string    `json:"validator_oper_addr"`
	Height              int64     `json:"height"`
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
		}

		val2 := UnbondingDelegationStorage{
			CompletionTimestamp: val.CompletionTimestamp,
			Coin:                string(bytes),
			DelegatorAddress:    val.DelegatorAddress,
			ValidatorAddress:    val.ValidatorAddress,
			Height:              val.Height,
		}

		v.db.GetGormDB(ctx).Table("unbonding_delegation").Create(val2)

		// v.db.SaveValidator() // interface implementation in adapter
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
