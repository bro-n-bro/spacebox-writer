package clickhouse

import (
	jsoniter "github.com/json-iterator/go"

	storageModel "github.com/bro-n-bro/spacebox-writer/adapter/clickhouse/models"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	tableGovParams              = "gov_params"
	tableProposal               = "proposal"
	tableProposalDepositMessage = "proposal_deposit_message"
)

func (ch *Clickhouse) GovParams(val model.GovParams) (err error) {
	var (
		tallyParamsBytes   []byte
		votingParamsBytes  []byte
		depositParamsBytes []byte
	)

	if tallyParamsBytes, err = jsoniter.Marshal(val.TallyParams); err != nil {
		return err
	}
	if votingParamsBytes, err = jsoniter.Marshal(val.VotingParams); err != nil {
		return err
	}
	if depositParamsBytes, err = jsoniter.Marshal(val.DepositParams); err != nil {
		return err
	}

	return ch.gorm.Table(tableGovParams).Create(storageModel.GovParams{
		DepositParams: string(depositParamsBytes),
		VotingParams:  string(votingParamsBytes),
		TallyParams:   string(tallyParamsBytes),
		Height:        val.Height,
	}).Error
}

func (ch *Clickhouse) Proposal(val model.Proposal) (err error) {
	return ch.gorm.Table(tableProposal).Create(storageModel.Proposal{
		ID:              int64(val.ID),
		Title:           val.Title,
		Description:     val.Description,
		ProposalRoute:   val.ProposalRoute,
		ProposalType:    val.ProposalType,
		ProposerAddress: val.ProposerAddress,
		Status:          val.Status,
		Content:         string(val.Content),
		SubmitTime:      val.SubmitTime,
		DepositEndTime:  val.DepositEndTime,
		VotingStartTime: val.VotingStartTime,
		VotingEndTime:   val.VotingEndTime,
	}).Error
}

func (ch *Clickhouse) ProposalDepositMessage(val model.ProposalDepositMessage) (err error) {
	var (
		coinsBytes []byte
	)

	if coinsBytes, err = jsoniter.Marshal(val.Coins); err != nil {
		return err
	}

	return ch.gorm.Table(tableProposalDepositMessage).Create(storageModel.ProposalDepositMessage{
		DepositorAddress: val.DepositorAddress,
		Coins:            string(coinsBytes),
		TxHash:           val.TxHash,
		ProposalID:       int64(val.ProposalID),
		Height:           val.Height,
		MsgIndex:         val.MsgIndex,
	}).Error
}
