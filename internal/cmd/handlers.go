package cmd

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/rawdaGastan/gridify/internal/deployer"
)

func HandleDeploy(mnemonics, port string) error {
	deployer, err := deployer.NewDeployer(mnemonics)
	if err != nil {
		return err
	}
	repoURLCmd := exec.Command("git", "config", "--get", "remote.origin.url")
	repoURL, err := repoURLCmd.Output()
	if err != nil {
		return err
	}
	FQDN, err := deployer.Deploy(context.Background(), string(repoURL), port)
	if err != nil {
		return err
	}
	fmt.Println(FQDN)
	return err
}

func HandleDestroy(mnemonics string) error {
	return nil
}
