package modules

import (
	"context"
	"sync"

	"github.com/rs/zerolog"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	bank2 "github.com/bro-n-bro/spacebox-writer/modules/bank"
	core2 "github.com/bro-n-bro/spacebox-writer/modules/core"
	distribution2 "github.com/bro-n-bro/spacebox-writer/modules/distribution"
	gov2 "github.com/bro-n-bro/spacebox-writer/modules/gov"
	mint2 "github.com/bro-n-bro/spacebox-writer/modules/mint"
	staking2 "github.com/bro-n-bro/spacebox-writer/modules/staking"
)

const (
	msgSubscribed = "topic: %v subscribed"

	keyModule = "module"
)

var (
	// moduleHandlers is a map of module name and topic handlers
	moduleHandlers = map[string][]topicHandler{
		"core": {
			{"transaction", core2.TransactionHandler},
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
			{"delegation_reward_message", distribution2.DelegationRewardMessageHandler},
			{"distribution_commission", distribution2.DistributionCommissionHandler},
		},
		"gov": {
			{"gov_params", gov2.GovParamsHandler},
			{"proposal", gov2.ProposalHandler},
			{"proposal_deposit_message", gov2.ProposalDepositMessageHandler},
			{"submit_proposal_message", gov2.SubmitProposalMessageHandler},
		},
		"mint": {
			{"mint_params", mint2.MintParamsHandler},
		},
		"staking": {
			{"delegation", staking2.DelegationHandler},
			{"delegation_message", staking2.DelegationMessageHandler},
			{"redelegation", staking2.RedelegationHandler},
			{"redelegation_message", staking2.RedelegationMessageHandler},
			{"staking_params", staking2.StakingParamsHandler},
			{"unbonding_delegation", staking2.UnbondingDelegationHandler},
			{"unbonding_delegation_message", staking2.UnbondingDelegationMessageHandler},
		},
	}
)

type (
	// Modules is a universal struct for all modules
	Modules struct {
		brk           rep.Broker
		st            rep.Storage
		log           *zerolog.Logger
		consumersWG   *sync.WaitGroup
		stopConsumers context.CancelFunc
		cfg           Config
	}

	// topicHandler is a struct for topic name and her handler
	topicHandler struct { //nolint:govet
		topicName string
		handler   func(ctx context.Context, msg [][]byte, db rep.Storage) error
	}
)

// New creates new instance of Modules
func New(cfg Config, st rep.Storage, log zerolog.Logger, brk rep.Broker) *Modules {
	return &Modules{
		log:         &log,
		cfg:         cfg,
		brk:         brk,
		st:          st,
		consumersWG: &sync.WaitGroup{},
	}
}

// Start starts all modules
func (m *Modules) Start(_ context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	m.stopConsumers = cancel

	// subscribe to all topics
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

// Stop stops all modules
func (m *Modules) Stop(ctx context.Context) error {
	m.stopConsumers()
	m.consumersWG.Wait()
	return nil
}
