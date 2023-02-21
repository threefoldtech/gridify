// Package deployer for project deployment
package deployer

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	gridDeployer "github.com/threefoldtech/grid3-go/deployer"
	"github.com/threefoldtech/grid3-go/graphql"
	"github.com/threefoldtech/grid_proxy_server/pkg/types"
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

	projectName, err := getProjectName(repoURL)
	if err != nil {
		return Deployer{}, err
	}

	return Deployer{
			tfPluginClient: &tfPluginClient,
			logger:         logger,
			repoURL:        repoURL,
			projectName:    projectName,
		},
		nil
}

// Deploy deploys a project and map each port to a domain
func (d *Deployer) Deploy(ctx context.Context, ports []uint) (map[uint]string, error) {

	d.logger.Debug().Msg("getting nodes with free resources")

	filter := constructNodeFilter()

	nodes, _, err := d.tfPluginClient.GridProxyClient.Nodes(filter, types.Limit{})
	if err != nil {
		return map[uint]string{}, errors.Wrapf(
			err,
			"failed to get a node with enough resources on network %s",
			d.tfPluginClient.Network,
		)
	}
	if len(nodes) == 0 {
		return map[uint]string{}, fmt.Errorf("no node with free resources available on %s", d.tfPluginClient.Network)
	}

	node := uint32(nodes[0].NodeID)

	d.logger.Info().Msg("deploying network")
	network := constructNetwork(d.projectName, node)
	err = d.tfPluginClient.NetworkDeployer.Deploy(ctx, &network)
	if err != nil {
		return map[uint]string{}, errors.Wrapf(err, "could not deploy network %s on node %d", network.Name, node)
	}

	d.logger.Info().Msg("deploying vm")
	dl := constructDeployment(network.Name, d.projectName, d.repoURL, node)
	err = d.tfPluginClient.DeploymentDeployer.Deploy(ctx, &dl)
	if err != nil {
		return map[uint]string{}, errors.Wrapf(err, "could not deploy vm %s on node %d", dl.Name, node)
	}

	resVM, err := d.tfPluginClient.State.LoadVMFromGrid(node, dl.Name, dl.Name)
	if err != nil {
		return map[uint]string{}, errors.Wrapf(err, "could not load vm %s on node %d", dl.Name, node)
	}

	portlessBackend := constructPortlessBackend(resVM.ComputedIP)
	if err != nil {
		return map[uint]string{}, errors.Wrapf(err, "could not construct backend for vm %s on node %d", dl.Name, node)
	}

	FQDNs := make(map[uint]string, 0)
	// TODO: deploy each gateway in a separate goroutine
	for _, port := range ports {
		backend := fmt.Sprintf("%s:%d", portlessBackend, port)
		d.logger.Info().Msgf("deploying gateway for port %d", port)
		gateway := constructGateway(backend, d.projectName, node)
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
	contractsSlice := make([]graphql.Contract, 0)
	contractsSlice = append(contractsSlice, contracts.NameContracts...)
	contractsSlice = append(contractsSlice, contracts.NodeContracts...)
	// TODO: cancel each contract in a separate goroutine
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
