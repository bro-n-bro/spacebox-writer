package distribution

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// CommunityPoolHandler is a handler for community pool event
func CommunityPoolHandler(ctx context.Context, msg []byte, ch rep.Storage) error {
	val := model.CommunityPool{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	return ch.CommunityPool(val)
}
