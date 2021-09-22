package lib

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Stage string

const (
	DEV  Stage = "dev"
	TEST Stage = "test"
	PROD Stage = "prod"
)

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
