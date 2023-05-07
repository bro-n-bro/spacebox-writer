package rep

import "github.com/bro-n-bro/spacebox/broker/model"

type (
	// Storage is a repository for storage
	Storage interface {
		AccountBalance([]model.AccountBalance) error
		MultiSendMessage([]model.MultiSendMessage) error
		SendMessage([]model.SendMessage) error
		Supply([]model.Supply) error

		Delegation([]model.Delegation) error
		DelegationMessage([]model.DelegationMessage) error
		Redelegation([]model.Redelegation) error
		RedelegationMessage([]model.RedelegationMessage) error
		StakingParams([]model.StakingParams) error
		UnbondingDelegation([]model.UnbondingDelegation) error
		UnbondingDelegationMessage([]model.UnbondingDelegationMessage) error
		EditValidatorMessage([]model.EditValidatorMessage) error

		MintParams([]model.MintParams) error

		GovParams([]model.GovParams) error
		Proposal([]model.Proposal) error
		ProposalDepositMessage([]model.ProposalDepositMessage) error
		SubmitProposalMessage([]model.SubmitProposalMessage) error
		VoteWeightedMessage([]model.VoteWeightedMessage) error

		CommunityPool([]model.CommunityPool) error
		DistributionParams([]model.DistributionParams) error
		WithdrawValidatorCommissionMessage([]model.WithdrawValidatorCommissionMessage) error
		DelegationRewardMessage([]model.DelegationRewardMessage) error
		ProposerReward([]model.ProposerReward) error
		DistributionCommission([]model.DistributionCommission) error
		DistributionReward([]model.DistributionReward) error

		Transaction([]model.Transaction) error
	}
)
