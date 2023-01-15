package staking

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"spacebox-writer/adapter/clickhouse"

	"github.com/hexy-dev/spacebox/broker/model"
)

func ValidatorInfoHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.ValidatorInfo{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return err
	}

	return ch.ValidatorInfo(val)
}
