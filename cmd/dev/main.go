package main

//go:generate go install github.com/ahmetb/govvv@latest
//go:generate go install github.com/volatiletech/sqlboiler/v4@latest
//go:generate go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest
//go:generate go install github.com/golang/mock/mockgen@latest

import (
	"sync"

	"github.com/smartnuance/saas-kit/pkg/auth"
	"github.com/smartnuance/saas-kit/pkg/lib"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		env, err := auth.Load()
		if err != nil {
			return
		}
		// here we might need to adjust some env values for running services as go routines
		authService, err := env.Setup()
		if err != nil {
			return
		}

		err = lib.RunInterruptible(authService.Run)
		wg.Done()
		return
	}()

	wg.Wait()
}
