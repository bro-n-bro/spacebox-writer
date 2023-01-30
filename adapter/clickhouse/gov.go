package clickhouse

import (
	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	storageModel "github.com/bro-n-bro/spacebox-writer/adapter/clickhouse/models"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	tableGovParams              = "gov_params"
	tableProposal               = "proposal"
	tableProposalDepositMessage = "proposal_deposit_message"
	tableProposalVoteMessage    = "proposal_vote_message"
	tableProposalTallyResult    = "proposal_tally_result"
)

func (ch *Clickhouse) GovParams(val model.GovParams) (err error) {
	var (
		tallyParamsBytes   []byte
		votingParamsBytes  []byte
		depositParamsBytes []byte
		valStorage         storageModel.GovParams
		prevValStorage     storageModel.GovParams
		updates            storageModel.GovParams
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

	valStorage = storageModel.GovParams{
		DepositParams: string(depositParamsBytes),
		VotingParams:  string(votingParamsBytes),
		TallyParams:   string(tallyParamsBytes),
		Height:        val.Height,
	}

	if err = ch.gorm.Table(tableGovParams).
		First(&prevValStorage).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = ch.gorm.Table(tableGovParams).Create(val).Error; err != nil {
				return errors.Wrap(err, "error of create")
			}
			return nil
		}
		return errors.Wrap(err, "error of database: %w")
	}

	if valStorage.Height > prevValStorage.Height {
		if err = copier.Copy(&valStorage, &updates); err != nil {
			return errors.Wrap(err, "error of prepare update")
		}
		if err = ch.gorm.Table(tableGovParams).
			Where("height = ?", valStorage.Height).
			Updates(&updates).Error; err != nil {
			return errors.Wrap(err, "error of update")
		}
	}

	return nil
}

// Proposal TODO: если status = pass / error / reject -  то это конечное сообщение и его не надо обновлять
func (ch *Clickhouse) Proposal(val model.Proposal) (err error) {
	if err = ch.gorm.Table(tableProposal).Create(storageModel.Proposal{
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
		exists     bool
	)

	if coinsBytes, err = jsoniter.Marshal(val.Coins); err != nil {
		return err
	}

	if exists, err = ch.ExistsTx(tableProposalDepositMessage, val.TxHash, val.MsgIndex); err != nil {
		return err
	}

	if !exists {
		if err = ch.gorm.Table(tableProposalDepositMessage).Create(storageModel.ProposalDepositMessage{
			DepositorAddress: val.DepositorAddress,
			Coins:            string(coinsBytes),
			TxHash:           val.TxHash,
			ProposalID:       val.ProposalID,
			Height:           val.Height,
			MsgIndex:         val.MsgIndex,
		}).Error; err != nil {
			return err
		}
	}

	return nil
}

func (ch *Clickhouse) ProposalTallyResult(val model.ProposalTallyResult) (err error) {
	var (
		updates        model.ValidatorCommission
		prevValStorage model.ValidatorCommission
	)

	if err = ch.gorm.Table(tableProposalTallyResult).
		Where("proposal_id = ?", val.ProposalID).
		First(&prevValStorage).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = ch.gorm.Table(tableProposalTallyResult).Create(val).Error; err != nil {
				return errors.Wrap(err, "error of create")
			}
			return nil
		}
		return errors.Wrap(err, "error of database: %w")
	}

	if val.Height > prevValStorage.Height {
		if err = copier.Copy(&val, &updates); err != nil {
			return errors.Wrap(err, "error of prepare update")
		}
		if err = ch.gorm.Table(tableProposalTallyResult).
			Where("proposal_id = ?", val.ProposalID).
			Updates(&updates).Error; err != nil {
			return errors.Wrap(err, "error of update")
		}
	}

	return nil
}

// ProposalVoteMessage TODO: добавить в ProposalVoteMessage поле TxHash
func (ch *Clickhouse) ProposalVoteMessage(val model.ProposalVoteMessage) (err error) {
	if err = ch.gorm.Table(tableProposalVoteMessage).Create(storageModel.ProposalVoteMessage{
		ProposalID:   val.ProposalID,
		VoterAddress: val.VoterAddress,
		Option:       val.Option,
		Height:       val.Height,
	}).Error; err != nil {
		return err
	}

	return nil
}
