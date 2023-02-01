package clickhouse

import (
	"fmt"

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

// GovParams TODO: unsupported Scan, storing driver.Value type *map[string]interface {} into type *string
func (ch *Clickhouse) GovParams(val model.GovParams) (err error) {
	var (
		//tallyParamsBytes   []byte
		//votingParamsBytes  []byte
		//depositParamsBytes []byte
		//valToStorage       storageModel.GovParams
		prevValStorage storageModel.GovParams
		//updates        storageModel.GovParams
	)

	//if tallyParamsBytes, err = jsoniter.Marshal(val.TallyParams); err != nil {
	//	return err
	//}
	//if votingParamsBytes, err = jsoniter.Marshal(val.VotingParams); err != nil {
	//	return err
	//}
	//if depositParamsBytes, err = jsoniter.Marshal(val.DepositParams); err != nil {
	//	return err
	//}

	//valToStorage = storageModel.GovParams{
	//	DepositParams: string(depositParamsBytes),
	//	VotingParams:  string(votingParamsBytes),
	//	TallyParams:   string(tallyParamsBytes),
	//	Height:        val.Height,
	//}

	row := ch.sql.QueryRow("SELECT deposit_params, voting_params, tally_params, height FROM spacebox.gov_params LIMIT 1")
	if err = row.Scan(
		&prevValStorage.DepositParams,
		&prevValStorage.VotingParams,
		&prevValStorage.TallyParams,
		&prevValStorage.Height,
	); err != nil {
		return err
	}

	fmt.Println(prevValStorage.DepositParams)
	fmt.Println(prevValStorage.VotingParams)
	fmt.Println(prevValStorage.TallyParams)
	fmt.Println(prevValStorage.Height)

	//if err = ch.gorm.Table(tableGovParams).
	//	First(&prevValStorage).Error; err != nil {
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//if err = ch.gorm.Table(tableGovParams).Create(valToStorage).Error; err != nil {
	//	return errors.Wrap(err, "error of create")
	//}
	//		return nil
	//	}
	//	return errors.Wrap(err, "error of database: %w")
	//}
	//
	//fmt.Println(prevValStorage, valToStorage)
	//if valToStorage.Height > prevValStorage.Height {
	//	if err = copier.Copy(&updates, &valToStorage); err != nil {
	//		return errors.Wrap(err, "error of prepare update")
	//	}
	//	if err = ch.gorm.Table(tableGovParams).
	//		Where("height = ?", prevValStorage.Height).
	//		Updates(&updates).Error; err != nil {
	//		return errors.Wrap(err, "error of update")
	//	}
	//}

	return nil
}

// Proposal TODO: если status = pass / error / reject -  то это конечное сообщение и его не надо обновлять
// TODO: unsupported Scan, storing driver.Value type *map[string]interface {} into type *string
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
		updates        model.ProposalTallyResult
		prevValStorage model.ProposalTallyResult
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

func (ch *Clickhouse) ProposalVoteMessage(val model.ProposalVoteMessage) (err error) {
	var (
		exists bool
	)

	if exists, err = ch.ExistsTx(tableProposalVoteMessage, val.TxHash, val.MsgIndex); err != nil {
		return err
	}

	if !exists {
		if err = ch.gorm.Table(tableProposalVoteMessage).Create(storageModel.ProposalVoteMessage{
			ProposalID:   val.ProposalID,
			VoterAddress: val.VoterAddress,
			Option:       val.Option,
			Height:       val.Height,
			TxHash:       val.TxHash,
		}).Error; err != nil {
			return err
		}
	}

	return nil
}
