package main

import (
	"embed"
	"os"
	"sync"

	"github.com/friendsofgo/errors"
	"github.com/rs/zerolog/log"
	"github.com/smartnuance/saas-kit/pkg/auth"
	"github.com/smartnuance/saas-kit/pkg/event"
	"github.com/smartnuance/saas-kit/pkg/lib"
	"github.com/smartnuance/saas-kit/pkg/lib/service"
	"github.com/smartnuance/saas-kit/pkg/webbff"
)

//go:embed migrations/*
var migrationDir embed.FS

func Main() (err error) {
	if len(os.Args) >= 2 {
		// Switch on the subcommand and parse the flags for appropriate FlagSet
		// os.Args[2:] will be all arguments starting after the subcommand at os.Args[1]
		switch os.Args[1] {
		case "init":
			err = execSQL("schemas.up.sql")
			return
		case "deinit":
			err = execSQL("schemas.down.sql")
			return
		case "webbff":
			log.Info().Msg("Redirect to webbff service...")
			os.Args = append(os.Args[0:1], os.Args[2:]...)
			_, err = webbff.Main()
			return
		case "auth":
			log.Info().Msg("Redirect to auth service...")
			os.Args = append(os.Args[0:1], os.Args[2:]...)
			_, err = auth.Main()
			return
		case "event":
			log.Info().Msg("Redirect to event service...")
			os.Args = append(os.Args[0:1], os.Args[2:]...)
			_, err = event.Main()
			return
		default:
			err = errors.Errorf("invalid command: %s", os.Args[1])
			return
		}
	} else {
		log.Info().Msg("Starting all services in separate goroutines...")
		err = runAll()
		return
	}
}

func runAll() error {
	var wg sync.WaitGroup

	errs := make(chan error, 1)
	waitCh := make(chan struct{}, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()

		env, err := auth.Load()
		if err != nil {
			errs <- err
			return
		}

		authService, err := env.Setup()
		if err != nil {
			errs <- err
			return
		}

		err = lib.RunInterruptible(authService.Run)
		if err != nil {
			errs <- err
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		env, err := event.Load()
		if err != nil {
			errs <- err
			return
		}

		eventService, err := env.Setup()
		if err != nil {
			errs <- err
			return
		}

		err = lib.RunInterruptible(eventService.Run)
		if err != nil {
			errs <- err
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		env, err := webbff.Load()
		if err != nil {
			errs <- err
			return
		}

		webbffService, err := env.Setup()
		if err != nil {
			errs <- err
			return
		}

		err = lib.RunInterruptible(webbffService.Run)
		if err != nil {
			errs <- err
			return
		}
	}()

	go func() {
		wg.Wait()
		close(waitCh)
	}()

	select {
	case err := <-errs:
		return err
	case <-waitCh:
	}

	return nil
}

func execSQL(script string) error {
	envs, err := lib.EnvMux("")
	if err != nil {
		return err
	}

	databaseEnv := service.LoadDBEnv(envs)

	db, err := service.SetupDB(databaseEnv, embed.FS{})
	if err != nil {
		return err
	}

	c, err := migrationDir.ReadFile("migrations/" + script)
	if err != nil {
		return err
	}
	sql := string(c)
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}

	db.Close()
	return nil
}

func main() {
	lib.SetupLogger("dev", "", false)

	err := Main()
	if err != nil {
		panic(err)
	}
}
