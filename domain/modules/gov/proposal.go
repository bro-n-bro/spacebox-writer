package gov

import (
	"context"

	"spacebox-writer/adapter/clickhouse"
)

func ProposalHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	return nil // TODO: clarify logic and do this
}
