package event

import (
	"time"

	"github.com/pkg/errors"
	m "github.com/smartnuance/saas-kit/pkg/event/dbmodels"

	"github.com/gin-gonic/gin"
)

// CreateWorkshopBody describes the workshop to be created
type CreateWorkshopBody struct {
	WorkshopInfo
	// Starts must be provided as RFC 3339 strings
	Starts time.Time `json:"starts"`
	// Ends must be provided as RFC 3339 strings
	Ends    time.Time `json:"ends,omitempty"`
	EventID string    `json:"eventID"`
}

func (s *Service) CreateWorkshop(ctx *gin.Context) (workshop *m.Workshop, err error) {
	var body CreateWorkshopBody
	err = ctx.ShouldBind(&body)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	eventID := body.EventID
	if eventID == "" {
		// create event here
	}

	workshop, err = s.DBAPI.CreateWorkshop(ctx, &body, eventID)
	return
}

var (
	ErrUnauthorized = errors.New("role insufficient to act on desired instances")
)
