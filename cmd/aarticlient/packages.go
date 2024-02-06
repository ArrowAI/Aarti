package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"aarti/pkg/packages/apk"
	"aarti/pkg/packages/deb"
	"aarti/pkg/packages/helm"
	"aarti/pkg/packages/rpm"
)

var PkgGroup = &cobra.Group{ID: "2_packages", Title: "Package Commands:"}

func newPkgCmd(typ string) *cobra.Command {
	pkgCmd := &cobra.Command{
		Use:               typ,
		Short:             fmt.Sprintf("Manage %s packages", typ),
		GroupID:           PkgGroup.ID,
		PersistentPreRunE: setup,
	}
	pkgCmd.AddCommand(
		newPkgListCmd(typ),
		newPkgPushCmd(typ),
		newPkgPullCmd(typ),
		newPkgDeleteCmd(typ),
		newPkgSetupCmd(typ),
	)
	return pkgCmd
}

func init() {
	rootCmd.AddGroup(PkgGroup)
	for _, v := range []string{apk.Name, deb.Name, rpm.Name, helm.Name} {
		rootCmd.AddCommand(newPkgCmd(v))
	}
}
