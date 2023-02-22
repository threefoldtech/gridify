// Package cmd for handling command line arguments
package cmd

import (
	"os"

	"github.com/rawdaGastan/gridify/internal/cmd"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the project in the current directory on threefold grid",
	RunE:  cmd.Deploy,
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().UintSliceP("ports", "p", []uint{}, "ports to forward the FQDNs to")
	err := deployCmd.MarkFlagRequired("ports")
	if err != nil {
		log.Error().Err(err).Send()
		os.Exit(1)
	}
}
