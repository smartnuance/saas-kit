package tokens

import (
	"crypto/rsa"

	"github.com/pkg/errors"
	"github.com/smartnuance/saas-kit/pkg/lib/roles"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
)

const (
	AccessPurpose  = "access"
	RefreshPurpose = "refresh"
)

// AccessTokenClaims contain temporary authorization information.
type AccessTokenClaims struct {
	Purpose  string `json:"purp"`
	Role     string `json:"role"`
	Instance string `json:"inst"`
	jwt.StandardClaims
}

// RefreshTokenClaims contain everything necessary to recreate an accesstoken,
// i.e. identify the right profile to load role and user meta information from.
type RefreshTokenClaims struct {
	Purpose  string `json:"purp"`
	Instance string `json:"inst"`
	jwt.StandardClaims
}

const BearerSchema = "Bearer "

// AuthorizeJWT creates a middleware that checks the presence and validity of the authorization header.
// If this middleware is installed on an endpoint, the authorization header is required.
// When the header is present and the access token (JWT) inside is valid, user, role and instance are set to context.
// The middleware creation is parameterized by service specifics.
func AuthorizeJWT(validationKey *rsa.PublicKey, issuer, audience string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if len(authHeader) <= len(BearerSchema) {
			log.Error().Msgf("missing/invalid authorization header, needs to start with '%s'", BearerSchema)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenString := authHeader[len(BearerSchema):]
		var claims AccessTokenClaims
		err := CheckAccessToken(tokenString, &claims, validationKey, issuer, audience)
		if err != nil {
			log.Error().Err(err).Msg("")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// set default context from JWT attributes
		ctx.Set(roles.UserKey, claims.Subject)      // acting subject (immutable)
		ctx.Set(roles.InstanceKey, claims.Instance) // instance (switchable by super admins onld)
		ctx.Set(roles.RoleKey, claims.Role)         // role (switchable if permission to)

		// order matters: first check if default JWT role allows for instance switch if header is present
		switchInstance := ctx.GetHeader(roles.InstanceHeader)
		if switchInstance != "" && switchInstance != claims.Instance {
			if !roles.CanActIn(ctx, roles.RoleSuperAdmin) {
				log.Error().Err(err).Msg("")
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			ctx.Set(roles.InstanceKey, switchInstance)
		}

		switchRole := ctx.GetHeader(roles.RoleHeader)
		err = roles.SwitchTo(ctx, switchRole)
		if err != nil {
			log.Error().Err(err).Msg("")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

func CheckAccessToken(tokenStr string, claims *AccessTokenClaims, validationKey *rsa.PublicKey, issuer, audience string) error {
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodRSA); !isvalid {
			return nil, fmt.Errorf("invalid token signing method: %s", token.Header["alg"])
		}
		return validationKey, nil
	})
	if err != nil {
		return errors.Wrap(err, "invalid token")
	}
	if !token.Valid {
		return errors.Wrap(err, "invalid token claims")
	}
	if claims.Purpose != AccessPurpose {
		return errors.Wrap(err, "invalid token purpose")
	}
	ok := claims.VerifyIssuer(issuer, true)
	if !ok {
		return errors.New("invalid token issuer")
	}
	ok = claims.VerifyAudience(audience, true)
	if !ok {
		return errors.New("invalid token audience")
	}
	return nil
}

func CheckRefreshToken(tokenStr string, claims *RefreshTokenClaims, validationKey *rsa.PublicKey, issuer, audience string) error {
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodRSA); !isvalid {
			return nil, fmt.Errorf("invalid token signing method: %s", token.Header["alg"])
		}
		return validationKey, nil
	})
	if err != nil {
		return errors.Wrap(err, "invalid token")
	}
	if !token.Valid {
		return errors.Wrap(err, "invalid token claims")
	}
	if claims.Purpose != RefreshPurpose {
		return errors.Errorf("invalid token purpose %s", claims.Purpose)
	}
	ok := claims.VerifyIssuer(issuer, true)
	if !ok {
		return errors.New("invalid token issuer")
	}
	ok = claims.VerifyAudience(audience, true)
	if !ok {
		return errors.New("invalid token audience")
	}
	return nil
}
