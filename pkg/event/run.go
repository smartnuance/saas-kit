package event

import (
	"context"
	"embed"
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"github.com/RichardKnop/go-fixtures"
	"github.com/friendsofgo/errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/smartnuance/saas-kit/pkg/auth/tokens"
	"github.com/smartnuance/saas-kit/pkg/lib"
	"github.com/smartnuance/saas-kit/pkg/lib/service"
)

//go:embed migrations/*
var migrationDir embed.FS

const ServiceName = "event"

// Build Variables picked up by govvv
// go get github.com/ahmetb/govvv
var (
	GitCommit string
	Version   string
)

// Env is a hierarchical environment configuration for the authentication service and it's API handlers.
type Env struct {
	service.DBEnv
	tokens.TokenEnv
	service.HTTPEnv
	AllowOrigins []string
	release      bool

	modelInfoPath string
}

// Service offers the APIs of the authentication service.
// This struct holds hierarchically structured state that is shared between requests.
type Service struct {
	Env
	service.DBConn
	DBAPI DBAPI
	service.HTTPServer
	TokenAPI     *tokens.TokenController
	AllowOrigins map[string]struct{}
}

var migrateDownFlag bool
var fakeMigrationVersion int
var clearDBFlag bool

func Main() (s Service, err error) {
	// Common steps for all command options
	var env Env
	env, err = Load()
	if err != nil {
		return
	}
	s, err = env.Setup()
	if err != nil {
		return
	}

	// Parse command options
	migrateCommand := flag.NewFlagSet("migrate", flag.ExitOnError)
	migrateCommand.BoolVar(&migrateDownFlag, "down", false, "migrate DB all down to empty")
	migrateCommand.IntVar(&fakeMigrationVersion, "fake", -1, "fakes DB version to specific version without actually migrating")
	migrateCommand.BoolVar(&clearDBFlag, "clear", false, "clear DB")
	fixtureCommand := flag.NewFlagSet("fixture", flag.ExitOnError)
	flag.Parse()

	// Check if a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) >= 2 {
		// Switch on the subcommand and parse the flags for appropriate FlagSet
		// os.Args[2:] will be all arguments starting after the subcommand at os.Args[1]
		switch os.Args[1] {
		case "migrate":
			err = migrateCommand.Parse(os.Args[2:])
			if err != nil {
				return
			}

			if fakeMigrationVersion != -1 {
				err = s.FakeMigration(fakeMigrationVersion)
			} else if migrateDownFlag {
				err = s.MigrateDown()
			} else if clearDBFlag {
				err = s.ClearDB()
			} else {
				err = s.Migrate()
			}
			return
		case "fixture":
			err = fixtureCommand.Parse(os.Args[2:])
			if err != nil {
				return
			}

			var data []byte
			data, err = ioutil.ReadFile(fixtureCommand.Arg(0))
			if err != nil {
				return
			}
			err = fixtures.Load(data, s.DB, "postgres")
			if err != nil {
				return
			}
		default:
			err = errors.Errorf("invalid command: %s", os.Args[1])
			return
		}
	} else {
		// Just migrate up and run the service
		err = s.Migrate()
		if err != nil {
			return
		}
		err = lib.RunInterruptible(s.Run)
		return
	}

	return
}

func Load() (env Env, err error) {
	envs, err := lib.EnvMux(ServiceName)
	if err != nil {
		return
	}

	env.HTTPEnv.Port = envs[strings.ToUpper(ServiceName)+"_SERVICE_PORT"]
	env.release = lib.Stage(envs["SAAS_KIT_ENV"]) == lib.PROD
	var ok bool
	env.modelInfoPath, ok = envs["MODEL_INFO_PATH"]
	if !ok {
		env.modelInfoPath = "./pkg/event/modelinfo"
	}

	env.DBEnv = service.LoadDBEnv(envs)
	env.TokenEnv = tokens.Load(envs, ServiceName)
	env.AllowOrigins = strings.Split(envs["ALLOW_ORIGINS"], ",")
	return
}

func (env Env) Setup() (s Service, err error) {
	s.Env = env

	lib.SetupLogger(ServiceName, Version, env.release)

	s.DBConn, err = service.SetupDB(env.DBEnv, migrationDir)
	if err != nil {
		return
	}
	s.DBAPI = &dbAPI{DB: s.DB}

	s.TokenAPI, err = tokens.Setup(s.TokenEnv)
	if err != nil {
		return
	}

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
