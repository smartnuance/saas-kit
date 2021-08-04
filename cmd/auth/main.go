package main

import (
	"github.com/smartnuance/saas-kit/pkg/auth"
)

func main() {
	_, err := auth.Main()
	if err != nil {
		panic(err)
	}
}
