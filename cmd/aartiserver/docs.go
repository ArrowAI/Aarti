//go:build docs

package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsCmd = &cobra.Command{
	Use:    "docs",
	Short:  "Generate documentation",
	Hidden: true,
	Args:   cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := os.MkdirAll(args[0], 0755); err != nil {
			return err
		}
		cmd.Root().DisableAutoGenTag = true
		if err := doc.GenMarkdownTree(cmd.Root(), args[0]); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	cmd.AddCommand(docsCmd)
}
