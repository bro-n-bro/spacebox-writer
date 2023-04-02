package bank

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// SendMessageHandler is a handler for send message event
func SendMessageHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.SendMessage](msgs)
	if err != nil {
		return err
	}

	return ch.SendMessage(vals)
}
