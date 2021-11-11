package webbff

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/smartnuance/saas-kit/pkg/lib"
	"github.com/smartnuance/saas-kit/pkg/lib/service"
)

const ServiceName = "webbff"

// Build Variables picked up by govvv
// go get github.com/ahmetb/govvv
var (
	GitCommit string
	Version   string
)

// Env is a hierarchical environment configuration for the authentication service and it's API handlers.
type Env struct {
	service.HTTPEnv
	authServiceAddress  string
	eventServiceAddress string
	AllowOrigins        []string
	release             bool
}

// Service offers the APIs of the authentication service.
// This struct holds hierarchically structured state that is shared between requests.
type Service struct {
	Env
	service.HTTPServer
	AllowOrigins map[string]struct{}
}

func Main() (webbffService Service, err error) {
	// Common steps for all command options
	var env Env
	env, err = Load()
	if err != nil {
		return
	}
	webbffService, err = env.Setup()
	if err != nil {
		return
	}
	err = lib.RunInterruptible(webbffService.Run)
	return
}

func Load() (env Env, err error) {
	envs, err := lib.EnvMux(ServiceName)
	if err != nil {
		return
	}

	env.Port = envs[strings.ToUpper(ServiceName)+"_SERVICE_PORT"]
	env.authServiceAddress = envs["AUTH_SERVICE_HOST"] + ":" + envs["AUTH_SERVICE_PORT"]
	env.eventServiceAddress = envs["EVENT_SERVICE_HOST"] + ":" + envs["EVENT_SERVICE_PORT"]
	env.release = lib.Stage(envs["SAAS_KIT_ENV"]) == lib.PROD

	env.AllowOrigins = strings.Split(envs["ALLOW_ORIGINS"], ",")
	return
}

func (env Env) Setup() (s Service, err error) {
	s.Env = env

	lib.SetupLogger(ServiceName, Version, env.release)

	s.HTTPServer = service.SetupHTTP(env.HTTPEnv, router(&s))

	s.AllowOrigins = map[string]struct{}{}
	for _, o := range env.AllowOrigins {
		s.AllowOrigins[o] = struct{}{}
	}

	if env.release {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Info().Str("port", s.HTTPServer.Port).Str("gitCommit", GitCommit).Msg("setup")

	return
}

func (s *Service) Run(ctx context.Context) (err error) {
	return s.Serve(ctx)
}
