package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/rawdaGastan/gridify/internal/config"
	"github.com/rawdaGastan/gridify/internal/deployer"
)

func Destroy(showLogs bool) error {
	config, err := config.LoadConfigData()
	if err != nil {
		return err
	}
	deployer, err := deployer.NewDeployer(config.Mnemonics, config.Network, showLogs)
	if err != nil {
		return err
	}
	repoURL, err := exec.Command("git", "config", "--get", "remote.origin.url").Output()
	if err != nil {
		return err
	}
	splitURL := strings.Split(string(repoURL), "/")
	projectName, _, found := strings.Cut(splitURL[len(splitURL)-1], ".git")
	if !found {
		return fmt.Errorf("couldn't get project name")
	}
	err = deployer.Destroy(string(projectName))
	if err != nil {
		return err
	}
	fmt.Println("Project Destroyed!")
	return nil
}
