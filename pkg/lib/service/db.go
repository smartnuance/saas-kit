package service

import (
	"database/sql"
	"embed"
	"fmt"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	migrateDriver "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/pkg/errors"
)

type DBEnv struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Schema   string
}

type DBConn struct {
	DBEnv
	*sql.DB
	migrationDir embed.FS
}

func LoadDBEnv(envs map[string]string) DBEnv {
	return DBEnv{
		Host:     envs["DB_HOST"],
		Port:     envs["DB_PORT"],
		User:     envs["DB_USER"],
		Password: envs["DB_PASSWORD"],
		DBName:   envs["DB_NAME"],
		Schema:   envs["DB_SCHEMA"],
	}
}

func SetupDB(env DBEnv, migrationDir embed.FS) (conn DBConn, err error) {
	conn.migrationDir = migrationDir

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", env.Host, env.User, env.Password, env.DBName, env.Port)
	if env.Schema != "" {
		dsn += fmt.Sprintf(" search_path=%s", env.Schema)
	}
	conn.DB, err = sql.Open("postgres", dsn)
	if err != nil {
		err = errors.Wrap(err, "failed to connect database at "+dsn)
		return
	}
	return
}

func (s *DBConn) migrator() (*migrate.Migrate, error) {
	driver, err := migrateDriver.WithInstance(s.DB, &migrateDriver.Config{})
	if err != nil {
		return nil, err
	}

	src, err := httpfs.New(http.FS(s.migrationDir), "migrations")
	if err != nil {
		return nil, err
	}
	migrator, err := migrate.NewWithInstance("httpfs", src, s.DBName, driver)
	if err != nil {
		return nil, errors.Wrap(err, "failed to migrate database "+s.DBName)
	}
	return migrator, nil
}

// Migrate migrates the DB up to the newest version
// Uses the DB instance of the service.
func (s *DBConn) Migrate() error {
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
func (s *DBConn) ClearDB() error {
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
func (s *DBConn) MigrateDown() error {
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
func (s *DBConn) FakeMigration(version int) error {
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
