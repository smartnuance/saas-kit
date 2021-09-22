// Code generated by SQLBoiler 4.6.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package dbmodels

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/v4/types"
	"github.com/volatiletech/strmangle"
)

// Workshop is an object representing the database table.
type Workshop struct {
	ID           string     `boil:"id" json:"id" toml:"id" yaml:"id"`
	Info         types.JSON `boil:"info" json:"info" toml:"info" yaml:"info"`
	Starts       time.Time  `boil:"starts" json:"starts" toml:"starts" yaml:"starts"`
	Ends         null.Time  `boil:"ends" json:"ends,omitempty" toml:"ends" yaml:"ends,omitempty"`
	EventID      string     `boil:"event_id" json:"event_id" toml:"event_id" yaml:"event_id"`
	Participants types.JSON `boil:"participants" json:"participants" toml:"participants" yaml:"participants"`
	CreatedAt    time.Time  `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt    time.Time  `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	DeletedAt    null.Time  `boil:"deleted_at" json:"deleted_at,omitempty" toml:"deleted_at" yaml:"deleted_at,omitempty"`

	R *workshopR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L workshopL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var WorkshopColumns = struct {
	ID           string
	Info         string
	Starts       string
	Ends         string
	EventID      string
	Participants string
	CreatedAt    string
	UpdatedAt    string
	DeletedAt    string
}{
	ID:           "id",
	Info:         "info",
	Starts:       "starts",
	Ends:         "ends",
	EventID:      "event_id",
	Participants: "participants",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
}

var WorkshopTableColumns = struct {
	ID           string
	Info         string
	Starts       string
	Ends         string
	EventID      string
	Participants string
	CreatedAt    string
	UpdatedAt    string
	DeletedAt    string
}{
	ID:           "workshops.id",
	Info:         "workshops.info",
	Starts:       "workshops.starts",
	Ends:         "workshops.ends",
	EventID:      "workshops.event_id",
	Participants: "workshops.participants",
	CreatedAt:    "workshops.created_at",
	UpdatedAt:    "workshops.updated_at",
	DeletedAt:    "workshops.deleted_at",
}

// Generated where

var WorkshopWhere = struct {
	ID           whereHelperstring
	Info         whereHelpertypes_JSON
	Starts       whereHelpertime_Time
	Ends         whereHelpernull_Time
	EventID      whereHelperstring
	Participants whereHelpertypes_JSON
	CreatedAt    whereHelpertime_Time
	UpdatedAt    whereHelpertime_Time
	DeletedAt    whereHelpernull_Time
}{
	ID:           whereHelperstring{field: "\"event\".\"workshops\".\"id\""},
	Info:         whereHelpertypes_JSON{field: "\"event\".\"workshops\".\"info\""},
	Starts:       whereHelpertime_Time{field: "\"event\".\"workshops\".\"starts\""},
	Ends:         whereHelpernull_Time{field: "\"event\".\"workshops\".\"ends\""},
	EventID:      whereHelperstring{field: "\"event\".\"workshops\".\"event_id\""},
	Participants: whereHelpertypes_JSON{field: "\"event\".\"workshops\".\"participants\""},
	CreatedAt:    whereHelpertime_Time{field: "\"event\".\"workshops\".\"created_at\""},
	UpdatedAt:    whereHelpertime_Time{field: "\"event\".\"workshops\".\"updated_at\""},
	DeletedAt:    whereHelpernull_Time{field: "\"event\".\"workshops\".\"deleted_at\""},
}

// WorkshopRels is where relationship names are stored.
var WorkshopRels = struct {
	Event string
}{
	Event: "Event",
}

// workshopR is where relationships are stored.
type workshopR struct {
	Event *Event `boil:"Event" json:"Event" toml:"Event" yaml:"Event"`
}

// NewStruct creates a new relationship struct
func (*workshopR) NewStruct() *workshopR {
	return &workshopR{}
}

// workshopL is where Load methods for each relationship are stored.
type workshopL struct{}

var (
	workshopAllColumns            = []string{"id", "info", "starts", "ends", "event_id", "participants", "created_at", "updated_at", "deleted_at"}
	workshopColumnsWithoutDefault = []string{"id", "info", "starts", "ends", "event_id", "participants", "deleted_at"}
	workshopColumnsWithDefault    = []string{"created_at", "updated_at"}
	workshopPrimaryKeyColumns     = []string{"id"}
)

type (
	// WorkshopSlice is an alias for a slice of pointers to Workshop.
	// This should almost always be used instead of []Workshop.
	WorkshopSlice []*Workshop
	// WorkshopHook is the signature for custom Workshop hook methods
	WorkshopHook func(context.Context, boil.ContextExecutor, *Workshop) error

	workshopQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	workshopType                 = reflect.TypeOf(&Workshop{})
	workshopMapping              = queries.MakeStructMapping(workshopType)
	workshopPrimaryKeyMapping, _ = queries.BindMapping(workshopType, workshopMapping, workshopPrimaryKeyColumns)
	workshopInsertCacheMut       sync.RWMutex
	workshopInsertCache          = make(map[string]insertCache)
	workshopUpdateCacheMut       sync.RWMutex
	workshopUpdateCache          = make(map[string]updateCache)
	workshopUpsertCacheMut       sync.RWMutex
	workshopUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var workshopBeforeInsertHooks []WorkshopHook
var workshopBeforeUpdateHooks []WorkshopHook
var workshopBeforeDeleteHooks []WorkshopHook
var workshopBeforeUpsertHooks []WorkshopHook

var workshopAfterInsertHooks []WorkshopHook
var workshopAfterSelectHooks []WorkshopHook
var workshopAfterUpdateHooks []WorkshopHook
var workshopAfterDeleteHooks []WorkshopHook
var workshopAfterUpsertHooks []WorkshopHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Workshop) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range workshopBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Workshop) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range workshopBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Workshop) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range workshopBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Workshop) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range workshopBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Workshop) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range workshopAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Workshop) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range workshopAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Workshop) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range workshopAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Workshop) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range workshopAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Workshop) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range workshopAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddWorkshopHook registers your hook function for all future operations.
func AddWorkshopHook(hookPoint boil.HookPoint, workshopHook WorkshopHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		workshopBeforeInsertHooks = append(workshopBeforeInsertHooks, workshopHook)
	case boil.BeforeUpdateHook:
		workshopBeforeUpdateHooks = append(workshopBeforeUpdateHooks, workshopHook)
	case boil.BeforeDeleteHook:
		workshopBeforeDeleteHooks = append(workshopBeforeDeleteHooks, workshopHook)
	case boil.BeforeUpsertHook:
		workshopBeforeUpsertHooks = append(workshopBeforeUpsertHooks, workshopHook)
	case boil.AfterInsertHook:
		workshopAfterInsertHooks = append(workshopAfterInsertHooks, workshopHook)
	case boil.AfterSelectHook:
		workshopAfterSelectHooks = append(workshopAfterSelectHooks, workshopHook)
	case boil.AfterUpdateHook:
		workshopAfterUpdateHooks = append(workshopAfterUpdateHooks, workshopHook)
	case boil.AfterDeleteHook:
		workshopAfterDeleteHooks = append(workshopAfterDeleteHooks, workshopHook)
	case boil.AfterUpsertHook:
		workshopAfterUpsertHooks = append(workshopAfterUpsertHooks, workshopHook)
	}
}

