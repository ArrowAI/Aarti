package main

import (
	"fmt"

	"github.com/spf13/cobra"

	artifact_registry "aarti"
)

var (
	cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Print the version information and exit",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\n", artifact_registry.Version)
			fmt.Printf("Commit: %s\n", artifact_registry.Commit)
			fmt.Printf("Date: %s\n", artifact_registry.Date)
			fmt.Printf("Repo: https://github.com/%s\n", artifact_registry.Repo)
		},
	}
)

func init() {
	rootCmd.AddCommand(cmdVersion)
}
