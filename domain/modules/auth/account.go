package auth

import (
	"context"
	"encoding/json"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"spacebox-writer/adapter/broker"
	"spacebox-writer/adapter/clickhouse"
	"spacebox-writer/internal/configs"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type account struct {
	db     *clickhouse.Clickhouse // interface with needed methods
	cancel context.CancelFunc
	ch     chan any
	log    *zerolog.Logger
}

func (v *account) handle(ctx context.Context) {
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

		val := model.Account{}
		if err := json.Unmarshal(bytes, &val); err != nil {
			v.log.Error().Err(err).Msg("unmarshall error")
			continue
		}

		var (
			db      = v.db.GetGormDB(ctx)
			updates model.Account
			getVal  model.Account
		)

		if err := db.Table("account").
			Where("address = ?", val.Address).
			First(&getVal).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				if err = db.Table("account").Create(val).Error; err != nil {
					v.log.Error().Err(err).Msg("error of create")
					continue
				}
			} else {
				v.log.Error().Err(err).Msg("error of database")
				continue
			}

		}

		if val.Height < getVal.Height {

			if err := copier.Copy(&val, &updates); err != nil {
				v.log.Error().Err(err).Msg("error of prepare update")
				continue
			}

			if err := db.Table("account").
				Where("address = ?", val.Height).
				Updates(&updates).Error; err != nil {
				v.log.Error().Err(err).Msg("error of update")
				continue
			}

		}

	}
}

func (v *account) subscribe(cfg configs.Config, db *clickhouse.Clickhouse, log *zerolog.Logger) error {
	log.Info().Str("consumer", "account").Msg("start consumer")

	v.log = log
	v.ch = make(chan any, 10)
	v.db = db

	b := broker.New(cfg, v.ch)
	ctx, cancel := context.WithCancel(context.Background())
	v.cancel = cancel

	if err := b.Subscribe(ctx, "account"); err != nil {
		return err
	}

	go v.handle(ctx)

	return nil
}

func (v *account) stop() {
	close(v.ch)
	v.cancel()
}
