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

type redelegation struct {
	db     *clickhouse.Clickhouse // interface with needed methods
	cancel context.CancelFunc
	ch     chan any
	log    *zerolog.Logger
}

type RedelegationStorage struct {
	CompletionTime      time.Time `json:"completion_time"`
	Coin                string    `json:"coin"`
	DelegatorAddress    string    `json:"delegator_address"`
	SrcValidatorAddress string    `json:"src_validator"`
	DstValidatorAddress string    `json:"dst_validator"`
	Height              int64     `json:"height"`
}

func (v *redelegation) handle(ctx context.Context) {
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

		val := model.Redelegation{}
		if err := json.Unmarshal(bytes, &val); err != nil {
			v.log.Error().Err(err).Msg("unmarshall error")
			continue
		}

		bytes, err := json.Marshal(val.Coin)
		if err != nil {
			v.log.Error().Err(err).Msg("marshall error")
		}

		val2 := RedelegationStorage{
			CompletionTime:      val.CompletionTime,
			Coin:                string(bytes),
			DelegatorAddress:    val.DelegatorAddress,
			SrcValidatorAddress: val.SrcValidatorAddress,
			DstValidatorAddress: val.DstValidatorAddress,
			Height:              val.Height,
		}

		v.db.GetGormDB(ctx).Table("redelegation").Create(val2)

		// v.db.SaveValidator() // interface implementation in adapter
	}
}

func (v *redelegation) subscribe(cfg configs.Config, db *clickhouse.Clickhouse, log *zerolog.Logger) error {
	log.Info().Str("consumer", "redelegation").Msg("start consumer")

	v.log = log
	v.ch = make(chan any, 10)
	v.db = db

	b := broker.New(cfg, v.ch)
	ctx, cancel := context.WithCancel(context.Background())
	v.cancel = cancel

	if err := b.Subscribe(ctx, "redelegation"); err != nil {
		return err
	}

	go v.handle(ctx)

	return nil
}

func (v *redelegation) stop() {
	close(v.ch)
	v.cancel()
}
