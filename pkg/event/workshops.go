package event

import (
	"time"

	"github.com/pkg/errors"
	m "github.com/smartnuance/saas-kit/pkg/event/dbmodels"
	"github.com/smartnuance/saas-kit/pkg/lib/roles"

	"github.com/gin-gonic/gin"
)

// WorkshopData describes the workshop to be created
type WorkshopData struct {
	WorkshopInfo
	// Starts must be provided as RFC 3339 strings
	Starts time.Time `json:"starts"`
	// Ends must be provided as RFC 3339 strings
	Ends    time.Time `json:"ends,omitempty"`
	EventID string    `json:"eventID"`
}

// EventData describes an Event
type EventData struct {
	EventInfo
	// Starts must be provided as RFC 3339 strings
	Starts time.Time `json:"starts"`
	// Ends must be provided as RFC 3339 strings
	Ends time.Time `json:"ends,omitempty"`
}

func (s *Service) CreateWorkshop(ctx *gin.Context) (workshop *m.Workshop, err error) {
	var instanceID string
	_, instanceID, err = roles.ApplyHeaders(ctx)
	if err != nil {
		return
	}

	// Check permission
	if !roles.CanActIn(ctx, roles.RoleEventOrganizer) {
		err = errors.WithStack(ErrUnauthorized)
		return
	}

	var data WorkshopData
	err = ctx.ShouldBind(&data)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	var event *m.Event
	if data.EventID == "" {
		// create event for this specific event
		event, err = s.DBAPI.CreateEvent(ctx, instanceID, &EventData{
			EventInfo: EventInfo{
				Title:        data.Title,
				LocationName: data.LocationName,
				LocationURL:  data.LocationURL,
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
	}

	if event.InstanceID == "" {
		_, _, event.InstanceID, err = roles.From(ctx)
		if err != nil {
			err = errors.WithStack(err)
			return
		}
	}

	workshop, err = s.DBAPI.CreateWorkshop(ctx, &data)
	return
}

var (
	ErrUnauthorized = errors.New("role insufficient to act on desired instances")
)
