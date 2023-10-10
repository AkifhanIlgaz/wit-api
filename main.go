package main

import (
	"github.com/AkifhanIlgaz/wit-api/setup"
)

func main() {
	err := setup.Run(":3000")
	if err != nil {
		panic(err)
	}
}
