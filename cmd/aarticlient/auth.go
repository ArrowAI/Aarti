package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
	"oras.land/oras-go/v2/registry/remote/auth"

	"aarti/pkg/api"
	hclient "aarti/pkg/http/client"
)

var authGroup = &cobra.Group{ID: "0_auth", Title: "Authentication Commands:"}

var (
	passStdin bool

	loginCmd = &cobra.Command{
		Use:   "login [registry]",
		Short: "Login to an Artifact Registry repository",
		Example: `
Log in with username and password from command line flags:
  aarticlient login -u username -p password localhost:5000

Log in with username and password from stdin:
  aarticlient login -u username --password-stdin localhost:5000

Log in with username and password in an interactive terminal and no TLS check:
  aarticlient login --insecure localhost:5000
`,
		Args:    cobra.ExactArgs(1),
		GroupID: authGroup.ID,
		PreRunE: setup,
		RunE: func(cmd *cobra.Command, args []string) error {
			if user == "" {
				reader := bufio.NewReader(cmd.InOrStdin())
				cmd.Print("Username: ")
				u, err := reader.ReadString('\n')
				if err != nil {
					return err
				}
				user = strings.TrimSpace(u)
				if user == "" {
					return fmt.Errorf("username is required")
				}
			}
			if passStdin {
				reader := bufio.NewReader(cmd.InOrStdin())
				b, err := reader.ReadString('\n')
				if err != nil {
					return err
				}
				pass = strings.TrimSpace(b)
			}
			if pass == "" {
				cmd.Print("Password: ")
				b, err := term.ReadPassword(int(os.Stdin.Fd()))
				if err != nil {
					return err
				}
				fmt.Println()
				pass = strings.TrimSpace(string(b))
				if pass == "" {
					return fmt.Errorf("password is required")
				}
			}
			c, err := api.NewClient(registry, repository, append(opts, hclient.WithBasicAuth(user, pass))...)
			if err != nil {
				return err
			}
			if err := c.Login(cmd.Context()); err != nil {
				return err
			}
			if err := credsStore.Put(cmd.Context(), repoURL(), auth.Credential{Username: user, Password: pass}); err != nil {
				return err
			}
			return nil
		},
	}
	logoutCmd = &cobra.Command{
		Use:     "logout [repository]",
		Short:   "Logout from an Artifact Registry repository",
		GroupID: authGroup.ID,
		Args:    cobra.ExactArgs(1),
		PreRunE: setup,
		RunE: func(cmd *cobra.Command, args []string) error {
			creds, err := credsStore.Get(cmd.Context(), repoURL())
			if err != nil {
				return err
			}
			if creds.Username == "" && creds.Password == "" {
				return nil
			}
			return credsStore.Delete(cmd.Context(), repoURL())
		},
	}
)

func init() {
	loginCmd.Flags().BoolVar(&passStdin, "password-stdin", false, "Take the password from stdin")
	rootCmd.AddCommand(loginCmd, logoutCmd)
	rootCmd.AddGroup(authGroup)
}
