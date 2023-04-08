package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/kasterism/astermule/cmd/app"
)

// ActionMap is used globally to store the list of function entry parameters
var ActionMap map[string]ParamMap

// the list of function entry parameters
type ParamMap map[string]string

// ExitName stores the names of all exit microservices
var ExitName []string

func main() {
	newRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	newRand.Seed(time.Now().UnixNano())

	command := app.NewRootCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
