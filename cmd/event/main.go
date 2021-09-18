package main

import (
	"github.com/smartnuance/saas-kit/pkg/event"
)

func main() {
	_, err := event.Main()
	if err != nil {
		panic(err)
	}
}
