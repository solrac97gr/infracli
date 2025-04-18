package main

import (
	"fmt"
	"os"

	"github.com/solrac97gr/infrastructure/cli/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing infracli: %s\n", err)
		os.Exit(1)
	}
}