package stacking

import (
	"context"
	"github.com/rs/zerolog"
	"spacebox-writer/adapter/clickhouse"
	"spacebox-writer/internal/configs"
)

type H interface {
	subscribe(configs.Config, *clickhouse.Clickhouse, *zerolog.Logger) error
	handle(context.Context)
}

var (
	modules = []H{
		&validator{},
		&validatorStatus{},
		&validatorInfo{},
		&stakingParams{},
		&stakingPool{},
		&redelegation{},
		&redelegationMessage{},
		&unbondingDelegation{},
		&unbondingDelegationMessage{},
		&delegation{},
		&delegationMessage{},
	}
)

type (
	Module struct {
		cfg     configs.Config
		storage *clickhouse.Clickhouse
		log     *zerolog.Logger
	}
)

func New(cfg configs.Config, cl *clickhouse.Clickhouse, log *zerolog.Logger) *Module {
	return &Module{
		cfg:     cfg,
		storage: cl,
		log:     log,
	}
}

func (m *Module) Subscribe() error {
	for _, handler := range modules {
		if err := handler.subscribe(m.cfg, m.storage, m.log); err != nil {
			return err
		}
	}

	return nil
}
