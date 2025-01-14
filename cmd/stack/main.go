// Package main implements the cast CLI.
package main

import (
	"os"

	"github.com/endobit/stack/internal/stack"
)

var version string

func main() {
	cmd := stack.NewRootCmd()
	cmd.Version = version

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
