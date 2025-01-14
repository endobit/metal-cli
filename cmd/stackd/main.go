// Package main implements the metald service.
package main

import (
	"os"

	"github.com/endobit/stack/internal/stackd"
)

var version string

func main() {
	cmd := stackd.NewRootCmd()
	cmd.Version = version

	if err := cmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
