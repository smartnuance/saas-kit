package event

import (
	"github.com/friendsofgo/errors"
	m "github.com/smartnuance/saas-kit/pkg/event/dbmodels"
	"github.com/smartnuance/saas-kit/pkg/graph/models"
	"github.com/smartnuance/saas-kit/pkg/lib/paging"
	"github.com/smartnuance/saas-kit/pkg/lib/roles"

	"github.com/gin-gonic/gin"
)

// WorkshopList describes the returned workshop list with paging.
type WorkshopList struct {
	Workshops []models.Workshop `json:"items"`
	Paging    paging.Paging     `json:"paging"`
}

func (s *Service) CreateWorkshop(ctx *gin.Context) (workshop *m.Workshop, err error) {
	// Check permission
	if !roles.CanActIn(ctx, roles.RoleEventOrganizer) {
		err = errors.WithStack(ErrUnauthorized)
		return
	}

	var data models.WorkshopInput
	err = ctx.ShouldBind(&data)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	// fallback to instance from context
	if data.Instance == "" {
		// fallback to default instance from headers
		data.Instance, err = roles.Instance(ctx)
		if err != nil {
			return
		}
	}

	if !roles.CanActFor(ctx, data.Instance) {
		err = errors.WithStack(ErrUnauthorized)
		return
	}

	var event *m.Event
	if data.EventID == "" {
		// create event for this specific event
		event, err = s.DBAPI.CreateEvent(ctx, &models.Event{
			Instance: &models.Instance{ID: data.Instance},
			EventInfo: &models.EventInfo{
				Title:        data.WorkshopInfo.Title,
				LocationName: data.WorkshopInfo.LocationName,
				LocationURL:  data.WorkshopInfo.LocationURL,
			},
			// assume same start/end of workshop
			Starts: data.Starts,
			Ends:   data.Ends,
		})
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		data.EventID = event.ID
	} else {
		event, err = s.DBAPI.GetEvent(ctx, data.EventID)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		data.EventID = event.ID
	}

	workshop, err = s.DBAPI.CreateWorkshop(ctx, &data)
	return
}

func (s *Service) ListWorkshops(ctx *gin.Context) (list WorkshopList, err error) {
	// Check permission
	if !roles.CanActIn(ctx, roles.RoleEventOrganizer) {
		r, _ := roles.FromContext(ctx)
		err = errors.Wrapf(ErrUnauthorized, "'%s' can not act as %s", r, roles.RoleEventOrganizer)
		return
	}

	var instanceID string
	instanceID, err = roles.Instance(ctx)
	if err != nil {
		err = errors.Wrap(ErrUnauthorized, err.Error())
		return
	}

	list, err = s.DBAPI.ListWorkshops(ctx, instanceID, paging.FromQuery(ctx))
	if err != nil {
		return
	}

	return
}

func (s *Service) DeleteWorkshop(ctx *gin.Context) (err error) {
	// Check permission
	if !roles.CanActIn(ctx, roles.RoleEventOrganizer) {
		r, _ := roles.FromContext(ctx)
		err = errors.Wrapf(ErrUnauthorized, "'%s' can not act as %s", r, roles.RoleEventOrganizer)
		return
	}

	_, err = roles.Instance(ctx)
	if err != nil {
		err = errors.Wrap(ErrUnauthorized, err.Error())
		return
	}

	err = s.DBAPI.DeleteWorkshop(ctx, ctx.Param("id"))
	if err != nil {
		return
	}

	return
}

var (
	ErrUnauthorized = errors.New("role insufficient to act on desired instances")
)
