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

type AccessTokenClaims struct {
	User     bool   `json:"user"`
	Role     string `json:"role"`
	Instance int    `json:"instance"`
	jwt.StandardClaims
}

const BearerSchema = "Bearer "

func AuthorizeJWT(validationKey *rsa.PublicKey) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if len(authHeader) <= len(BearerSchema) {
			log.Error().Msgf("missing/invalid authorization header, needs to start with '%s'", BearerSchema)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenString := authHeader[len(BearerSchema):]
		var claims AccessTokenClaims
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			if _, isvalid := token.Method.(*jwt.SigningMethodRSA); !isvalid {
				return nil, fmt.Errorf("invalid token signing method: %s", token.Header["alg"])
			}
			return validationKey, nil
		})
		if err != nil {
			log.Error().Err(err).Msg("error parsing authorization header")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			log.Error().Msg("invalid token in authorization header")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("role", claims.Role)
		ctx.Set("instance", claims.Instance)
	}
}

var (
	ErrMissingToken = errors.New("missing token")
	ErrInvalidToken = errors.New("invalid token")
)
