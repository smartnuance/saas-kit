package auth

//go:generate go-bindata -ignore=\\.gitignore -o ./bindata/migrations.go -prefix "migrations/" -pkg gen migrations/...

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	migrateDriver "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	migrateBindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	migrations "github.com/smartnuance/saas-kit/pkg/auth/bindata"
	"github.com/smartnuance/saas-kit/pkg/lib"
)

// Build Variables picked up by govvv
// go get github.com/ahmetb/govvv
var (
	GitCommit string
	Version   string
	BuildDate string
)

type Env struct {
	lib.DatabaseEnv
	release bool
}

type Service struct {
	Env
	DB *gorm.DB
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

	env.release = envs["SAAS_KIT_ENV"] == "prod"

	env.DatabaseEnv = lib.NewDatabaseEnv(envs)
	return
}

func (env Env) Setup() (service Service, err error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", env.Host, env.User, env.Password, env.DBName, env.Port)
	service.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		err = errors.Wrap(err, "failed to connect database at "+dsn)
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

	router := gin.Default()

	// v1 := router.Group("/api/v1")
	// userAPI := v1.Group("/users")
	// {
	// 	userAPI.GET("/:id", user)
	// 	userAPI.PUT("/:id", updateUser)
	// 	userAPI.DELETE("/:id", deleteUser)
	// }
	// tokenAPI := router.Group("/token")
	// {
	// 	tokenAPI.POST("/:id", login)
	// 	tokenAPI.POST("/:id", refresh)
	// 	tokenAPI.POST("/:id", revoke)
	// }

	err = router.Run()
	return
}
