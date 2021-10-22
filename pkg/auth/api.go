package auth

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/smartnuance/saas-kit/pkg/lib/roles"
	"github.com/smartnuance/saas-kit/pkg/lib/tokens"
)

func router(s *Service) *gin.Engine {
	var router = gin.Default()

	config := cors.DefaultConfig()
	config.AddAllowMethods("PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS")
	config.AddAllowHeaders("Authorization")
	config.AddAllowHeaders(roles.RoleHeader)
	if s.release {
		config.AllowOriginFunc = func(origin string) bool {
			_, ok := s.AllowOrigins[origin]
			return ok
		}
	} else {
		config.AllowAllOrigins = true
	}
	router.Use(cors.New(config))

	api := router.Group("/")

	// without authorization middleware
	api.PUT("/signup", func(ctx *gin.Context) {
		SignupHandler(ctx, s)
	})
	api.POST("/login", func(ctx *gin.Context) {
		LoginHandler(ctx, s)
	})
	api.POST("/refresh", func(ctx *gin.Context) {
		RefreshHandler(ctx, s)
	})

	// with authorization middleware
	tokenAPI := api.Group("/revoke", tokens.AuthorizeJWT(s.TokenAPI.ValidationKey, s.Issuer, s.Audience))
	{
		tokenAPI.DELETE("/", func(ctx *gin.Context) {
			RevokeHandler(ctx, s)
		})
		tokenAPI.DELETE("/all", func(ctx *gin.Context) {
			RevokeAllHandler(ctx, s)
		})
	}

	return router
}

// SignupHandler creates a new user.
func SignupHandler(ctx *gin.Context, s *Service) {
	userID, err := s.Signup(ctx)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		ctx.AbortWithStatus(http.StatusUnauthorized)
	} else {
		ctx.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "user created successfully!", "userID": userID})
	}
}

// LoginHandler logs a user in and returs a fresh set of tokens.
func LoginHandler(ctx *gin.Context, s *Service) {
	accessToken, refreshToken, role, err := s.Login(ctx)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		ctx.AbortWithStatus(http.StatusUnauthorized)
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
			"role": role,
		})
	}
}

// RefreshHandler refreshes a user's access token.
func RefreshHandler(ctx *gin.Context, s *Service) {
	accessToken, err := s.Refresh(ctx)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		// Different errors might allow to differentiate between the user does not exist or the provided credentials are wrong.
		// Do not leak this detail to protect information who created an account on the platform!
		ctx.AbortWithStatus(http.StatusUnauthorized)
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"accessToken": accessToken,
		})
	}
}

// RevokeHandler revokes a user's tokens for a specific instance or falls back to the authorization tokens instance.
func RevokeHandler(ctx *gin.Context, s *Service) {
	err := s.Revoke(ctx)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		if err == ErrTokenNotFound {
			ctx.Status(http.StatusNotModified)
			return
		}
		ctx.AbortWithStatus(http.StatusUnauthorized)
	} else {
		ctx.Status(http.StatusOK)
	}
}

// RevokeAllHandler revokes a user's tokens for a specific instance or falls back to the authorization tokens instance.
func RevokeAllHandler(ctx *gin.Context, s *Service) {
	err := s.RevokeAll(ctx)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		if err == ErrTokenNotFound {
			ctx.Status(http.StatusNotModified)
			return
		}
		ctx.AbortWithStatus(http.StatusUnauthorized)
	} else {
		ctx.Status(http.StatusOK)
	}
}
