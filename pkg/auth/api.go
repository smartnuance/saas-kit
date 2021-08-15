package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func router(s *Service) *gin.Engine {
	var router = gin.Default()

	v1 := router.Group("/api/v1")
	// userAPI := v1.Group("/users")
	// {
	// 	userAPI.GET("/:id", user)
	// 	userAPI.PUT("/:id", updateUser)
	// 	userAPI.DELETE("/:id", deleteUser)
	// }
	tokenAPI := v1.Group("/token")
	{
		tokenAPI.POST("/login", func(ctx *gin.Context) {
			LoginHandler(ctx, s)
		})
		// tokenAPI.POST("/:id", refresh)
		// tokenAPI.POST("/:id", revoke)
	}

	return router
}

// LoginHandler logs a user in and returs a fresh set of tokens.
var LoginHandler = func(ctx *gin.Context, s *Service) {
	token, err := s.Login(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, nil)
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}
