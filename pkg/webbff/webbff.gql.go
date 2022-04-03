package webbff

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/smartnuance/saas-kit/pkg/graph/models"
	"github.com/smartnuance/saas-kit/pkg/graph/queries"
	"github.com/smartnuance/saas-kit/pkg/lib/roles"
)

func (r *mutationResolver) CreateWorkshop(ctx context.Context, input models.WorkshopInput) (*models.Workshop, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Workshops(ctx context.Context) (*models.WorkshopList, error) {
	req, err := http.NewRequest("GET", "http://"+r.Service.eventServiceAddress+"/workshop/list", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", ctx.Value("Authorization").(string))
	req.Header.Add(roles.RoleHeader, ctx.Value(roles.RoleKey).(string))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var workshopList models.WorkshopList
	err = json.NewDecoder(resp.Body).Decode(&workshopList)
	if err != nil {
		return nil, err
	}
	return &workshopList, nil
}

// Mutation returns queries.MutationResolver implementation.
func (r *Resolver) Mutation() queries.MutationResolver { return &mutationResolver{r} }

// Query returns queries.QueryResolver implementation.
func (r *Resolver) Query() queries.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
