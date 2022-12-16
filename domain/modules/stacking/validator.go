package stacking

import (
	"context"
	"encoding/json"

	"github.com/hexy-dev/spacebox/broker/model"
	"spacebox-writer/adapter/broker"
	"spacebox-writer/adapter/clickhouse"
)

type validator struct {
	db *clickhouse.Clickhouse // interface with needed methods

	cancel context.CancelFunc

	ch chan any
}

func (v *validator) handle(ctx context.Context) {
	for message := range v.ch {
		select {
		case <-ctx.Done():
			return
		default:
		}

		bytes, ok := message.([]byte)
		if !ok {
			continue
		}
		val := model.Validator{}
		if err := json.Unmarshal(bytes, &val); err != nil {
			// logl
			continue
		}

		// business logic

		// v.db.SaveValidator() // interface implementation in adapter
	}
}

func (v *validator) subscribe(cfg broker.Config, db *clickhouse.Clickhouse) error {
	v.ch = make(chan any, 10)
	v.db = db

	b := broker.New(cfg, v.ch)
	ctx, cancel := context.WithCancel(context.Background())
	v.cancel = cancel

	if err := b.Subscribe(ctx, "validator"); err != nil {
		return err
	}

	go v.handle(ctx)

	return nil
}

func (v *validator) stop() {
	close(v.ch)
	v.cancel()
}