// One returns a single workshop record from the query.
func (q workshopQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Workshop, error) {
	o := &Workshop{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "dbmodels: failed to execute a one query for workshops")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Workshop records from the query.
func (q workshopQuery) All(ctx context.Context, exec boil.ContextExecutor) (WorkshopSlice, error) {
	var o []*Workshop

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "dbmodels: failed to assign all query results to Workshop slice")
	}

	if len(workshopAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Workshop records in the query.
func (q workshopQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "dbmodels: failed to count workshops rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q workshopQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "dbmodels: failed to check if workshops exists")
	}

	return count > 0, nil
}

// Event pointed to by the foreign key.
func (o *Workshop) Event(mods ...qm.QueryMod) eventQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.EventID),
		qmhelper.WhereIsNull("deleted_at"),
	}

	queryMods = append(queryMods, mods...)

	query := Events(queryMods...)
	queries.SetFrom(query.Query, "\"event\".\"events\"")

	return query
}

// LoadEvent allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (workshopL) LoadEvent(ctx context.Context, e boil.ContextExecutor, singular bool, maybeWorkshop interface{}, mods queries.Applicator) error {
	var slice []*Workshop
	var object *Workshop

	if singular {
		object = maybeWorkshop.(*Workshop)
	} else {
		slice = *maybeWorkshop.(*[]*Workshop)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &workshopR{}
		}
		args = append(args, object.EventID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &workshopR{}
			}

			for _, a := range args {
				if a == obj.EventID {
					continue Outer
				}
			}

			args = append(args, obj.EventID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`event.events`),
		qm.WhereIn(`event.events.id in ?`, args...),
		qmhelper.WhereIsNull(`event.events.deleted_at`),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Event")
	}

	var resultSlice []*Event
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Event")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for events")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for events")
	}

	if len(workshopAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(ctx, e); err != nil {
				return err
			}
		}
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Event = foreign
		if foreign.R == nil {
			foreign.R = &eventR{}
		}
		foreign.R.Workshops = append(foreign.R.Workshops, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.EventID == foreign.ID {
				local.R.Event = foreign
				if foreign.R == nil {
					foreign.R = &eventR{}
				}
				foreign.R.Workshops = append(foreign.R.Workshops, local)
				break
			}
		}
	}

	return nil
}

