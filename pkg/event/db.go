package event

//go:generate sqlboiler --config sqlboiler.toml psql
//go:generate mockgen -destination db_mock.go -package event . DBAPI

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/friendsofgo/errors"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	m "github.com/smartnuance/saas-kit/pkg/event/dbmodels"
	"github.com/smartnuance/saas-kit/pkg/lib/paging"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/types"
	"google.golang.org/protobuf/types/known/timestamppb"
	// . "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type DBAPI interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	Commit(tx *sql.Tx) error
	Rollback(tx *sql.Tx) error
	CreateWorkshop(ctx context.Context, data *Workshop) (workshop *m.Workshop, err error)
	ListWorkshops(ctx context.Context, instanceID string, page paging.Page) (list WorkshopList, err error)
	GetWorkshop(ctx context.Context, workshopID string) (workshop *m.Workshop, err error)
	DeleteWorkshop(ctx context.Context, workshopID string) (err error)
	CreateEvent(ctx context.Context, data *Event) (event *m.Event, err error)
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

func (db *dbAPI) CreateWorkshop(ctx context.Context, data *Workshop) (workshop *m.Workshop, err error) {
	var info types.JSON
	info, err = json.Marshal(data.WorkshopInfo)
	if err != nil {
		return
	}
	id := data.Id
	if id == "" {
		id = xid.New().String()
	}
	var eventID string
	switch e := data.BelongsTo.(type) {
	case *Workshop_Event:
		eventID = e.Event.Id
	case *Workshop_EventID:
		eventID = e.EventID
	}
	workshop = &m.Workshop{
		ID:           id,
		Info:         info,
		Starts:       data.Starts.AsTime(),
		Ends:         null.NewTime(data.GetEnds().AsTime(), data.GetEnds() != nil),
		EventID:      eventID,
		Participants: types.JSON("{}"),
	}
	err = workshop.Upsert(ctx, db.DB, true, boil.None().Cols, boil.Infer(), boil.Infer())
	if err != nil {
		return
	}
	return
}

func (db *dbAPI) ListWorkshops(ctx context.Context, instanceID string, page paging.Page) (list WorkshopList, err error) {
	results, err := m.Workshops(
		qm.InnerJoin(fmt.Sprintf("%s on %s = %s", m.TableNames.Events, m.EventTableColumns.ID, m.WorkshopColumns.EventID)),
		qm.Load(m.WorkshopRels.Event),
		m.EventWhere.InstanceID.EQ(instanceID),
		m.EventWhere.ID.Page(page),
	).All(ctx, db.DB)
	if err == sql.ErrNoRows {
		// wrap sql error in specific error of event context
		err = errors.Wrap(ErrRetrieveWorkshopList, err.Error())
		return
	}

	list.Items = []*Workshop{}
	for _, w := range results {
		if w == nil {
			err = errors.New("got nil workshop row")
			return
		}

		var workshop Workshop
		workshop, err = loadWorkshop(*w)
		if err != nil {
			return
		}

		list.Items = append(list.Items, &workshop)
	}

	list.Paging = &paging.Paging{Cur: &paging.Paging_Current{PageSize: int32(len(list.Items))}}
	if len(list.Items) > 0 {
		list.Paging.Cur.Start = list.Items[0].Id
		list.Paging.Cur.End = list.Items[len(list.Items)-1].Id

		_, isFirst := page.(*paging.Paging_First)
		if !isFirst {
			list.Paging.Prev = &paging.Paging_Previous{
				End:      list.Items[0].Id,
				PageSize: int32(page.Size()),
			}
		}
		isLast := len(list.Items) < page.Size()
		if !isLast {
			list.Paging.Next = &paging.Paging_Next{
				Start:    list.Items[len(list.Items)-1].Id,
				PageSize: int32(page.Size()),
			}
		}
	}

	return
}

func (db *dbAPI) DeleteWorkshop(ctx context.Context, workshopID string) (err error) {
	workshop, err := db.GetWorkshop(ctx, workshopID)
	if errors.Is(err, ErrWorkshopDoesNotExist) {
		// for idempotence: return no error
		return nil
	}
	if err != nil {
		return
	}
	_, err = workshop.Delete(ctx, db.DB, false)
	return
}

func (db *dbAPI) CreateEvent(ctx context.Context, data *Event) (event *m.Event, err error) {
	var info types.JSON
	info, err = json.Marshal(data.EventInfo)
	if err != nil {
		return
	}
	log.Debug().Msg(string(info))
	event = &m.Event{
		ID:         xid.New().String(),
		Info:       info,
		Starts:     data.Starts.AsTime(),
		Ends:       null.NewTime(data.GetEnds().AsTime(), data.GetEnds() != nil),
		InstanceID: data.Instance.Id,
	}
	err = event.Upsert(ctx, db.DB, true, boil.None().Cols, boil.Infer(), boil.Infer())
	if err != nil {
		return
	}
	return
}

func (db *dbAPI) GetWorkshop(ctx context.Context, workshopID string) (workshop *m.Workshop, err error) {
	workshop, err = m.Workshops(m.WorkshopWhere.ID.EQ(workshopID)).One(ctx, db.DB)
	if err == sql.ErrNoRows {
		// transform sql error in specific error of event context
		err = errors.WithStack(ErrWorkshopDoesNotExist)
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

func loadEvent(row m.Event) (event Event, err error) {
	var eventInfo Event_Info
	err = json.Unmarshal(row.Info, &eventInfo)
	if err != nil {
		return
	}
	event = Event{
		Id:        row.InstanceID,
		EventInfo: &eventInfo,
		Starts:    event.Starts,
		Ends:      event.Ends,
	}
	return
}

func loadWorkshop(row m.Workshop) (workshop Workshop, err error) {
	var info Workshop_Info
	err = json.Unmarshal(row.Info, &info)
	if err != nil {
		return
	}

	eventRow := row.R.Event
	if eventRow == nil {
		err = errors.New("workshop with no associated event found")
		log.Error().Err(err).Msg("bug")
		return
	}
	var event Event
	event, err = loadEvent(*eventRow)
	if err != nil {
		return
	}
	var ends *timestamppb.Timestamp
	if row.Ends.Valid {
		ends = timestamppb.New(row.Ends.Time)
	}
	workshop = Workshop{
		Id:           row.ID,
		WorkshopInfo: &info,
		Starts:       timestamppb.New(row.Starts),
		Ends:         ends,
		BelongsTo:    &Workshop_Event{Event: &event},
	}
	return
}

var (
	ErrEventDoesNotExist    = errors.New("event does not exist")
	ErrWorkshopDoesNotExist = errors.New("workshop does not exist")
	ErrRetrieveWorkshopList = errors.New("retrieving workshop list failed")
)
