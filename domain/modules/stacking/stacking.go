package stacking

import (
	"context"
	"spacebox-writer/adapter/broker"
	"spacebox-writer/adapter/clickhouse"
)

type H interface {
	subscribe(broker.Config, *clickhouse.Clickhouse) error
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
		bCfg    broker.Config
		storage *clickhouse.Clickhouse
	}
)

func New(bCfg broker.Config, cl *clickhouse.Clickhouse) *Module {
	return &Module{
		bCfg:    bCfg,
		storage: cl,
	}
}

func (m *Module) Subscribe() error {
	for _, handler := range modules {
		if err := handler.subscribe(m.bCfg, m.storage); err != nil {
			return err
		}
	}

	return nil
}
