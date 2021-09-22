package main

import (
	"github.com/smartnuance/saas-kit/pkg/webbff"
)

func main() {
	_, err := webbff.Main()
	if err != nil {
		panic(err)
	}
}
