package main

import (
	"os"

	"github.com/ujjwaltalwar/betteruk-bot/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
