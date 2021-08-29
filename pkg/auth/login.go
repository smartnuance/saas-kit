package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	m "github.com/smartnuance/saas-kit/pkg/auth/dbmodels"
	"github.com/smartnuance/saas-kit/pkg/lib/roles"
	"github.com/smartnuance/saas-kit/pkg/lib/tokens"
	"golang.org/x/crypto/bcrypt"
)

// CredentialsBody describes the login credentials
type CredentialsBody struct {
	InstanceURL string `json:"url"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

func (s *Service) Login(ctx *gin.Context) (accessToken, refreshToken string, err error) {
	var body CredentialsBody
	err = ctx.ShouldBind(&body)
	if err != nil {
		err = errors.WithStack(ErrMissingCredentials)
		return
	}
	var user *m.User
	user, err = loginWithCredentials(ctx, body.Email, body.Password)
	if err != nil {
		return
	}

	var instance *m.Instance
	instance, err = GetInstance(ctx, body.InstanceURL)
	if err != nil {
		return
	}

	var expiresAt time.Time
	refreshToken, expiresAt, err = s.TokenAPI.GenerateRefreshToken(int(user.ID), int(instance.ID), true)
	if err != nil {
		return
	}

	profile, err := GetProfile(ctx, user.ID, instance.ID)
	if err != nil {
		err = errors.WithStack(ErrProfileDoesNotExist)
		return
	}

	var role string
	if profile.Role.Valid {
		role = profile.Role.String
	} else {
		role = roles.NoRole
	}
	accessToken, err = s.TokenAPI.GenerateAccessToken(int(user.ID), int(instance.ID), true, role)
	if err != nil {
		return
	}

	err = SaveToken(ctx, profile, refreshToken, expiresAt)
	if err != nil {
		return
	}

	return accessToken, refreshToken, nil
}

func loginWithCredentials(ctx *gin.Context, email string, password string) (*m.User, error) {
	user, err := FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return nil, errors.WithStack(ErrInvalidCredentials)
	}
	return user, nil
}

// RefreshTokenBody describes the refresh body
type RefreshTokenBody struct {
	RefreshToken string `json:"refreshToken"`
}

func (s *Service) Refresh(ctx *gin.Context) (string, error) {
	var body RefreshTokenBody
	err := ctx.ShouldBind(&body)
	if err != nil {
		return "", errors.WithStack(ErrMissingRefreshToken)
	}

	var claims tokens.RefreshTokenClaims
	token, err := jwt.ParseWithClaims(body.RefreshToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodRSA); !isvalid {
			return nil, fmt.Errorf("invalid token signing method: %s", token.Header["alg"])
		}
		return s.TokenAPI.ValidationKey, nil
	})
	if err != nil {
		return "", errors.WithStack(ErrTokenInvalid)
	}
	if !token.Valid {
		return "", errors.WithStack(ErrTokenInvalid)
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return "", errors.WithStack(ErrInvalidUserID)
	}

	profile, err := GetProfile(ctx, int64(userID), int64(claims.Instance))
	if err != nil {
		return "", errors.WithStack(ErrProfileDoesNotExist)
	}
	var role string
	if profile.Role.Valid {
		role = profile.Role.String
	} else {
		role = roles.NoRole
	}

	return s.TokenAPI.GenerateAccessToken(userID, claims.Instance, true, role)
}

func (s *Service) Revoke(ctx *gin.Context) error {
	_userID := ctx.Param("user_id")
	if len(_userID) == 0 {
		return errors.WithStack(ErrMissingUserID)
	}
	userID, err := strconv.Atoi(_userID)
	if err != nil {
		return errors.WithStack(ErrInvalidUserID)
	}

	_, instanceID, err := roles.FromHeaders(ctx)
	if err != nil {
		return err
	}

	// Check permission to revoke token for potentially different user
	if !(roles.CanActFor(ctx, instanceID) || roles.CanActIn(ctx, roles.RoleSuperAdmin)) {
		return ErrUnauthorized
	}

	numDeleted, err := DeleteToken(ctx, int64(userID))
	if err != nil {
		return err
	}
	if numDeleted == 0 {
		return errors.WithStack(ErrTokenNotFound)
	}
	return nil
}

var (
	ErrUnauthorized         = errors.New("role insufficient to act on desired instances")
	ErrMissingCredentials   = errors.New("missing credentials, email and password have to be provided")
	ErrInvalidCredentials   = errors.New("invalid credentials, email/password combination wrong")
	ErrMissingUserID        = errors.New("missing user id")
	ErrInvalidUserID        = errors.New("invalid user id provided")
	ErrMissingRefreshToken  = errors.New("missing refresh token in JSON body")
	ErrTokenInvalid         = errors.New("token invalid")
	ErrTokenNotFound        = errors.New("token not found")
	ErrUserDoesNotExist     = errors.New("user does not exist")
	ErrInstanceDoesNotExist = errors.New("instance does not exist")
	ErrProfileDoesNotExist  = errors.New("profile does not exist")
)
