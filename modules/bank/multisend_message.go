package bank

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// MultiSendMessageHandler is a handler for multi send message event
func MultiSendMessageHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.MultiSendMessage](msgs)
	if err != nil {
		return err
	}

	return ch.MultiSendMessage(vals)
}
