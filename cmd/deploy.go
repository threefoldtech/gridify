package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rawdaGastan/gridify/internal/config"
	"github.com/rawdaGastan/gridify/internal/deployer"
)

func Deploy(ports string, showLogs bool) error {
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
	portsSlice := strings.Split(ports, ",")

	FQDNs, err := deployer.Deploy(context.Background(), string(repoURL), projectName, portsSlice)
	if err != nil {
		return err
	}

	fmt.Println("Project Deployed!")
	for port, FQDN := range FQDNs {
		fmt.Printf("%s: %s\n", port, FQDN)
	}
	return nil
}
