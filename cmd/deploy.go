// Package cmd for handling command line arguments
package cmd

import (
	"context"
	"os"
	"os/exec"

	"github.com/rawdaGastan/gridify/internal/config"
	"github.com/rawdaGastan/gridify/internal/deployer"
	"github.com/rs/zerolog"
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
			log.Error().Err(err).Send()
			return err
		}

		config, err := config.LoadConfigData()
		if err != nil {
			log.Error().Err(err).Msg("failed to load configuration try logging again")
			return err
		}

		repoURL, err := exec.Command("git", "config", "--get", "remote.origin.url").Output()
		if err != nil {
			log.Error().Err(err).Msg("failed to get remote repository url")
			return err
		}

		logLevel := zerolog.InfoLevel
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			log.Error().Err(err).Send()
			return err
		}
		if debug {
			logLevel = zerolog.DebugLevel
		}

		logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).Level(logLevel).
			With().
			Timestamp().
			Logger()

		deployer, err := deployer.NewDeployer(config.Mnemonics, config.Network, string(repoURL), logger)
		if err != nil {
			log.Error().Err(err).Send()
			return err
		}

		FQDNs, err := deployer.Deploy(context.Background(), ports)
		if err != nil {
			log.Error().Err(err).Send()
			return err
		}
		for port, FQDN := range FQDNs {
			logger.Info().Msgf("%d: %s\n", port, FQDN)
		}
		return nil
	},
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
