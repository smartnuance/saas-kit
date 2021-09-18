package event

import (
	"context"
	"database/sql"
	"embed"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/RichardKnop/go-fixtures"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	migrateDriver "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartnuance/saas-kit/pkg/auth/tokens"
	"github.com/smartnuance/saas-kit/pkg/lib"
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
	lib.DatabaseEnv
	tokens.TokenEnv
	port         string
	AllowOrigins []string
	release      bool
}

// Service offers the APIs of the authentication service.
// This struct holds hierarchically structured state that is shared between requests.
type Service struct {
	Env
	DB           *sql.DB
	DBAPI        DBAPI
	TokenAPI     *tokens.TokenController
	AllowOrigins map[string]struct{}
}

var migrateDownFlag bool
var fakeMigrationVersion int
var clearDBFlag bool

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

	env.port = envs["HTTP_PORT"]
	env.release = lib.Stage(envs["SAAS_KIT_ENV"]) == lib.PROD

	env.DatabaseEnv = lib.LoadDatabaseEnv(envs)
	env.TokenEnv = tokens.Load(envs, ServiceName)
	env.AllowOrigins = strings.Split(envs["ALLOW_ORIGINS"], ",")
	return
}

func (env Env) Setup() (s Service, err error) {
	s.Env = env

	lib.SetupLogger(ServiceName, Version, env.release)

	log.Info().Str("GitCommit", GitCommit).Msg("Setup service")

	s.DB, err = lib.SetupDatabase(env.DatabaseEnv)
	if err != nil {
		return
	}
	s.DBAPI = &dbAPI{DB: s.DB}

	s.TokenAPI, err = tokens.Setup(s.TokenEnv)
	if err != nil {
		return
	}

	s.AllowOrigins = map[string]struct{}{}
	for _, o := range env.AllowOrigins {
		s.AllowOrigins[o] = struct{}{}
	}

	if env.release {
		gin.SetMode(gin.ReleaseMode)
	}

	return
}

func (s *Service) migrator() (*migrate.Migrate, error) {
	driver, err := migrateDriver.WithInstance(s.DB, &migrateDriver.Config{})
	if err != nil {
		return nil, err
	}

	src, err := httpfs.New(http.FS(migrationDir), "migrations")
	if err != nil {
		return nil, err
	}
	migrator, err := migrate.NewWithInstance("httpfs", src, s.Env.DB, driver)
	if err != nil {
		return nil, errors.Wrap(err, "failed to migrate database "+s.Env.DB)
	}
	return migrator, nil
}

// Migrate migrates the DB up to the newest version
// Uses the DB instance of the service.
func (s *Service) Migrate() error {
	migrator, err := s.migrator()
	if err != nil {
		return err
	}
	err = migrator.Up()
	// ignore error happing on no change to database necessary
	if err == migrate.ErrNoChange {
		return nil
	}
	return err
}

// ClearDB migrates the DB down to an empty database.
// Uses the DB instance of the service.
func (s *Service) ClearDB() error {
	migrator, err := s.migrator()
	if err != nil {
		return err
	}
	err = migrator.Drop()
	// ignore error happing on no change to database necessary
	if err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// MigrateDown migrates the DB down to an empty database.
// Uses the DB instance of the service.
func (s *Service) MigrateDown() error {
	migrator, err := s.migrator()
	if err != nil {
		return err
	}
	err = migrator.Down()
	// ignore error happing on no change to database necessary
	if err != migrate.ErrNoChange {
		return err
	}
	return nil
}

// FakeMigration fakes a specific version without migrating.
// Uses the DB instance of the service.
func (s *Service) FakeMigration(version int) error {
	migrator, err := s.migrator()
	if err != nil {
		return err
	}
	err = migrator.Force(version)
	// ignore error happing on no change to database necessary
	if err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func (s *Service) Run(ctx context.Context) (err error) {
	srv := &http.Server{
		Addr:    ":" + s.port,
		Handler: router(s),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err)
		}
	}()

	<-ctx.Done()
	log.Info().Msg("gracefully shutdown service...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Stack().Err(err).Msg("error during shutdown")
	}
	log.Info().Msg("...shutdown done")

	return
}
