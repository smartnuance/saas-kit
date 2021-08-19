package auth

//go:generate go-bindata -ignore=\\.gitignore -o ./bindata/migrations.go -prefix "migrations/" -pkg gen migrations/...
//go:generate go-bindata -ignore=\\.gitignore -o ./bindata/migrations.go -prefix "migrations/" -pkg gen migrations/...

import (
	"database/sql"
	"flag"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	migrateDriver "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	migrateBindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/pkg/errors"

	"github.com/rs/zerolog/log"
	migrations "github.com/smartnuance/saas-kit/pkg/auth/bindata"
	"github.com/smartnuance/saas-kit/pkg/auth/tokens"
	"github.com/smartnuance/saas-kit/pkg/lib"
)

const ServiceName = "auth"

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
	release bool
}

// Service offers the APIs of the authentication service.
// This struct holds hierarchically structured state that is shared between requests.
type Service struct {
	Env
	DB       *sql.DB
	TokenAPI *tokens.TokenController
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
	flag.Parse()

	// Check if a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) >= 2 {
		// Switch on the subcommand and parse the flags for appropriate FlagSet
		// os.Args[2:] will be all arguments starting after the subcommand at os.Args[1]
		switch os.Args[1] {
		case "migrate":
			migrateCommand.Parse(os.Args[2:])

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
		default:
			os.Exit(1)
		}
	} else {
		// Just migrate up and run the service
		err = authService.Migrate()
		if err != nil {
			return
		}
		err = authService.Run()
		return
	}

	return
}

func Load() (env Env, err error) {
	envs, err := lib.EnvMux()
	if err != nil {
		return
	}

	env.release = lib.Stage(envs["SAAS_KIT_ENV"]) == lib.PROD

	env.DatabaseEnv = lib.LoadDatabaseEnv(envs)
	env.TokenEnv = tokens.Load(envs, ServiceName)
	return
}

func (env Env) Setup() (s Service, err error) {
	s.Env = env

	lib.SetupLogger(ServiceName, Version, env.release)

	log.Debug().Str("GitCommit", GitCommit).Interface("ServiceStruct", s).Msg("Setup service")

	s.DB, err = lib.SetupDatabase(env.DatabaseEnv)
	if err != nil {
		return
	}

	s.TokenAPI, err = tokens.Setup(s.TokenEnv)
	if err != nil {
		return
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

	migrationsRessource := migrateBindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		})
	d, err := migrateBindata.WithInstance(migrationsRessource)
	if err != nil {
		return nil, err
	}

	migrator, err := migrate.NewWithInstance("go-bindata", d, s.Env.DBName, driver)
	if err != nil {
		return nil, errors.Wrap(err, "failed to migrate database "+s.Env.DBName)
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

func (s *Service) Run() (err error) {
	err = router(s).Run()
	return
}
