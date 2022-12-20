package stacking

import (
	"context"
	"spacebox-writer/adapter/clickhouse"
	"spacebox-writer/internal/configs"
)

type H interface {
	subscribe(configs.Config, *clickhouse.Clickhouse) error
	handle(context.Context)
}

var (
	topics = []string{
		"validator",
		"validator_status",
		"validator_info",
	}

	modules = []H{&validator{}}
)

type (
	Module struct {
		cfg     configs.Config
		storage *clickhouse.Clickhouse
	}
)

func New(cfg configs.Config, cl *clickhouse.Clickhouse) *Module {
	return &Module{
		cfg:     cfg,
		storage: cl,
	}
}

func (m *Module) Subscribe() error {
	for _, handler := range modules {
		if err := handler.subscribe(m.cfg, m.storage); err != nil {
			return err
		}
	}

	return nil
}
