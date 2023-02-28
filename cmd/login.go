// Package cmd for parsing command line arguments
package cmd

import (
	command "github.com/rawdaGastan/gridify/internal/cmd"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login with mnemonics to a grid network",
	RunE: func(cmd *cobra.Command, args []string) error {
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}

		command.Login(debug)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
