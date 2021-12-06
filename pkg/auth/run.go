package auth

import (
	"context"
	"embed"
	"flag"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/RichardKnop/go-fixtures"
	"github.com/friendsofgo/errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	m "github.com/smartnuance/saas-kit/pkg/auth/dbmodels"
	"github.com/smartnuance/saas-kit/pkg/auth/tokens"
	"github.com/smartnuance/saas-kit/pkg/lib"
	"github.com/smartnuance/saas-kit/pkg/lib/service"
)

//go:embed migrations/*
var migrationDir embed.FS

const ServiceName = "auth"

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
var userName string
var userEmail string
var userPassword string
var userInstanceURL string

func Main() (authService Service, err error) {
	// Common steps for all command options
	var env Env
	env, err = Load()
	if err != nil {
		return
	}
	authService, err = env.Setup()
	if err != nil {
		return
	}

	// Parse command options
	migrateCommand := flag.NewFlagSet("migrate", flag.ExitOnError)
	migrateCommand.BoolVar(&migrateDownFlag, "down", false, "migrate DB all down to empty")
	migrateCommand.IntVar(&fakeMigrationVersion, "fake", -1, "fakes DB version to specific version without actually migrating")
	migrateCommand.BoolVar(&clearDBFlag, "clear", false, "clear DB")
	fixtureCommand := flag.NewFlagSet("fixture", flag.ExitOnError)
	userCommand := flag.NewFlagSet("adduser", flag.ExitOnError)
	userCommand.StringVar(&userName, "name", "", "name of user to add")
	userCommand.StringVar(&userEmail, "email", "", "email of user to add")
	userCommand.StringVar(&userPassword, "password", "", "password of user to add")
	userCommand.StringVar(&userInstanceURL, "instance", "smartnuance.com", "instance URL for which to add user's default profile")
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
				err = authService.FakeMigration(fakeMigrationVersion)
			} else if migrateDownFlag {
				err = authService.MigrateDown()
			} else if clearDBFlag {
				err = authService.ClearDB()
			} else {
				err = authService.Migrate()
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
			err = fixtures.Load(data, authService.DB, "postgres")
			if err != nil {
				return
			}
		case "adduser":
			err = userCommand.Parse(os.Args[2:])
			if err != nil {
				return
			}

			r := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(r)

			var instance *m.Instance
			instance, err = authService.DBAPI.GetInstance(ctx, userInstanceURL)
			if err != nil {
				return
			}

			_, err = authService.signup(ctx, instance.ID, SignupBody{Name: userName, Email: userEmail, Password: userPassword}, "super admin")
			if err != nil {
				return
			}
		default:
			err = errors.Errorf("invalid command: %s", os.Args[1])
			return
		}
	} else {
		// Just migrate up and run the service
		err = authService.Migrate()
		if err != nil {
			return
		}
		err = lib.RunInterruptible(authService.Run)
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

	log.Info().Str("port", s.HTTPServer.Port).Str("gitCommit", GitCommit).Str("schema", env.DBEnv.Schema).Msg("setup")

	return
}

func (s *Service) Run(ctx context.Context) (err error) {
	return s.Serve(ctx)
}
