package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/kasterism/astermule/cmd/app"
)

func main() {
	newRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	newRand.Seed(time.Now().UnixNano())

	command := app.NewRootCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
