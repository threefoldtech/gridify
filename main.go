package main

import (
	"context"
	"flag"
	"fmt"
	"os/exec"

	"github.com/rawdaGastan/gridify/internal/deployer"
)

func main() {
	cmd := flag.String("command", "", "")
	flag.Parse()
	fmt.Println(*cmd)
	repoURLCmd := exec.Command("git", "config", "--get", "remote.origin.url")
	repoURL, err := repoURLCmd.Output()
	if err != nil {
		panic(err)
	}
	deployer := deployer.NewDeployer()
	vmIP, err := deployer.Deploy(context.Background(), string(repoURL))
	if err != nil {
		panic(err)
	}
	fmt.Println(vmIP)
}
