// Package cmd for handling command line arguments
package cmd

import (
	"github.com/rawdaGastan/gridify/internal/cmd"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login with mnemonics to a grid network",
	RunE:  cmd.Login,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
