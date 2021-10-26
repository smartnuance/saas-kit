package event

//go:generate sqlboiler --config sqlboiler.toml psql
//go:generate mockgen -destination db_mock.go -package event . DBAPI

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	m "github.com/smartnuance/saas-kit/pkg/event/dbmodels"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
	// . "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type DBAPI interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	Commit(tx *sql.Tx) error
	Rollback(tx *sql.Tx) error
	CreateWorkshop(ctx context.Context, data *CreateWorkshopData) (workshop *m.Workshop, err error)
	ListWorkshops(ctx context.Context, instanceID string) (workshop []WorkshopData, err error)
	CreateEvent(ctx context.Context, data *EventData) (event *m.Event, err error)
	GetEvent(ctx context.Context, eventID string) (event *m.Event, err error)
}

type dbAPI struct {
	DB *sql.DB
}

func (db *dbAPI) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return db.DB.BeginTx(ctx, nil)
}

func (db *dbAPI) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (db *dbAPI) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

func (db *dbAPI) CreateWorkshop(ctx context.Context, data *CreateWorkshopData) (workshop *m.Workshop, err error) {
	var info types.JSON
	info, err = json.Marshal(data.WorkshopInfo)
	if err != nil {
		return
	}
	id := data.ID
	if id == "" {
		id = xid.New().String()
	}
	workshop = &m.Workshop{
		ID:           id,
		Info:         info,
		Starts:       data.Starts,
		Ends:         null.TimeFromPtr(data.Ends),
		EventID:      data.EventID,
		Participants: types.JSON("{}"),
	}
	err = workshop.Upsert(ctx, db.DB, true, boil.None().Cols, boil.Infer(), boil.Infer())
	if err != nil {
		return
	}
	return
}

func (db *dbAPI) ListWorkshops(ctx context.Context, instanceID string) ([]WorkshopData, error) {
	results, err := m.Workshops(qm.Load(m.WorkshopRels.Event, m.EventWhere.InstanceID.EQ(instanceID)), qm.Limit(10)).All(ctx, db.DB)
	if err == sql.ErrNoRows {
		// wrap sql error in specific error of event context
		return nil, errors.Wrap(ErrRetrieveWorkshopList, err.Error())
	}

	workshops := []WorkshopData{}
	for _, w := range results {
		if w == nil {
			return nil, errors.New("got nil workshop row")
		}

		workshop, err := loadWorkshop(*w)
		if err != nil {
			if errors.Is(err, ErrWorkshopWithNoEvent) {
				log.Warn().Err(err).Str("workshop.ID", w.ID).Msg("skipping")
				continue
			} else {
				return nil, err
			}
		}

		workshops = append(workshops, workshop)
	}

	return workshops, nil
}

func (db *dbAPI) CreateEvent(ctx context.Context, data *EventData) (event *m.Event, err error) {
	var info types.JSON
	info, err = json.Marshal(data.EventInfo)
	if err != nil {
		return
	}
	log.Debug().Msg(string(info))
	event = &m.Event{
		ID:         xid.New().String(),
		Info:       info,
		Starts:     data.Starts,
		Ends:       null.TimeFromPtr(data.Ends),
		InstanceID: data.InstanceID,
	}
	err = event.Upsert(ctx, db.DB, true, boil.None().Cols, boil.Infer(), boil.Infer())
	if err != nil {
		return
	}
	return
}

func (db *dbAPI) GetEvent(ctx context.Context, eventID string) (event *m.Event, err error) {
	event, err = m.Events(m.EventWhere.ID.EQ(eventID)).One(ctx, db.DB)
	if err == sql.ErrNoRows {
		// transform sql error in specific error of event context
		err = errors.WithStack(ErrEventDoesNotExist)
		return
	}
	return
}

type EventInfo struct {
	Title        string `json:"title,omitempty"`
	LocationName string `json:"locationName,omitempty"`
	LocationURL  string `json:"locationURL,omitempty"`
}

type WorkshopInfo struct {
	Title        string `json:"title,omitempty"`
	LocationName string `json:"locationName,omitempty"`
	LocationURL  string `json:"locationURL,omitempty"`
	Couples      bool   `json:"couples"`
}

func loadEvent(row m.Event) (event EventData, err error) {
	var eventInfo EventInfo
	err = json.Unmarshal(row.Info, &eventInfo)
	if err != nil {
		return
	}
	event = EventData{
		InstanceID: row.InstanceID,
		EventInfo:  eventInfo,
		Starts:     event.Starts,
		Ends:       event.Ends,
	}
	return
}

func loadWorkshop(row m.Workshop) (workshop WorkshopData, err error) {
	var info WorkshopInfo
	err = json.Unmarshal(row.Info, &info)
	if err != nil {
		return
	}

	eventRow := row.R.Event
	if eventRow == nil {
		err = ErrWorkshopWithNoEvent
		return
	}
	var event EventData
	event, err = loadEvent(*eventRow)
	if err != nil {
		return
	}
	workshop = WorkshopData{
		ID:           row.ID,
		WorkshopInfo: info,
		Starts:       row.Starts,
		Ends:         row.Ends.Ptr(),
		EventID:      row.EventID,
		EventData:    event,
	}
	return
}

var (
	ErrEventDoesNotExist    = errors.New("event does not exist")
	ErrRetrieveWorkshopList = errors.New("retrieving workshop list failed")
	ErrWorkshopWithNoEvent  = errors.New("workshop with no associated event found")
)
