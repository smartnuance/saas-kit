package auth

import (
	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/smartnuance/saas-kit/pkg/lib/roles"
)

// RevokeBody describes the user/instance to revoke tokens for
type RevokeBody struct {
	Email       string `json:"email"`
	InstanceURL string `json:"url"`
}

func (s *Service) Revoke(ctx *gin.Context) error {
	var body RevokeBody
	err := ctx.ShouldBind(&body)
	if err != nil {
		return err
	}

	var userID string
	if len(body.Email) > 0 {
		user, err := s.DBAPI.FindUserByEmail(ctx, body.Email)
		if err != nil {
			return err
		}
		userID = user.ID
	} else {
		userID, err = roles.User(ctx)
		if err != nil {
			return err
		}
	}

	var instanceID string
	if len(body.InstanceURL) > 0 {
		instance, err := s.DBAPI.GetInstance(ctx, body.InstanceURL)
		if err != nil {
			return err
		}
		instanceID = instance.ID
	} else {
		// fallback to default instance from headers
		_, instanceID, err = roles.ApplyHeaders(ctx)
		if err != nil {
			return err
		}
	}

	// Check permission to revoke token for potentially different user
	if !(roles.CanActAs(ctx, userID) ||
		(roles.CanActFor(ctx, instanceID) && roles.CanActIn(ctx, roles.RoleInstanceAdmin))) {
		return errors.WithStack(ErrUnauthorized)
	}

	profile, err := s.DBAPI.GetProfile(ctx, userID, instanceID)
	if err != nil {
		return errors.WithStack(ErrProfileDoesNotExist)
	}

	_, err = s.DBAPI.DeleteToken(ctx, profile.ID)
	if err != nil {
		return err
	}
	return nil
}

// RevokeAllBody describes the user to revoke all tokens for
type RevokeAllBody struct {
	Email string `json:"email"`
}

func (s *Service) RevokeAll(ctx *gin.Context) error {
	var body RevokeAllBody
	err := ctx.ShouldBind(&body)
	if err != nil {
		return err
	}

	var userID string
	if len(body.Email) > 0 {
		user, err := s.DBAPI.FindUserByEmail(ctx, body.Email)
		if err != nil {
			return err
		}
		userID = user.ID
	} else {
		userID, err = roles.User(ctx)
		if err != nil {
			return err
		}
	}

	// Check permission to revoke token for potentially different user
	if !(roles.CanActAs(ctx, userID) ||
		roles.CanActIn(ctx, roles.RoleSuperAdmin)) {
		return errors.WithStack(ErrUnauthorized)
	}

	_, err = s.DBAPI.DeleteAllTokens(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

var (
	ErrMissingRevokeEmail = errors.New("missing user id")
	ErrTokenRevoked       = errors.New("refresh token was revoked and is no longer valid")
)
