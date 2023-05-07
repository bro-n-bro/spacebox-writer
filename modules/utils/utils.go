package utils

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox/broker/model"
)

type models interface {
	model.SendMessage | model.Supply | model.AccountBalance | model.MultiSendMessage |
		model.Transaction | model.CommunityPool | model.DelegationRewardMessage | model.DistributionParams |
		model.GovParams | model.Proposal | model.ProposalDepositMessage | model.MintParams |
		model.Delegation | model.DelegationMessage | model.Redelegation | model.RedelegationMessage |
		model.StakingParams | model.UnbondingDelegation | model.UnbondingDelegationMessage | model.ProposerReward |
		model.DistributionCommission | model.SubmitProposalMessage | model.WithdrawValidatorCommissionMessage |
		model.DistributionReward | model.VoteWeightedMessage | model.EditValidatorMessage | model.ExecMessage |
		model.GrantAllowanceMessage | model.GrantMessage
}

func ConvertMessages[T models](msgs [][]byte) ([]T, error) {
	vals := make([]T, 0, len(msgs))
	for _, msg := range msgs {
		val := new(T)
		if err := jsoniter.Unmarshal(msg, &val); err != nil {
			return nil, errors.Wrap(err, "unmarshall error")
		}
		vals = append(vals, *val)
	}

	return vals, nil
}
