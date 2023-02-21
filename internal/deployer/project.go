// Package deployer for project deployment
package deployer

import (
	"fmt"
	"strings"
)

func getProjectName(repoURL string) (string, error) {

	splitURL := strings.Split(string(repoURL), "/")
	projectName, _, found := strings.Cut(splitURL[len(splitURL)-1], ".git")
	if !found {
		return "", fmt.Errorf("couldn't get project name")
	}
	return projectName, nil
}
