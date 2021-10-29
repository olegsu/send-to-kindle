package main

import (
	"github.com/olegsu/send-to-kindle/cmd/kindle"
)

func main() {
	if err := kindle.Build().Execute(); err != nil {
		panic(err)
	}
}
