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
	config.AddAllowMethods("PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS")
	config.AddAllowHeaders("Authorization")
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
	s.AddInfoHandlers(api.Group("/info"))

	// without authorization middleware
	api.PUT("/workshop", CreateWorkshopHandler(s))

	return router
}

// AddInfoHandlers adds new handlers to retrieve model structure info.
func (s *Service) AddInfoHandlers(routerGroup *gin.RouterGroup) {
	dir := http.Dir(s.modelInfoPath)
	routerGroup.GET("/workshop", func(ctx *gin.Context) {
		ctx.FileFromFS("/workshop.json", dir)
	})
}

// CreateWorkshopHandler creates a new workshop.
func CreateWorkshopHandler(s *Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		workshop, err := s.CreateWorkshop(ctx)
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			ctx.AbortWithStatus(http.StatusUnauthorized)
		} else {
			ctx.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "workshop created successfully!", "workshopID": workshop.ID})
		}
	}
}
