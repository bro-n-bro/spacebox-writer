package staking

import (
	"context"
	"encoding/json"

	"github.com/hexy-dev/spacebox-writer/adapter/clickhouse"
	"github.com/hexy-dev/spacebox/broker/model"
)

func RedelegationMessageHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.RedelegationMessage{}
	if err := json.Unmarshal(msg, &val); err != nil {
		return err
	}

	return ch.RedelegationMessage(val)
}
