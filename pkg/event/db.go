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
	"github.com/volatiletech/sqlboiler/v4/types"
	// . "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type DBAPI interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	Commit(tx *sql.Tx) error
	Rollback(tx *sql.Tx) error
	CreateWorkshop(ctx context.Context, data *WorkshopData) (workshop *m.Workshop, err error)
	CreateEvent(ctx context.Context, instanceID string, data *EventData) (event *m.Event, err error)
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

func (db *dbAPI) CreateWorkshop(ctx context.Context, data *WorkshopData) (workshop *m.Workshop, err error) {
	var info types.JSON
	info, err = json.Marshal(data.WorkshopInfo)
	if err != nil {
		return
	}
	workshop = &m.Workshop{
		ID:           xid.New().String(),
		Info:         info,
		Starts:       data.Starts,
		Ends:         null.TimeFrom(data.Ends),
		EventID:      data.EventID,
		Participants: types.JSON("{}"),
	}
	err = workshop.Upsert(ctx, db.DB, true, boil.None().Cols, boil.Infer(), boil.Infer())
	if err != nil {
		return
	}
	return
}

func (db *dbAPI) CreateEvent(ctx context.Context, instanceID string, data *EventData) (event *m.Event, err error) {
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
		Ends:       null.TimeFrom(data.Ends),
		InstanceID: instanceID,
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

var (
	ErrEventDoesNotExist = errors.New("event does not exist")
)
