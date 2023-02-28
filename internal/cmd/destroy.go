// Package cmd for handling commands
package cmd

import (
	"os"

	"github.com/rawdaGastan/gridify/internal/config"
	"github.com/rawdaGastan/gridify/internal/deployer"
	"github.com/rawdaGastan/gridify/internal/repository"
	"github.com/rawdaGastan/gridify/internal/tfplugin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Destroy handles destroy command logic
func Destroy(debug bool) {
	path, err := config.GetConfigPath()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get configuration file")
	}

	var cfg config.Config
	err = cfg.Load(path)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration try to login again using gridify login")
	}

	repoURL, err := repository.GetRepositoryURL(".")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get remote repository url")
	}

	logLevel := zerolog.InfoLevel
	if debug {
		logLevel = zerolog.DebugLevel
	}

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).Level(logLevel).
		With().
		Timestamp().
		Logger()

	tfPluginClient, err := tfplugin.NewTFPluginClient(cfg.Mnemonics, cfg.Network)
	if err != nil {
		log.Fatal().
			Err(err).
			Msgf("failed to get threefold plugin client using mnemonics: '%s' on grid network '%s'", cfg.Mnemonics, cfg.Network)
	}
	deployer, err := deployer.NewDeployer(&tfPluginClient, string(repoURL), logger)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	err = deployer.Destroy()
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
