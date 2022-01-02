package webbff

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/smartnuance/saas-kit/pkg/graph/models"
	"github.com/smartnuance/saas-kit/pkg/graph/queries"
)

func (r *mutationResolver) CreateWorkshop(ctx context.Context, input models.NewWorkshop) (*models.Workshop, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Workshops(ctx context.Context) ([]*models.Workshop, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns queries.MutationResolver implementation.
func (r *Resolver) Mutation() queries.MutationResolver { return &mutationResolver{r} }

// Query returns queries.QueryResolver implementation.
func (r *Resolver) Query() queries.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
