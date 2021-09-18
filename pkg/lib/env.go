package lib

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"

	"fmt"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Stage string

const (
	DEV  Stage = "dev"
	TEST Stage = "test"
	PROD Stage = "prod"
)

type DatabaseEnv struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
	Schema   string
}

func LoadDatabaseEnv(envs map[string]string) DatabaseEnv {
	return DatabaseEnv{
		Host:     envs["DB_HOST"],
		Port:     envs["DB_PORT"],
		User:     envs["DB_USER"],
		Password: envs["DB_PASSWORD"],
		DB:       envs["DB_NAME"],
		Schema:   envs["DB_SCHEMA"],
	}
}

func SetupDatabase(env DatabaseEnv) (db *sql.DB, err error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", env.Host, env.User, env.Password, env.DB, env.Port)
	if env.Schema != "" {
		dsn += fmt.Sprintf(" search_path=%s", env.Schema)
	}
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		err = errors.Wrap(err, "failed to connect database at "+dsn)
		return
	}
	boil.SetDB(db)
	return
}

func EnvMux(serviceName string) (envs map[string]string, err error) {
	env := os.Getenv("SAAS_KIT_ENV")
	if env == "" {
		env = string(DEV)
	}

	envs, err = godotenv.Read()
	if err != nil {
		err = errors.Wrap(err, "error loading base .env file")
		return
	}

	p := ".env." + env
	envOverrides, err := godotenv.Read(p)
	if err != nil {
		err = errors.Wrap(err, "error loading env file from "+p)
		return
	}

	for k, v := range envOverrides {
		// environment-specific values take precedence
		envs[k] = v
	}

	if serviceName != "" {
		p = ".env." + serviceName
		envOverrides, err = godotenv.Read(p)
		if err != nil {
			err = errors.Wrap(err, "error loading env file from "+p)
			return
		}
	}

	for k, v := range envOverrides {
		// service-specific values take precedence
		envs[k] = v
	}

	envs["SAAS_KIT_ENV"] = env

	return
}
