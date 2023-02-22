// Package deployer for project deployment
package deployer

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	gridDeployer "github.com/threefoldtech/grid3-go/deployer"
)

// Deployer struct manages project deployment
type Deployer struct {
	tfPluginClient *gridDeployer.TFPluginClient

	repoURL     string
	projectName string

	logger zerolog.Logger
}

// NewDeployer return new project deployer
func NewDeployer(mnemonics, network string, repoURL string, logger zerolog.Logger) (Deployer, error) {
	rand.Seed(time.Now().UnixNano())

	tfPluginClient, err := gridDeployer.NewTFPluginClient(mnemonics, "sr25519", network, "", "", "", true, false)
	if err != nil {
		return Deployer{}, err
	}
	deployer := Deployer{
		tfPluginClient: &tfPluginClient,
		logger:         logger,
		repoURL:        repoURL,
	}

	projectName, err := deployer.getProjectName()
	if err != nil {
		return Deployer{}, err
	}
	deployer.projectName = projectName

	return deployer, nil
}

// Deploy deploys a project and map each port to a domain
func (d *Deployer) Deploy(ctx context.Context, ports []uint) (map[uint]string, error) {

	d.logger.Debug().Msg("getting nodes with free resources")

	node, err := findNode(d.tfPluginClient.GridProxyClient)
	if err != nil {
		return map[uint]string{}, errors.Wrapf(
			err,
			"failed to get a node with enough resources on network %s",
			d.tfPluginClient.Network,
		)
	}

	d.logger.Info().Msg("deploying network")
	network := buildNetwork(d.projectName, node)
	err = d.tfPluginClient.NetworkDeployer.Deploy(ctx, &network)
	if err != nil {
		return map[uint]string{}, errors.Wrapf(err, "could not deploy network %s on node %d", network.Name, node)
	}

	d.logger.Info().Msg("deploying vm")
	dl := buildDeployment(network.Name, d.projectName, d.repoURL, node)
	err = d.tfPluginClient.DeploymentDeployer.Deploy(ctx, &dl)
	if err != nil {
		return map[uint]string{}, errors.Wrapf(err, "could not deploy vm %s on node %d", dl.Name, node)
	}

	resVM, err := d.tfPluginClient.State.LoadVMFromGrid(node, dl.Name, dl.Name)
	if err != nil {
		return map[uint]string{}, errors.Wrapf(err, "could not load vm %s on node %d", dl.Name, node)
	}

	portlessBackend := buildPortlessBackend(resVM.ComputedIP)

	FQDNs := make(map[uint]string, 0)
	// TODO: deploy each gateway in a separate goroutine
	for _, port := range ports {
		backend := fmt.Sprintf("%s:%d", portlessBackend, port)
		d.logger.Info().Msgf("deploying gateway for port %d", port)
		gateway := buildGateway(backend, d.projectName, node)
		err := d.tfPluginClient.GatewayNameDeployer.Deploy(ctx, &gateway)
		if err != nil {
			return map[uint]string{}, errors.Wrapf(err, "could not deploy gateway %s on node %d", gateway.Name, node)
		}
		resGateway, err := d.tfPluginClient.State.LoadGatewayNameFromGrid(node, gateway.Name, gateway.Name)
		if err != nil {
			return map[uint]string{}, errors.Wrapf(err, "could not load gateway %s on node %d", gateway.Name, node)
		}
		FQDNs[port] = resGateway.FQDN
	}

	d.logger.Info().Msg("Project Deployed!")

	return FQDNs, nil
}

// Destroy destroys all the contracts of a project
func (d *Deployer) Destroy() error {
	d.logger.Info().Msgf("canceling contracts for project %s", d.projectName)
	contracts, err := d.tfPluginClient.ContractsGetter.ListContractsOfProjectName(d.projectName)
	if err != nil {
		return errors.Wrapf(err, "could not load contracts for project %s", d.projectName)
	}
	contractsSlice := append(contracts.NameContracts, contracts.NodeContracts...)
	for _, contract := range contractsSlice {
		contractID, err := strconv.ParseUint(contract.ContractID, 0, 64)
		if err != nil {
			return errors.Wrapf(err, "could not parse contract %s into uint64", contract.ContractID)
		}
		d.logger.Debug().Msgf("canceling contract %d", contractID)

		err = d.tfPluginClient.SubstrateConn.CancelContract(d.tfPluginClient.Identity, contractID)
		if err != nil {
			return errors.Wrapf(err, "could not cancel contract %d", contractID)
		}
	}
	d.logger.Info().Msg("Project Destroyed!")
	return nil
}

func (d *Deployer) getProjectName() (string, error) {

	splitURL := strings.Split(string(d.repoURL), "/")
	projectName, _, found := strings.Cut(splitURL[len(splitURL)-1], ".git")
	if !found {
		return "", fmt.Errorf("couldn't get project name")
	}
	return projectName, nil
}
