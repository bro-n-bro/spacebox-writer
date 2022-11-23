package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"spacebox-writer/adapter/clickhouse"
	ct "spacebox-writer/graph/custom_types"
	"spacebox-writer/graph/generated"
)

// GetBlocks is the resolver for the getBlocks field.
func (r *queryResolver) GetBlocks(ctx context.Context) ([]*ct.Block, error) {
	_ctx := clickhouse.GetContext(ctx)
	var blocks []*ct.Block
	err := _ctx.Database.Find(&blocks).Error
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
