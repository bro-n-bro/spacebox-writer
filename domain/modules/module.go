package modules

import (
	"context"
	"spacebox-writer/adapter/clickhouse"
	"spacebox-writer/domain/modules/stacking"
	"spacebox-writer/internal/configs"
)

type Modules struct {
	cfg configs.Config
	st  *clickhouse.Clickhouse
}

type subscriber interface {
	Subscribe() error
}

func New(cfg configs.Config, s *clickhouse.Clickhouse) *Modules {
	return &Modules{
		cfg: cfg,
		st:  s,
	}
}

func (m *Modules) Start(ctx context.Context) error {
	activeModules := make([]subscriber, 0)
	for _, moduleName := range m.cfg.Modules {
		switch moduleName {
		case "stacking":
			activeModules = append(activeModules, stacking.New(m.cfg, m.st))
		}
	}

	for _, am := range activeModules {
		if err := am.Subscribe(); err != nil {
			return err
		}
	}

	return nil
}

func (m *Modules) Stop(ctx context.Context) error {
	return nil

}
