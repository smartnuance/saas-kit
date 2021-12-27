package auth

import (
	"time"

	"github.com/friendsofgo/errors"

	"github.com/gin-gonic/gin"
	m "github.com/smartnuance/saas-kit/pkg/auth/dbmodels"
	"github.com/smartnuance/saas-kit/pkg/lib/roles"
	"github.com/smartnuance/saas-kit/pkg/lib/tokens"
	"golang.org/x/crypto/bcrypt"
)

// CredentialsBody describes the login credentials
type CredentialsBody struct {
	InstanceURL string `json:"instance"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

func (s *Service) Login(ctx *gin.Context) (accessToken, refreshToken string, role roles.Role, err error) {
	var body CredentialsBody
	err = ctx.ShouldBind(&body)
	if err != nil {
		err = errors.WithStack(ErrMissingCredentials)
		return
	}
	var user *m.User
	user, err = s.loginWithCredentials(ctx, body.Email, body.Password)
	if err != nil {
		return
	}

	var instance *m.Instance
	instance, err = s.DBAPI.GetInstance(ctx, body.InstanceURL)
	if err != nil {
		return
	}

	var expiresAt time.Time
	refreshToken, expiresAt, err = s.TokenAPI.GenerateRefreshToken(user.ID, instance.ID)
	if err != nil {
		return
	}

	profile, err := s.DBAPI.GetProfile(ctx, user.ID, instance.ID)
	if err != nil {
		err = errors.WithStack(ErrProfileDoesNotExist)
		return
	}

	if profile.Role.Valid {
		role = roles.Role(profile.Role.String)
	} else {
		role = roles.NoRole
	}
	accessToken, err = s.TokenAPI.GenerateAccessToken(user.ID, instance.ID, role)
	if err != nil {
		return
	}

	err = s.DBAPI.SaveToken(ctx, profile, refreshToken, expiresAt)
	if err != nil {
		return
	}

	return
}

func (s *Service) loginWithCredentials(ctx *gin.Context, email string, password string) (*m.User, error) {
	user, err := s.DBAPI.FindUserByEmail(ctx, email)
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
	err = tokens.CheckRefreshToken(body.RefreshToken, &claims, s.TokenAPI.ValidationKey, s.Issuer, s.Audience)
	if err != nil {
		return "", errors.WithStack(errors.Wrap(err, ErrTokenInvalid.Error()))
	}

	userID := claims.Subject
	profile, err := s.DBAPI.GetProfile(ctx, userID, claims.Instance)
	if err != nil {
		return "", errors.WithStack(ErrProfileDoesNotExist)
	}

	// check if revoked in the meanwhile
	ok, err := s.DBAPI.HasToken(ctx, userID, profile.ID, body.RefreshToken)
	if err != nil {
		return "", errors.WithStack(err)
	}
	if !ok {
		return "", errors.WithStack(ErrTokenRevoked)
	}

	var role roles.Role
	if profile.Role.Valid {
		role = roles.Role(profile.Role.String)
	} else {
		role = roles.NoRole
	}

	return s.TokenAPI.GenerateAccessToken(userID, claims.Instance, role)
}

var (
	ErrMissingCredentials   = errors.New("missing credentials, email and password have to be provided")
	ErrInvalidCredentials   = errors.New("invalid credentials, email/password combination wrong")
	ErrMissingRefreshToken  = errors.New("missing refresh token in JSON body")
	ErrTokenInvalid         = errors.New("token invalid")
	ErrTokenNotFound        = errors.New("token not found")
	ErrUserDoesNotExist     = errors.New("user does not exist")
	ErrInstanceDoesNotExist = errors.New("instance does not exist")
	ErrProfileDoesNotExist  = errors.New("profile does not exist")
)
