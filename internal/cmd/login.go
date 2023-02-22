// Package cmd for handling commands
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	bip39 "github.com/cosmos/go-bip39"
	"github.com/rawdaGastan/gridify/internal/config"
	"github.com/rs/zerolog/log"
)

// Login handles login command logic
func Login(debug bool) error {
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
	path, err := config.GetConfigPath()
	if err != nil {
		log.Error().Err(err).Msg("failed to get configuration file")
		return err
	}

	var cfg config.Config
	cfg.Mnemonics = mnemonics
	cfg.Network = network

	err = cfg.Save(path)
	if err != nil {
		log.Error().Err(err).Send()
		return err
	}
	log.Info().Msg("configuration saved")
	return nil
}
