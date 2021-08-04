package main

//go:generate go get github.com/ahmetb/govvv

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
