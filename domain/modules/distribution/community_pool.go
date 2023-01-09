package distribution

import (
	"context"

	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"

	"github.com/hexy-dev/spacebox/broker/model"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func CommunityPoolHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {

	val := model.CommunityPool{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	coinsBytes, err := jsoniter.Marshal(val.Coins)
	if err != nil {
		return err
	}

	if err = ch.GetGormDB(ctx).Table("community_pool").Create(storageModel.CommunityPool{
		Coins:  string(coinsBytes),
		Height: val.Height,
	}).Error; err != nil {
		return err
	}

	return nil
}
