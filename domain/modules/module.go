package modules

import (
	"context"
	"github.com/rs/zerolog"
	"spacebox-writer/adapter/clickhouse"
	"spacebox-writer/domain/modules/staking"
	"spacebox-writer/internal/configs"
)

type Modules struct {
	cfg configs.Config
	st  *clickhouse.Clickhouse
	log *zerolog.Logger
}

type subscriber interface {
	Subscribe() error
}

func New(cfg configs.Config, s *clickhouse.Clickhouse, log *zerolog.Logger) *Modules {
	return &Modules{
		cfg: cfg,
		st:  s,
		log: log,
	}
}

func (m *Modules) Start(ctx context.Context) error {
	activeModules := make([]subscriber, 0)
	for _, moduleName := range m.cfg.Modules {
		m.log.Info().Str("module", moduleName).Msg("start")
		switch moduleName {
		case "staking":
			activeModules = append(activeModules, staking.New(m.cfg, m.st, m.log))
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
