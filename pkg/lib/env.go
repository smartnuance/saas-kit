package lib

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"

	"fmt"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	DBName   string
}

func LoadDatabaseEnv(envs map[string]string) DatabaseEnv {
	return DatabaseEnv{
		Host:     envs["DB_HOST"],
		Port:     envs["DB_PORT"],
		User:     envs["DB_USER"],
		Password: envs["DB_PASSWORD"],
		DBName:   envs["DB_NAME"],
	}
}

func SetupDatabase(env DatabaseEnv) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", env.Host, env.User, env.Password, env.DBName, env.Port)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		err = errors.Wrap(err, "failed to connect database at "+dsn)
		return
	}
	return
}

func EnvMux() (envs map[string]string, err error) {
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

	envs["SAAS_KIT_ENV"] = env

	return
}
