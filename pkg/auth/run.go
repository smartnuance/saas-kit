package auth

//go:generate go-bindata -ignore=\\.gitignore -o ./bindata/migrations.go -prefix "migrations/" -pkg gen migrations/...

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	migrateDriver "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	migrateBindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/pkg/errors"
	"gorm.io/gorm"

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
	DB       *gorm.DB
	TokenAPI *tokens.TokenController
}

func Main() (authService Service, err error) {
	var env Env
	env, err = Load()
	if err != nil {
		return
	}
	authService, err = env.Setup()
	if err != nil {
		return
	}
	err = authService.Migrate()
	if err != nil {
		return
	}
	err = authService.Run()
	if err != nil {
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
	env.TokenEnv = tokens.Load(envs)
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

// Migrate migrates with the DB instance of the service up to the newest version
func (s *Service) Migrate() error {
	db, err := s.DB.DB()
	if err != nil {
		return err
	}
	driver, err := migrateDriver.WithInstance(db, &migrateDriver.Config{})
	if err != nil {
		return err
	}

	migrationsRessource := migrateBindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		})
	d, err := migrateBindata.WithInstance(migrationsRessource)
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithInstance("go-bindata", d, s.Env.DBName, driver)
	if err != nil {
		return errors.Wrap(err, "failed to migrate database "+s.Env.DBName)
	}
	err = migrator.Up()
	// ignore error happing on no change to database necessary
	if err == migrate.ErrNoChange {
		return nil
	}
	return err
}

// AutoMigrate alters tables and constraints with the DB instance of the service using gin's auto migration
//
// Only use this during development to write sql files in migrations/.
func (s *Service) AutoMigrate() error {
	return s.DB.AutoMigrate(&Profile{})
}

func (s *Service) Run() (err error) {
	err = router(s).Run()
	return
}
