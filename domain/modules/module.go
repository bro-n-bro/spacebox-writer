package modules

import (
	"context"
	"sync"

	"github.com/rs/zerolog"

	"spacebox-writer/adapter/clickhouse"
	"spacebox-writer/domain/modules/auth"
	"spacebox-writer/domain/modules/bank"
	"spacebox-writer/domain/modules/core"
	"spacebox-writer/domain/modules/distribution"
	"spacebox-writer/domain/modules/gov"
	"spacebox-writer/domain/modules/mint"
	"spacebox-writer/domain/modules/staking"
	"spacebox-writer/internal/rep"
)

type Modules struct {
	b             rep.Broker
	st            *clickhouse.Clickhouse
	log           *zerolog.Logger
	consumersWg   *sync.WaitGroup
	stopConsumers context.CancelFunc
	cfg           Config
}

type topicHandler struct { // nolint:govet
	topicName string
	handler   func(ctx context.Context, msg []byte, db *clickhouse.Clickhouse) error
}

var (
	moduleHandlers = map[string][]topicHandler{
		"core": {
			{"block", core.BlockHandler},
			{"message", core.MessageHandler},
			{"transaction", core.TransactionHandler},
		},
		"auth": {
			{"account", auth.AccountHandler},
		},
		"bank": {
			{"account_balance", bank.AccountBalanceHandler},
			{"multisend_message", bank.MultiSendMessageHandler},
			{"send_message", bank.SendMessageHandler},
			{"supply", bank.SupplyHandler},
		},
		"distribution": {
			{"distribution_params", distribution.DistributionParamsHandler},
			{"community_pool", distribution.CommunityPoolHandler},
			{"validator_commission", distribution.ValidatorCommissionHandler},
			{"delegation_reward", distribution.DelegationRewardHandler},
			{"delegation_reward_message", distribution.DelegationRewardMessageHandler},
		},
		"gov": {
			{"gov_params", gov.GovParamsHandler},
			{"proposal", gov.ProposalHandler},
			{"proposal_deposit", gov.ProposalDepositHandler},
			{"proposal_deposit_message", gov.ProposalDepositMessageHandler},
			{"proposal_tally_result", gov.ProposalTallyResultHandler},
			{"proposal_vote_message", gov.ProposalVoteMessageHandler},
		},
		"mint": {
			{"mint_params", mint.MintParamsHandler},
			{"inflation", mint.InflationHandler},
			{"annual_provision", mint.AnnualProvisionHandler},
		},
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

func New(cfg Config, s *clickhouse.Clickhouse, log zerolog.Logger, b rep.Broker) *Modules {
	return &Modules{
		cfg:         cfg,
		st:          s, // TODO: use interface
		log:         &log,
		b:           b,
		consumersWg: &sync.WaitGroup{},
	}
}

func (m *Modules) Start(_ context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	m.stopConsumers = cancel

	for _, moduleName := range m.cfg.Modules {
		if topicHandlers, ok := moduleHandlers[moduleName]; ok {
			for _, th := range topicHandlers {
				m.consumersWg.Add(1)
				if err := m.b.Subscribe(ctx, m.consumersWg, th.topicName, th.handler); err != nil {
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
	m.consumersWg.Wait()
	return nil
}
