package auth

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/smartnuance/saas-kit/pkg/lib/tokens"
)

func router(s *Service) *gin.Engine {
	var router = gin.Default()

	config := cors.DefaultConfig()
	config.AddAllowHeaders("PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS")
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
		tokenAPI.DELETE("/:user_id", func(ctx *gin.Context) {
			RevokeHandler(ctx, s)
		})
	}
	// userAPI := api.Group("/user", tokens.AuthorizeJWT(s.TokenAPI.ValidationKey))
	// {
	// 	userAPI.GET("/:id", user)
	// 	userAPI.DELETE("/:id", deleteUser)
	// }

	return router
}

// SignupHandler creates a new user.
var SignupHandler = func(ctx *gin.Context, s *Service) {
	userID, err := s.Signup(ctx)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		ctx.AbortWithStatus(http.StatusUnauthorized)
	} else {
		ctx.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "user created successfully!", "userID": userID})
	}
}

// LoginHandler logs a user in and returs a fresh set of tokens.
var LoginHandler = func(ctx *gin.Context, s *Service) {
	accessToken, refreshToken, err := s.Login(ctx)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		ctx.AbortWithStatus(http.StatusUnauthorized)
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		})
	}
}

// RefreshHandler refreshes a user's access token.
var RefreshHandler = func(ctx *gin.Context, s *Service) {
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

// RevokeHandler refreshes a user's access token.
var RevokeHandler = func(ctx *gin.Context, s *Service) {
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
