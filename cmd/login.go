// Package cmd for handling command line arguments
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bernix/cosmos-key-sign/cosmos/bip39"
	"github.com/rawdaGastan/gridify/internal/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login with mnemonics on a grid network",
	RunE: func(cmd *cobra.Command, args []string) error {
		scanner := bufio.NewReader(os.Stdin)

		fmt.Print("Please enter your mnemonics: ")
		mnemonics, err := scanner.ReadString('\n')
		if err != nil {
			log.Error().Err(err).Msg("failed to read mnemonics")
			return err
		}
		mnemonics = strings.TrimSpace(mnemonics)
		if !bip39.IsMnemonicValid(mnemonics) {
			log.Error().Msg("failed to validate mnemonics")
			return errors.New("login failed")
		}

		fmt.Print("Please enter grid network: ")
		network, err := scanner.ReadString('\n')
		if err != nil {
			log.Error().Err(err).Msg("failed to read network")
			return err
		}
		network = strings.TrimSpace(network)

		if network != "dev" && network != "qa" && network != "test" && network != "main" {
			log.Error().Msg("invalid network, must be one of: dev, test, qa and main")
			return errors.New("login failed")
		}

		err = config.SaveConfigData(mnemonics, network)
		if err != nil {
			log.Error().Err(err).Send()
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
