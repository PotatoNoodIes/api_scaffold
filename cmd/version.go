package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version is injected at build time via:
//
//	go build -ldflags "-X api-scaffold/cmd.version=v0.1.0"
var version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the api-scaffold version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("api-scaffold " + version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
