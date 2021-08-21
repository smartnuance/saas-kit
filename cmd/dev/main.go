package main

//go:generate go install github.com/ahmetb/govvv@latest
//go:generate go install github.com/volatiletech/sqlboiler/v4@latest
//go:generate go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest

import (
	"log"
	"sync"

	"github.com/smartnuance/saas-kit/pkg/auth"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		_, err := auth.Main()
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()
}
