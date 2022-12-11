package event

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/smartnuance/saas-kit/pkg/lib/roles"
	"github.com/smartnuance/saas-kit/pkg/lib/tokens"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
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

	// with authorization middleware
	api := router.Group("/", tokens.AuthorizeJWT(s.TokenAPI.ValidationKey, s.Issuer, s.Audience))
	api.PUT("/workshop", s.CreateWorkshopHandler())
	api.GET("/workshop/list", s.ListWorkshopHandler())
	api.DELETE("/workshop/:id", s.DeleteWorkshopHandler())

	return router
}

// CreateWorkshopHandler creates a new workshop.
func (s *Service) CreateWorkshopHandler() gin.HandlerFunc {
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

// ListWorkshopHandler lists workshops.
func (s *Service) ListWorkshopHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		workshops, err := s.ListWorkshops(ctx)
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			ctx.AbortWithStatus(http.StatusUnauthorized)
		} else {
			respondProto(ctx, &workshops)
		}
	}
}

// DeleteWorkshopHandler deletes a workshop.
func (s *Service) DeleteWorkshopHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := s.DeleteWorkshop(ctx)
		if err != nil {
			log.Error().Stack().Err(err).Msg("")
			ctx.AbortWithStatus(http.StatusUnauthorized)
		} else {
			ctx.Status(http.StatusOK)
		}
	}
}

func respondProto(ctx *gin.Context, m proto.Message) {
	jsonData, err := protojson.Marshal(m)
	if err != nil {
		log.Error().Stack().Err(err).Msg("")
		ctx.AbortWithStatus(http.StatusInternalServerError)
	} else {
		ctx.Data(http.StatusOK, "application/json", jsonData)
	}
}
