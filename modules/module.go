package modules

import (
	"context"
	"sync"

	"github.com/rs/zerolog"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/authz"
	"github.com/bro-n-bro/spacebox-writer/modules/bank"
	"github.com/bro-n-bro/spacebox-writer/modules/core"
	"github.com/bro-n-bro/spacebox-writer/modules/distribution"
	"github.com/bro-n-bro/spacebox-writer/modules/gov"
	"github.com/bro-n-bro/spacebox-writer/modules/mint"
	"github.com/bro-n-bro/spacebox-writer/modules/slashing"
	"github.com/bro-n-bro/spacebox-writer/modules/staking"
)

const (
	msgSubscribed = "topic: %v subscribed"

	keyModule = "module"
)

var (
	// moduleHandlers is a map of module name and topic handlers
	moduleHandlers = map[string][]topicHandler{
		"core": {
			{"transaction", core.TransactionHandler},
		},
		"bank": {
			{"supply", bank.SupplyHandler},
			{"account_balance", bank.AccountBalanceHandler},
			{"multisend_message", bank.MultiSendMessageHandler},
			{"send_message", bank.SendMessageHandler},
		},
		"distribution": {
			{"distribution_params", distribution.DistributionParamsHandler},
			{"community_pool", distribution.CommunityPoolHandler},
			{"delegation_reward_message", distribution.DelegationRewardMessageHandler},
			{"proposer_reward", distribution.ProposerRewardHandler},
			{"distribution_reward", distribution.DistributionRewardHandler},
			{
				"withdraw_validator_commission_message",
				distribution.WithdrawValidatorCommissionMessageHandler,
			},
		},
		"gov": {
			{"gov_params", gov.GovParamsHandler},
			{"proposal", gov.ProposalHandler},
			{"proposal_deposit_message", gov.ProposalDepositMessageHandler},
			{"submit_proposal_message", gov.SubmitProposalMessageHandler},
			{"vote_weighted_message", gov.VoteWeightedMessageHandler},
		},
		"mint": {
			{"mint_params", mint.MintParamsHandler},
		},
		"staking": {
			{"delegation", staking.DelegationHandler},
			{"delegation_message", staking.DelegationMessageHandler},
			{"redelegation", staking.RedelegationHandler},
			{"redelegation_message", staking.RedelegationMessageHandler},
			{"staking_params", staking.StakingParamsHandler},
			{"unbonding_delegation", staking.UnbondingDelegationHandler},
			{"unbonding_delegation_message", staking.UnbondingDelegationMessageHandler},
			{"edit_validator_message", staking.EditValidatorMessageHandler},
			{"create_validator_message", staking.CreateValidatorMessageHandler},
		},
		"authz": {
			{"exec_message", authz.ExecMessageHandler},
		},
		"slashing": {
			{"handle_validator_signature", slashing.HandleValidatorSignatureHandler},
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
