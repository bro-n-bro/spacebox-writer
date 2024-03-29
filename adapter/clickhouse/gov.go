package clickhouse

import (
	"database/sql"

	jsoniter "github.com/json-iterator/go"

	storageModel "github.com/bro-n-bro/spacebox-writer/adapter/clickhouse/models"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	tableGovParams              = "gov_params"
	tableProposal               = "proposal"
	tableProposalDepositMessage = "proposal_deposit_message"
	tableSubmitProposalMessage  = "submit_proposal_message"
	tableVoteWeightedMessage    = "vote_weighted_message"
)

// GovParams is a method for saving gov params data to clickhouse
func (ch *Clickhouse) GovParams(vals []model.GovParams) (err error) {
	var (
		tallyParams, votingParams, depositParams string
	)

	batch := make([]storageModel.GovParams, len(vals))
	for i, val := range vals {
		if tallyParams, err = jsoniter.MarshalToString(val.TallyParams); err != nil {
			return err
		}
		if votingParams, err = jsoniter.MarshalToString(val.VotingParams); err != nil {
			return err
		}
		if depositParams, err = jsoniter.MarshalToString(val.DepositParams); err != nil {
			return err
		}

		batch[i] = storageModel.GovParams{
			DepositParams: depositParams,
			VotingParams:  votingParams,
			TallyParams:   tallyParams,
			Height:        val.Height,
		}
	}

	return ch.gorm.Table(tableGovParams).CreateInBatches(batch, len(batch)).Error
}

// Proposal is a method for saving proposal data to clickhouse
func (ch *Clickhouse) Proposal(vals []model.Proposal) (err error) {
	batch := make([]storageModel.Proposal, len(vals))
	for i, val := range vals {
		batch[i] = storageModel.Proposal{
			ID:              int64(val.ID),
			Title:           val.Title,
			Description:     val.Description,
			ProposalRoute:   val.ProposalRoute,
			ProposalType:    val.ProposalType,
			ProposerAddress: val.ProposerAddress,
			Status:          val.Status,
			Content:         string(val.Content),
			SubmitTime: sql.NullTime{
				Time:  val.SubmitTime,
				Valid: !val.SubmitTime.IsZero(),
			},
			DepositEndTime: sql.NullTime{
				Time:  val.DepositEndTime,
				Valid: !val.DepositEndTime.IsZero(),
			},
			VotingStartTime: sql.NullTime{
				Time:  val.VotingStartTime,
				Valid: !val.VotingStartTime.IsZero(),
			},
			VotingEndTime: sql.NullTime{
				Time:  val.VotingEndTime,
				Valid: !val.VotingEndTime.IsZero(),
			},
		}
	}

	return ch.gorm.Table(tableProposal).CreateInBatches(batch, len(batch)).Error
}

// ProposalDepositMessage is a method for saving proposal deposit message data to clickhouse
func (ch *Clickhouse) ProposalDepositMessage(vals []model.ProposalDepositMessage) (err error) {
	var (
		coins string
	)

	batch := make([]storageModel.ProposalDepositMessage, len(vals))
	for i, val := range vals {
		if coins, err = jsoniter.MarshalToString(val.Coins); err != nil {
			return err
		}
		batch[i] = storageModel.ProposalDepositMessage{
			DepositorAddress: val.DepositorAddress,
			Coins:            coins,
			TxHash:           val.TxHash,
			ProposalID:       int64(val.ProposalID),
			Height:           val.Height,
			MsgIndex:         val.MsgIndex,
		}
	}

	return ch.gorm.Table(tableProposalDepositMessage).CreateInBatches(batch, len(batch)).Error
}

// SubmitProposalMessage is a method for saving submit proposal message data to clickhouse
func (ch *Clickhouse) SubmitProposalMessage(vals []model.SubmitProposalMessage) (err error) {
	var (
		initialDeposit string
	)

	batch := make([]storageModel.SubmitProposalMessage, len(vals))
	for i, val := range vals {
		if initialDeposit, err = jsoniter.MarshalToString(val.InitialDeposit); err != nil {
			return err
		}

		batch[i] = storageModel.SubmitProposalMessage{
			TxHash:         val.TxHash,
			Proposer:       val.Proposer,
			InitialDeposit: initialDeposit,
			Title:          val.Title,
			Description:    val.Description,
			Type:           val.Type,
			ProposalID:     int64(val.ProposalID),
			Height:         val.Height,
			MsgIndex:       val.MsgIndex,
		}
	}

	return ch.gorm.Table(tableSubmitProposalMessage).CreateInBatches(batch, len(batch)).Error
}

// VoteWeightedMessage is a method for saving vote weighted message data to clickhouse
func (ch *Clickhouse) VoteWeightedMessage(vals []model.VoteWeightedMessage) (err error) {
	var (
		weightedVoteOption string
	)

	batch := make([]storageModel.VoteWeightedMessage, len(vals))
	for i, val := range vals {
		if weightedVoteOption, err = jsoniter.MarshalToString(val.WeightedVoteOption); err != nil {
			return err
		}

		batch[i] = storageModel.VoteWeightedMessage{
			TxHash:             val.TxHash,
			ProposalID:         int64(val.ProposalId),
			Height:             val.Height,
			MsgIndex:           val.MsgIndex,
			Voter:              val.Voter,
			WeightedVoteOption: weightedVoteOption,
		}
	}

	return ch.gorm.Table(tableVoteWeightedMessage).CreateInBatches(batch, len(batch)).Error
}
