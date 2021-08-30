package auth

import (
	"database/sql"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	m "github.com/smartnuance/saas-kit/pkg/auth/dbmodels"
	"github.com/smartnuance/saas-kit/pkg/lib/roles"
	"golang.org/x/crypto/bcrypt"
)

// SignupBody describes the signup body with desired credentials
type SignupBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *Service) Signup(ctx *gin.Context) (userID int, err error) {
	var body SignupBody
	err = ctx.ShouldBind(&body)
	if err != nil {
		return
	}

	var role string
	var instanceID int
	role, instanceID, err = roles.FromHeaders(ctx)
	if err != nil {
		return
	}

	return s.signup(ctx, instanceID, body, role)
}

func (s *Service) signup(ctx *gin.Context, instanceID int, body SignupBody, role string) (userID int, err error) {
	log.Debug().Msgf("Signup user %s with email %s to %d with role %s", body.Name, body.Email, instanceID, role)
	if len(body.Email) == 0 {
		err = ErrInvalidEmail
		return
	}
	if len(body.Password) == 0 {
		err = ErrInvalidPassword
		return
	}

	hashedPassword, err := hashAndSaltPassword(body.Password)
	if err != nil {
		return
	}

	// use a transaction to ensure user is only created with a valid profile
	var tx *sql.Tx
	tx, err = s.DB.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	var user *m.User
	user, err = CreateUser(ctx, tx, body.Name, body.Email, hashedPassword)
	if err != nil {
		return
	}

	_, err = CreateProfile(ctx, tx, int64(instanceID), user, role)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		errRollback := tx.Rollback()
		if errRollback != nil {
			errors.Wrapf(err, errRollback.Error())
		}
		return
	}

	return int(user.ID), nil
}

func hashAndSaltPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

var (
	ErrInvalidEmail    = errors.New("invalid user email provided")
	ErrInvalidPassword = errors.New("invalid user password provided")
)
