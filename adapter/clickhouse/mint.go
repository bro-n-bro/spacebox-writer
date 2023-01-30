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
	tableAnnualProvision = "annual_provision"
	tableMintParams      = "mint_params"
)

func (ch *Clickhouse) AnnualProvision(val model.AnnualProvision) (err error) {
	return ch.gorm.Table(tableAnnualProvision).Create(val).Error
}

func (ch *Clickhouse) MintParams(val model.MintParams) (err error) {
	var (
		paramsBytes    []byte
		updates        storageModel.MintParams
		prevValStorage storageModel.MintParams
		valStorage     storageModel.MintParams
	)

	if paramsBytes, err = jsoniter.Marshal(val.Params); err != nil {
		return err
	}

	valStorage = storageModel.MintParams{
		Params: string(paramsBytes),
		Height: val.Height,
	}

	if err = ch.gorm.Table(tableMintParams).
		First(&prevValStorage).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = ch.gorm.Table(tableMintParams).Create(valStorage).Error; err != nil {
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
		if err = ch.gorm.Table(tableMintParams).
			Where("height = ?", valStorage.Height).
			Updates(&updates).Error; err != nil {
			return errors.Wrap(err, "error of update")
		}
	}
	return nil
}
