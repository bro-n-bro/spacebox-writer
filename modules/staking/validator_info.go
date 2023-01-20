package staking

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
)

func ValidatorInfoHandler(ctx context.Context, msg []byte, ch rep.Storage) error {
	val := model.ValidatorInfo{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return err
	}

	return ch.ValidatorInfo(val)
}
