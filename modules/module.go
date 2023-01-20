package modules

import (
	"context"
	"sync"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	auth2 "github.com/bro-n-bro/spacebox-writer/modules/auth"
	bank2 "github.com/bro-n-bro/spacebox-writer/modules/bank"
	core2 "github.com/bro-n-bro/spacebox-writer/modules/core"
	distribution2 "github.com/bro-n-bro/spacebox-writer/modules/distribution"
	gov2 "github.com/bro-n-bro/spacebox-writer/modules/gov"
	mint2 "github.com/bro-n-bro/spacebox-writer/modules/mint"
	staking2 "github.com/bro-n-bro/spacebox-writer/modules/staking"
	"github.com/rs/zerolog"
)

const (
	msgSubscribed = "topic: %v subscribed"

	keyModule = "module"
)

var (
	moduleHandlers = map[string][]topicHandler{
		"core": {
			{"block", core2.BlockHandler},
			{"message", core2.MessageHandler},
			{"transaction", core2.TransactionHandler},
		},
		"auth": {
			{"account", auth2.AccountHandler},
		},
		"bank": {
			{"supply", bank2.SupplyHandler},
			{"account_balance", bank2.AccountBalanceHandler},
			{"multisend_message", bank2.MultiSendMessageHandler},
			{"send_message", bank2.SendMessageHandler},
		},
		"distribution": {
			{"distribution_params", distribution2.DistributionParamsHandler},
			{"community_pool", distribution2.CommunityPoolHandler},
			{"validator_commission", distribution2.ValidatorCommissionHandler},
			{"delegation_reward_message", distribution2.DelegationRewardMessageHandler},
		},
		"gov": {
			{"gov_params", gov2.GovParamsHandler},
			{"proposal", gov2.ProposalHandler},
			{"proposal_deposit_message", gov2.ProposalDepositMessageHandler},
			{"proposal_tally_result", gov2.ProposalTallyResultHandler},
			{"proposal_vote_message", gov2.ProposalVoteMessageHandler},
		},
		"mint": {
			{"mint_params", mint2.MintParamsHandler},
			{"annual_provision", mint2.AnnualProvisionHandler},
		},
		"staking": {
			{"validator", staking2.ValidatorHandler},
			{"delegation", staking2.DelegationHandler},
			{"delegation_message", staking2.DelegationMessageHandler},
			{"redelegation", staking2.RedelegationHandler},
			{"redelegation_message", staking2.RedelegationMessageHandler},
			{"staking_params", staking2.StakingParamsHandler},
			{"staking_pool", staking2.StakingPoolHandler},
			{"unbonding_delegation", staking2.UnbondingDelegationHandler},
			{"unbonding_delegation_message", staking2.UnbondingDelegationMessageHandler},
			{"validator_info", staking2.ValidatorInfoHandler},
			{"validator_status", staking2.ValidatorStatusHandler},
			{"validator_description", staking2.ValidatorDescriptionHandler},
		},
	}
)

type (
	Modules struct {
		brk           rep.Broker
		st            rep.Storage
		log           *zerolog.Logger
		consumersWG   *sync.WaitGroup
		stopConsumers context.CancelFunc
		cfg           Config
	}

	topicHandler struct { //nolint:govet
		topicName string
		handler   func(ctx context.Context, msg []byte, db rep.Storage) error
	}
)

func New(cfg Config, st rep.Storage, log zerolog.Logger, brk rep.Broker) *Modules {
	return &Modules{
		log:         &log,
		cfg:         cfg,
		brk:         brk,
		st:          st,
		consumersWG: &sync.WaitGroup{},
	}
}

func (m *Modules) Start(_ context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	m.stopConsumers = cancel

	for _, moduleName := range m.cfg.Modules {
		if topicHandlers, ok := moduleHandlers[moduleName]; ok {
			for _, th := range topicHandlers {
				m.consumersWG.Add(1)
				if err := m.brk.Subscribe(ctx, m.consumersWG, th.topicName, th.handler); err != nil {
					return err
				}
				m.log.Info().Str(keyModule, moduleName).Msgf(msgSubscribed, th.topicName)
			}
		}
	}

	return nil
}

func (m *Modules) Stop(ctx context.Context) error {
	m.stopConsumers()
	m.consumersWG.Wait()
	return nil
}
