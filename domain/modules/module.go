package modules

import (
	"context"

	"github.com/rs/zerolog"

	"spacebox-writer/adapter/clickhouse"
	"spacebox-writer/domain/modules/auth"
	"spacebox-writer/domain/modules/staking"
	"spacebox-writer/internal/configs"
	"spacebox-writer/internal/rep"
)

type Modules struct {
	cfg configs.Config
	st  *clickhouse.Clickhouse
	log *zerolog.Logger
	b   rep.Broker

	stopConsumers context.CancelFunc
}

type topicHandler struct {
	topicName string
	handler   func(ctx context.Context, msg []byte, db *clickhouse.Clickhouse) error
}

var (
	moduleHandlers = map[string][]topicHandler{
		"auth": {{"account", auth.AccountHandler}},
		"staking": {
			{"validator", staking.ValidatorHandler},
			{"delegation", staking.DelegationHandler},
			{"delegation_message", staking.DelegationMessageHandler},
			{"redelegation", staking.RedelegationHandler},
			{"redelegation_message", staking.RedelegationMessageHandler},
			{"staking_params", staking.StakingParamsHandler},
			{"staking_pool", staking.StakingPoolHandler},
			{"unbonding_delegation", staking.UnbondingDelegationHandler},
			{"unbonding_delegation_message", staking.UnbondingDelegationMessageHandler},
			{"validator_info", staking.ValidatorInfoHandler},
			{"validator_status", staking.ValidatorStatusHandler},
		},
	}
)

func New(cfg configs.Config, s *clickhouse.Clickhouse, log zerolog.Logger, b rep.Broker) *Modules {
	return &Modules{
		cfg: cfg,
		st:  s, // TODO: use interface
		log: &log,
		b:   b,
	}
}

func (m *Modules) Start(_ context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	m.stopConsumers = cancel

	for _, moduleName := range m.cfg.Modules {
		if topicHandlers, ok := moduleHandlers[moduleName]; ok {
			for _, th := range topicHandlers {
				if err := m.b.Subscribe(ctx, th.topicName, th.handler); err != nil {
					return err
				}
				m.log.Info().Str("module", moduleName).Msgf("topic: %v subscribed", th.topicName)
			}
		}
	}

	return nil
}

func (m *Modules) Stop(ctx context.Context) error {
	m.stopConsumers()
	return nil
}
