// Package cmd for handling commands
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	bip39 "github.com/cosmos/go-bip39"
	"github.com/rawdaGastan/gridify/internal/config"
	"github.com/rs/zerolog/log"
)

// Login handles login command logic
func Login(debug bool) {
	scanner := bufio.NewReader(os.Stdin)

	fmt.Print("Please enter your mnemonics: ")
	mnemonics, err := scanner.ReadString('\n')
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read mnemonics")
	}
	mnemonics = strings.TrimSpace(mnemonics)
	if !bip39.IsMnemonicValid(mnemonics) {
		log.Fatal().Msg("failed to validate mnemonics")
	}

	fmt.Print("Please enter grid network (main,test): ")
	network, err := scanner.ReadString('\n')
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read grid network")
	}
	network = strings.TrimSpace(network)

	if network != "dev" && network != "qa" && network != "test" && network != "main" {
		log.Fatal().Msg("invalid grid network, must be one of: dev, test, qa and main")
	}
	path, err := config.GetConfigPath()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get configuration file")
	}

	var cfg config.Config
	cfg.Mnemonics = mnemonics
	cfg.Network = network

	err = cfg.Save(path)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Info().Msg("configuration saved")
}
