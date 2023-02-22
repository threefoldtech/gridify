// Package cmd for handling command line arguments
package cmd

import (
	"github.com/rawdaGastan/gridify/internal/cmd"
	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy the project in the current directory from threefold grid",
	RunE:  cmd.Destroy,
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}
