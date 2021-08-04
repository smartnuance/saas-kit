package lib

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type DatabaseEnv struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func NewDatabaseEnv(envs map[string]string) DatabaseEnv {
	return DatabaseEnv{
		Host:     envs["DB_HOST"],
		Port:     envs["DB_PORT"],
		User:     envs["DB_USER"],
		Password: envs["DB_PASSWORD"],
		DBName:   envs["DB_NAME"],
	}
}

func EnvMux() (envs map[string]string, err error) {
	env := os.Getenv("SAAS_KIT_ENV")
	if "" == env {
		env = "dev"
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
