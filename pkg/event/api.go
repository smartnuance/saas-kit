package event

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

	// with authorization middleware
	api := router.Group("/", tokens.AuthorizeJWT(s.TokenAPI.ValidationKey, s.Issuer, s.Audience))

	// without authorization middleware
	api.PUT("/workshop", func(ctx *gin.Context) {
		CreateWorkshopHandler(ctx, s)
	})

	return router
}

// CreateWorkshopHandler creates a new workshop.
var CreateWorkshopHandler = func(ctx *gin.Context, s *Service) {
	workshop, err := s.CreateWorkshop(ctx)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		ctx.AbortWithStatus(http.StatusUnauthorized)
	} else {
		ctx.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "workshop created successfully!", "workshopID": workshop.ID})
	}
}
