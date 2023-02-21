// Package cmd for handling command line arguments
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/rawdaGastan/gridify/internal/config"
	"github.com/rawdaGastan/gridify/internal/deployer"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy the project in the current directory from threefold grid",
	Run: func(cmd *cobra.Command, args []string) {
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

		err = deployer.Destroy()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}
