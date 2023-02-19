// Package cmd for handling command line arguments
package cmd

import (
	"github.com/rawdaGastan/gridify/internal/config"
	"github.com/rawdaGastan/gridify/internal/deployer"
)

func login(mnemonics, network string, showLogs bool) error {
	// TODO: better verification
	_, err := deployer.NewDeployer(mnemonics, network, showLogs)
	if err != nil {
		return err
	}
	return config.SaveConfigData(mnemonics, network)
}
