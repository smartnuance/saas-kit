package event

//go:generate sqlboiler --config sqlboiler.toml psql
//go:generate mockgen -destination db_mock.go -package event . DBAPI

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/rs/xid"
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
	CreateWorkshop(ctx context.Context, body *CreateWorkshopBody, eventID string) (workshop *m.Workshop, err error)
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

func (db *dbAPI) CreateWorkshop(ctx context.Context, body *CreateWorkshopBody, eventID string) (workshop *m.Workshop, err error) {
	var info types.JSON
	info, err = json.Marshal(body.WorkshopInfo)
	if err != nil {
		return
	}
	workshop = &m.Workshop{
		ID:      xid.New().String(),
		Info:    info,
		Starts:  body.Starts,
		Ends:    null.TimeFrom(body.Ends),
		EventID: eventID,
	}
	err = workshop.Upsert(ctx, db.DB, true, boil.None().Cols, boil.Infer(), boil.Infer())
	if err != nil {
		return
	}
	return
}

// func (db *dbAPI) CreateUser(ctx context.Context, tx *sql.Tx, name, email string, passwordHash []byte) (user *m.User, err error) {
// 	user = &m.User{
// 		ID:       xid.New().String(),
// 		Name:     null.StringFrom(name),
// 		Email:    email,
// 		Password: passwordHash,
// 	}
// 	err = user.Insert(ctx, tx, boil.Infer())
// 	if err != nil {
// 		return
// 	}
// 	return
// }

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
