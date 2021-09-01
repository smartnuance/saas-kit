package auth

//go:generate sqlboiler --config sqlboiler.toml psql

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

func FindUserByEmail(ctx context.Context, email string) (*m.User, error) {
	user, err := m.Users(m.UserWhere.Email.EQ(email)).OneG(ctx)
	if err == sql.ErrNoRows {
		// transform sql error in specific error of login context
		return nil, errors.WithStack(ErrUserDoesNotExist)
	}
	return user, err
}

func GetInstance(ctx context.Context, instanceURL string) (instance *m.Instance, err error) {
	instance, err = m.Instances(m.InstanceWhere.URL.EQ(instanceURL)).OneG(ctx)
	if err == sql.ErrNoRows {
		// transform sql error in specific error of login context
		err = errors.WithStack(ErrInstanceDoesNotExist)
		return
	}
	return instance, err
}

func GetProfile(ctx context.Context, userID, instanceID int64) (profile *m.Profile, err error) {
	where := &m.ProfileWhere
	profile, err = m.Profiles(where.UserID.EQ(userID), where.InstanceID.EQ(instanceID)).OneG(ctx)
	if err == sql.ErrNoRows {
		// transform sql error in specific error of login context
		err = errors.WithStack(ErrProfileDoesNotExist)
		return
	}
	return profile, err
}

func GetUserAndProfile(ctx context.Context, userID int64, instanceURL string) (user *m.User, profile *m.Profile, err error) {
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

func CreateProfile(ctx context.Context, tx *sql.Tx, instanceID int64, user *m.User, role string) (profile *m.Profile, err error) {
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

func CreateUser(ctx context.Context, tx *sql.Tx, name, email string, passwordHash []byte) (user *m.User, err error) {
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

func DeleteUser(ctx context.Context, userID int64) error {
	_, err := m.Users(m.UserWhere.ID.EQ(userID)).DeleteAllG(ctx, false)
	return err
}

func SaveToken(ctx context.Context, profile *m.Profile, token string, expiresAt time.Time) error {
	t := m.Token{
		UserID:    profile.UserID,
		ProfileID: profile.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	return t.InsertG(ctx, boil.Infer())
}

func DeleteToken(ctx context.Context, profileID int64) (int64, error) {
	where := &m.TokenWhere
	numDeleted, err := m.Tokens(
		where.ProfileID.EQ(profileID),
	).DeleteAllG(ctx)
	return numDeleted, err
}

func DeleteAllTokens(ctx context.Context, userID int64) (int64, error) {
	where := &m.TokenWhere
	numDeleted, err := m.Tokens(
		where.UserID.EQ(userID),
	).DeleteAllG(ctx)
	return numDeleted, err
}
