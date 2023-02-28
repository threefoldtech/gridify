// Package cmd for parsing command line arguments
package cmd

import (
	command "github.com/rawdaGastan/gridify/internal/cmd"
	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy the project in the current directory from threefold grid",
	RunE: func(cmd *cobra.Command, args []string) error {
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}

		command.Destroy(debug)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}
