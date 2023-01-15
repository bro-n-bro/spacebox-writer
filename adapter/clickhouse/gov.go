package clickhouse

import (
	jsoniter "github.com/json-iterator/go"

	storageModel "github.com/hexy-dev/spacebox-writer/adapter/clickhouse/models"
	"github.com/hexy-dev/spacebox/broker/model"
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

	if err = ch.gorm.Table("gov_params").Create(storageModel.GovParams{
		DepositParams: string(depositParamsBytes),
		VotingParams:  string(votingParamsBytes),
		TallyParams:   string(tallyParamsBytes),
		Height:        val.Height,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (ch *Clickhouse) Proposal(val model.Proposal) (err error) {
	if err = ch.gorm.Table("proposal").Create(storageModel.Proposal{
		ID:              val.ID,
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
	}).Error; err != nil {
		return err
	}

	return nil
}

func (ch *Clickhouse) ProposalDepositMessage(val model.ProposalDepositMessage) (err error) {
	var (
		coinsBytes []byte
	)

	if coinsBytes, err = jsoniter.Marshal(val.Coins); err != nil {
		return err
	}

	if err = ch.gorm.Table("proposal_deposit_message").Create(storageModel.ProposalDepositMessage{
		DepositorAddress: val.DepositorAddress,
		Coins:            string(coinsBytes),
		TxHash:           val.TxHash,
		ProposalID:       int64(val.ProposalID),
		Height:           val.Height,
		MsgIndex:         val.MsgIndex,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (ch *Clickhouse) ProposalTallyResult(val model.ProposalTallyResult) (err error) {
	if err = ch.gorm.Table("proposal_tally_result").Create(storageModel.ProposalTallyResult{
		ProposalID: int64(val.ProposalID),
		Yes:        val.Yes,
		No:         val.No,
		Abstain:    val.Abstain,
		NoWithVeto: val.NoWithVeto,
		Height:     val.Height,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (ch *Clickhouse) ProposalVoteMessage(val model.ProposalVoteMessage) (err error) {
	if err = ch.gorm.Table("proposal_vote_message").Create(storageModel.ProposalVoteMessage{
		ProposalID:   val.ProposalID,
		VoterAddress: val.VoterAddress,
		Option:       val.Option,
		Height:       val.Height,
	}).Error; err != nil {
		return err
	}

	return nil
}
