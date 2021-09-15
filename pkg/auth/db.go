package auth

//go:generate sqlboiler --config sqlboiler.toml psql
//go:generate mockgen -destination db_mock.go -package auth . DBAPI

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
	m "github.com/smartnuance/saas-kit/pkg/auth/dbmodels"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	// . "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type DBAPI interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	Commit(tx *sql.Tx) error
	Rollback(tx *sql.Tx) error
	FindUserByEmail(ctx context.Context, email string) (*m.User, error)
	GetInstance(ctx context.Context, instanceURL string) (instance *m.Instance, err error)
	GetProfile(ctx context.Context, userID, instanceID int64) (profile *m.Profile, err error)
	GetUserAndProfile(ctx context.Context, userID int64, instanceURL string) (user *m.User, profile *m.Profile, err error)
	CreateProfile(ctx context.Context, tx *sql.Tx, instanceID int64, user *m.User, role string) (profile *m.Profile, err error)
	CreateUser(ctx context.Context, tx *sql.Tx, name, email string, passwordHash []byte) (user *m.User, err error)
	DeleteUser(ctx context.Context, userID int64) error
	SaveToken(ctx context.Context, profile *m.Profile, token string, expiresAt time.Time) error
	DeleteToken(ctx context.Context, profileID int64) (int64, error)
	DeleteAllTokens(ctx context.Context, userID int64) (int64, error)
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

func (db *dbAPI) FindUserByEmail(ctx context.Context, email string) (*m.User, error) {
	user, err := m.Users(m.UserWhere.Email.EQ(email)).OneG(ctx)
	if err == sql.ErrNoRows {
		// transform sql error in specific error of login context
		return nil, errors.WithStack(ErrUserDoesNotExist)
	}
	return user, err
}

func (db *dbAPI) GetInstance(ctx context.Context, instanceURL string) (instance *m.Instance, err error) {
	instance, err = m.Instances(m.InstanceWhere.URL.EQ(instanceURL)).OneG(ctx)
	if err == sql.ErrNoRows {
		// transform sql error in specific error of login context
		err = errors.WithStack(ErrInstanceDoesNotExist)
		return
	}
	return instance, err
}

func (db *dbAPI) GetProfile(ctx context.Context, userID, instanceID int64) (profile *m.Profile, err error) {
	where := &m.ProfileWhere
	profile, err = m.Profiles(where.UserID.EQ(userID), where.InstanceID.EQ(instanceID)).OneG(ctx)
	if err == sql.ErrNoRows {
		// transform sql error in specific error of login context
		err = errors.WithStack(ErrProfileDoesNotExist)
		return
	}
	return profile, err
}

func (db *dbAPI) GetUserAndProfile(ctx context.Context, userID int64, instanceURL string) (user *m.User, profile *m.Profile, err error) {
	profile, err = m.Profiles(m.ProfileWhere.UserID.EQ(userID)).OneG(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			// transform sql error in specific error of login context
			err = errors.WithStack(ErrProfileDoesNotExist)
			return
		}
		return
	}
	user, err = profile.User().OneG(ctx)
	if err != nil {
		return
	}
	return
}

func (db *dbAPI) CreateProfile(ctx context.Context, tx *sql.Tx, instanceID int64, user *m.User, role string) (profile *m.Profile, err error) {
	profile = &m.Profile{
		InstanceID: instanceID,
		UserID:     user.ID,
		Role:       null.StringFrom(role),
	}
	err = profile.Upsert(ctx, tx, true, boil.None().Cols, boil.Infer(), boil.Infer())
	if err != nil {
		return
	}
	return
}

func (db *dbAPI) CreateUser(ctx context.Context, tx *sql.Tx, name, email string, passwordHash []byte) (user *m.User, err error) {
	user = &m.User{
		Name:     null.StringFrom(name),
		Email:    email,
		Password: passwordHash,
	}
	err = user.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return
	}
	return
}

func (db *dbAPI) DeleteUser(ctx context.Context, userID int64) error {
	_, err := m.Users(m.UserWhere.ID.EQ(userID)).DeleteAllG(ctx, false)
	return err
}

func (db *dbAPI) SaveToken(ctx context.Context, profile *m.Profile, token string, expiresAt time.Time) error {
	t := m.Token{
		UserID:    profile.UserID,
		ProfileID: profile.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	return t.InsertG(ctx, boil.Infer())
}

func (db *dbAPI) DeleteToken(ctx context.Context, profileID int64) (int64, error) {
	where := &m.TokenWhere
	numDeleted, err := m.Tokens(
		where.ProfileID.EQ(profileID),
	).DeleteAllG(ctx)
	return numDeleted, err
}

func (db *dbAPI) DeleteAllTokens(ctx context.Context, userID int64) (int64, error) {
	where := &m.TokenWhere
	numDeleted, err := m.Tokens(
		where.UserID.EQ(userID),
	).DeleteAllG(ctx)
	return numDeleted, err
}
