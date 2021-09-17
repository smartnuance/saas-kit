package tokens

import (
	"crypto/rsa"

	"github.com/pkg/errors"

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
	User     bool   `json:"user"`
	Role     string `json:"role"`
	Instance string    `json:"inst"`
	jwt.StandardClaims
}

// RefreshTokenClaims contain everything necessary to recreate an accesstoken,
// i.e. identify the right profile to load role and user meta information from.
type RefreshTokenClaims struct {
	Purpose  string `json:"purp"`
	User     bool   `json:"user"`
	Instance string    `json:"inst"`
	jwt.StandardClaims
}

const BearerSchema = "Bearer "

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

		ctx.Set("user", claims.User)
		ctx.Set("role", claims.Role)
		ctx.Set("instance", claims.Instance)
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
