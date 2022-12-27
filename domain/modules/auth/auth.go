package auth

import (
	"spacebox-writer/adapter/clickhouse"
	"spacebox-writer/internal/configs"

	"github.com/rs/zerolog"
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
	//for _, handler := range consumers {
	//	if err := handler.subscribe(m.cfg, m.storage, m.log); err != nil {
	//		return err
	//	}
	//}
	return nil
}