// SetEvent of the workshop to the related item.
// Sets o.R.Event to related.
// Adds o to related.R.Workshops.
func (o *Workshop) SetEvent(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Event) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"event\".\"workshops\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"event_id"}),
		strmangle.WhereClause("\"", "\"", 2, workshopPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.EventID = related.ID
	if o.R == nil {
		o.R = &workshopR{
			Event: related,
		}
	} else {
		o.R.Event = related
	}

	if related.R == nil {
		related.R = &eventR{
			Workshops: WorkshopSlice{o},
		}
	} else {
		related.R.Workshops = append(related.R.Workshops, o)
	}

	return nil
}

// Workshops retrieves all the records using an executor.
func Workshops(mods ...qm.QueryMod) workshopQuery {
	mods = append(mods, qm.From("\"event\".\"workshops\""), qmhelper.WhereIsNull("\"event\".\"workshops\".\"deleted_at\""))
	return workshopQuery{NewQuery(mods...)}
}

// FindWorkshop retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindWorkshop(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*Workshop, error) {
	workshopObj := &Workshop{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"event\".\"workshops\" where \"id\"=$1 and \"deleted_at\" is null", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, workshopObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "dbmodels: unable to select from workshops")
	}

	if err = workshopObj.doAfterSelectHooks(ctx, exec); err != nil {
		return workshopObj, err
	}

	return workshopObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Workshop) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("dbmodels: no workshops provided for insertion")
	}

	var err error
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		if o.UpdatedAt.IsZero() {
			o.UpdatedAt = currTime
		}
	}

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(workshopColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	workshopInsertCacheMut.RLock()
	cache, cached := workshopInsertCache[key]
	workshopInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			workshopAllColumns,
			workshopColumnsWithDefault,
			workshopColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(workshopType, workshopMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(workshopType, workshopMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"event\".\"workshops\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"event\".\"workshops\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "dbmodels: unable to insert into workshops")
	}

	if !cached {
		workshopInsertCacheMut.Lock()
		workshopInsertCache[key] = cache
		workshopInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Workshop.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Workshop) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		o.UpdatedAt = currTime
	}

	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	workshopUpdateCacheMut.RLock()
	cache, cached := workshopUpdateCache[key]
	workshopUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			workshopAllColumns,
			workshopPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("dbmodels: unable to update workshops, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"event\".\"workshops\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, workshopPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(workshopType, workshopMapping, append(wl, workshopPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "dbmodels: unable to update workshops row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "dbmodels: failed to get rows affected by update for workshops")
	}

	if !cached {
		workshopUpdateCacheMut.Lock()
		workshopUpdateCache[key] = cache
		workshopUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q workshopQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "dbmodels: unable to update all for workshops")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "dbmodels: unable to retrieve rows affected for workshops")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o WorkshopSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("dbmodels: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), workshopPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"event\".\"workshops\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, workshopPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "dbmodels: unable to update all in workshop slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "dbmodels: unable to retrieve rows affected all in update all workshop")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Workshop) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("dbmodels: no workshops provided for upsert")
	}
	if !boil.TimestampsAreSkipped(ctx) {
		currTime := time.Now().In(boil.GetLocation())

		if o.CreatedAt.IsZero() {
			o.CreatedAt = currTime
		}
		o.UpdatedAt = currTime
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(workshopColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	workshopUpsertCacheMut.RLock()
	cache, cached := workshopUpsertCache[key]
	workshopUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			workshopAllColumns,
			workshopColumnsWithDefault,
			workshopColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			workshopAllColumns,
			workshopPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("dbmodels: unable to upsert workshops, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(workshopPrimaryKeyColumns))
			copy(conflict, workshopPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"event\".\"workshops\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(workshopType, workshopMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(workshopType, workshopMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "dbmodels: unable to upsert workshops")
	}

	if !cached {
		workshopUpsertCacheMut.Lock()
		workshopUpsertCache[key] = cache
		workshopUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Workshop record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Workshop) Delete(ctx context.Context, exec boil.ContextExecutor, hardDelete bool) (int64, error) {
	if o == nil {
		return 0, errors.New("dbmodels: no Workshop provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	var (
		sql  string
		args []interface{}
	)
	if hardDelete {
		args = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), workshopPrimaryKeyMapping)
		sql = "DELETE FROM \"event\".\"workshops\" WHERE \"id\"=$1"
	} else {
		currTime := time.Now().In(boil.GetLocation())
		o.DeletedAt = null.TimeFrom(currTime)
		wl := []string{"deleted_at"}
		sql = fmt.Sprintf("UPDATE \"event\".\"workshops\" SET %s WHERE \"id\"=$2",
			strmangle.SetParamNames("\"", "\"", 1, wl),
		)
		valueMapping, err := queries.BindMapping(workshopType, workshopMapping, append(wl, workshopPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
		args = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), valueMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "dbmodels: unable to delete from workshops")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "dbmodels: failed to get rows affected by delete for workshops")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q workshopQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor, hardDelete bool) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("dbmodels: no workshopQuery provided for delete all")
	}

	if hardDelete {
		queries.SetDelete(q.Query)
	} else {
		currTime := time.Now().In(boil.GetLocation())
		queries.SetUpdate(q.Query, M{"deleted_at": currTime})
	}

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "dbmodels: unable to delete all from workshops")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "dbmodels: failed to get rows affected by deleteall for workshops")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o WorkshopSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor, hardDelete bool) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(workshopBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var (
		sql  string
		args []interface{}
	)
	if hardDelete {
		for _, obj := range o {
			pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), workshopPrimaryKeyMapping)
			args = append(args, pkeyArgs...)
		}
		sql = "DELETE FROM \"event\".\"workshops\" WHERE " +
			strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, workshopPrimaryKeyColumns, len(o))
	} else {
		currTime := time.Now().In(boil.GetLocation())
		for _, obj := range o {
			pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), workshopPrimaryKeyMapping)
			args = append(args, pkeyArgs...)
			obj.DeletedAt = null.TimeFrom(currTime)
		}
		wl := []string{"deleted_at"}
		sql = fmt.Sprintf("UPDATE \"event\".\"workshops\" SET %s WHERE "+
			strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 2, workshopPrimaryKeyColumns, len(o)),
			strmangle.SetParamNames("\"", "\"", 1, wl),
		)
		args = append([]interface{}{currTime}, args...)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "dbmodels: unable to delete all from workshop slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "dbmodels: failed to get rows affected by deleteall for workshops")
	}

	if len(workshopAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Workshop) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindWorkshop(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *WorkshopSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := WorkshopSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), workshopPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"event\".\"workshops\".* FROM \"event\".\"workshops\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, workshopPrimaryKeyColumns, len(*o)) +
		"and \"deleted_at\" is null"

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "dbmodels: unable to reload all in WorkshopSlice")
	}

	*o = slice

	return nil
}

// WorkshopExists checks if the Workshop row exists.
func WorkshopExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"event\".\"workshops\" where \"id\"=$1 and \"deleted_at\" is null limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "dbmodels: unable to check if workshops exists")
	}

	return exists, nil
}