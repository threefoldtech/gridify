// Package cmd for parsing command line arguments
package cmd

import (
	command "github.com/rawdaGastan/gridify/internal/cmd"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the project in the current directory on threefold grid",
	RunE: func(cmd *cobra.Command, args []string) error {
		ports, err := cmd.Flags().GetUintSlice("ports")
		if err != nil {
			return err
		}

		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}
		vmSpec, err := cmd.Flags().GetString("spec")
		if err != nil {
			return err
		}

		err = command.Deploy(vmSpec, ports, debug)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().UintSliceP("ports", "p", []uint{}, "ports to forward the FQDNs to")
	err := deployCmd.MarkFlagRequired("ports")
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	deployCmd.Flags().StringP("spec", "s", "eco", "vm spec can be (eco, standard, performance)")

}
