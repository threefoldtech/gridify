package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rawdaGastan/gridify/internal/deployer"
)

func HandleDeploy(mnemonics, ports string) error {
	deployer, err := deployer.NewDeployer(mnemonics)
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
	fmt.Println(FQDNs)
	return err
}

func HandleDestroy(mnemonics string) error {

	deployer, err := deployer.NewDeployer(mnemonics)
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
	return deployer.Destroy(string(projectName))
}