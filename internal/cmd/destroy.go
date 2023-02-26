// Package cmd for handling commands
package cmd

import (
	"os"

	"github.com/rawdaGastan/gridify/internal/config"
	"github.com/rawdaGastan/gridify/internal/deployer"
	"github.com/rawdaGastan/gridify/internal/repository"
	"github.com/rawdaGastan/gridify/internal/tfpluginclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Destroy handles destroy command logic
func Destroy(debug bool) error {
	path, err := config.GetConfigPath()
	if err != nil {
		log.Error().Err(err).Msg("failed to get configuration file")
		return err
	}

	var cfg config.Config
	err = cfg.Load(path)
	if err != nil {
		log.Error().Err(err).Msg("failed to load configuration try to login again using gridify login")
		return err
	}

	repoURL, err := repository.GetRepositoryURL(".")
	if err != nil {
		log.Error().Err(err).Msg("failed to get remote repository url")
		return err
	}

	logLevel := zerolog.InfoLevel
	if debug {
		logLevel = zerolog.DebugLevel
	}

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).Level(logLevel).
		With().
		Timestamp().
		Logger()

	tfPluginClient, err := tfpluginclient.NewTFPluginClient(cfg.Mnemonics, cfg.Network)
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to get threefold plugin client using mnemonics: '%s' on grid network '%s'", cfg.Mnemonics, cfg.Network)
		return err
	}
	deployer, err := deployer.NewDeployer(&tfPluginClient, string(repoURL), logger)
	if err != nil {
		log.Error().Err(err).Send()
		return err
	}

	err = deployer.Destroy()
	if err != nil {
		log.Error().Err(err).Send()
		return err
	}
	return nil
}
