package webbff

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
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

	// without authorization middleware
	router.Any("/auth/*proxyPath", ReverseProxy(s.authServiceAddress))
	router.Any("/event/*proxyPath", ReverseProxy(s.eventServiceAddress))

	return router
}

func ReverseProxy(address string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		remote, err := url.Parse("http://" + address)
		if err != nil {
			panic(err)
		}

		log.Info().Msg(remote.String())

		proxy := httputil.NewSingleHostReverseProxy(remote)
		//Define the director func
		//This is a good place to log, for example
		proxy.Director = func(req *http.Request) {
			req.Header = ctx.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = ctx.Param("proxyPath")
		}

		gin.WrapH(proxy)(ctx)
	}
}