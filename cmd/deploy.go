// Package cmd for handling command line arguments
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/rawdaGastan/gridify/internal/config"
	"github.com/rawdaGastan/gridify/internal/deployer"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the project in the current directory on threefold grid",
	Run: func(cmd *cobra.Command, args []string) {
		ports, err := cmd.Flags().GetUintSlice("ports")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		config, err := config.LoadConfigData()
		if err != nil {
			fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed to load configuration try logging again"))
			os.Exit(1)
		}

		repoURL, err := exec.Command("git", "config", "--get", "remote.origin.url").Output()
		if err != nil {
			fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed to get remote repository url"))
			os.Exit(1)
		}

		logLevel := zerolog.InfoLevel
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
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
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		FQDNs, err := deployer.Deploy(context.Background(), ports)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for port, FQDN := range FQDNs {
			logger.Info().Msgf("%d: %s\n", port, FQDN)
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().UintSliceP("ports", "p", []uint{}, "ports to forward the FQDNs to")
	deployCmd.MarkFlagRequired("ports")
}
