package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"aarti/pkg/packages"
	"aarti/pkg/packages/apk"
	"aarti/pkg/packages/deb"
	"aarti/pkg/packages/helm"
	"aarti/pkg/packages/rpm"
)

func newPkgDeleteCmd(typ string) *cobra.Command {
	return &cobra.Command{
		Use:     fmt.Sprintf("delete [repository] [path]"),
		Short:   fmt.Sprintf("Delete %s package from the repository", typ),
		Aliases: []string{"rm", "remove", "del"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				c   packages.Deleter
				err error
			)
			switch typ {
			case apk.Name:
				c, err = apk.NewClient(registry, repository, "", "", opts...)
			case deb.Name:
				c, err = deb.NewClient(registry, repository, "", "", opts...)
			case rpm.Name:
				c, err = rpm.NewClient(registry, repository, opts...)
			case helm.Name:
				c, err = helm.NewClient(registry, repository, opts...)
			default:
				return fmt.Errorf("unsupported package type: %s", typ)
			}
			if err != nil {
				return err
			}
			return c.Delete(cmd.Context(), args[1])
		},
	}
}
