package rep

import "github.com/bro-n-bro/spacebox/broker/model"

type Storage interface {
	Account(model.Account) error

	AccountBalance(model.AccountBalance) error
	MultiSendMessage(model.MultiSendMessage) error
	SendMessage(model.SendMessage) error
	Supply(model.Supply) error

	Delegation(model.Delegation) error
	DelegationMessage(model.DelegationMessage) error
	Redelegation(model.Redelegation) error
	RedelegationMessage(model.RedelegationMessage) error
	StakingParams(model.StakingParams) error
	StakingPool(model.StakingPool) error
	UnbondingDelegation(model.UnbondingDelegation) error
	UnbondingDelegationMessage(model.UnbondingDelegationMessage) error
	Validator(model.Validator) error
	ValidatorInfo(model.ValidatorInfo) error
	ValidatorStatus(model.ValidatorStatus) error
	ValidatorDescription(model.ValidatorDescription) error

	AnnualProvision(model.AnnualProvision) error
	MintParams(model.MintParams) error

	GovParams(model.GovParams) error
	Proposal(model.Proposal) error
	ProposalTallyResult(model.ProposalTallyResult) error
	ProposalVoteMessage(model.ProposalVoteMessage) error
	ProposalDepositMessage(val model.ProposalDepositMessage) error

	CommunityPool(model.CommunityPool) error
	DistributionParams(model.DistributionParams) error
	DelegationRewardMessage(message model.DelegationRewardMessage) error
	ValidatorCommission(model.ValidatorCommission) error

	Block(model.Block) error
	Message(model.Message) error
	Transaction(model.Transaction) error
}
