package deployer

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	gridDeployer "github.com/threefoldtech/grid3-go/deployer"
)

type deployer struct {
	tfPluginClient *gridDeployer.TFPluginClient

	logger zerolog.Logger
}

func NewDeployer(mnemonics, network string, showLogs bool) (deployer, error) {
	rand.Seed(time.Now().UnixNano())

	tfPluginClient, err := gridDeployer.NewTFPluginClient(mnemonics, "sr25519", network, "", "", true, false)
	if err != nil {
		return deployer{}, err
	}
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()
	if !showLogs {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: io.Discard})
	}

	return deployer{
			tfPluginClient: &tfPluginClient,
			logger:         logger,
		},
		nil
}

func (d *deployer) Deploy(ctx context.Context, repoURL, projectName string, ports []string) (map[string]string, error) {

	randomString := randString(10)

	d.logger.Info().Msg("getting nodes with free resources")
	filter := gridDeployer.NodeFilter{
		CRU:       2,
		MRU:       2,
		HRU:       5,
		PublicIPs: true,
		Gateway:   true,
		Status:    "up",
	}
	nodes, err := gridDeployer.FilterNodes(filter, gridDeployer.RMBProxyURLs[d.tfPluginClient.Network])
	if err != nil {
		return map[string]string{}, errors.Wrapf(
			err,
			"failed to get a node with enough resources on network %s",
			d.tfPluginClient.Network,
		)
	}
	if len(nodes) == 0 {
		return map[string]string{}, fmt.Errorf(
			"no available node with enough resources on network %s",
			d.tfPluginClient.Network,
		)
	}

	networkName := fmt.Sprintf("gnet%s", randomString)
	d.logger.Info().Msgf("deploying network '%s'", networkName)
	err = d.deployNetwork(ctx, networkName, projectName, nodes[0])
	if err != nil {
		return map[string]string{}, errors.Wrapf(err, "could not deploy network %s on node %d", networkName, nodes[0])
	}

	vmName := fmt.Sprintf("gvm%s", randomString)
	d.logger.Info().Msgf("deploying vm '%s'", vmName)
	err = d.deployVM(ctx, networkName, projectName, repoURL, vmName, nodes[0])
	if err != nil {
		return map[string]string{}, errors.Wrapf(err, "could not deploy vm %s on node %d", vmName, nodes[0])
	}

	portlessBackend, err := d.constructPortlessBackend(ctx, nodes[0], vmName)
	if err != nil {
		return map[string]string{}, errors.Wrapf(err, "could not construct backend for vm %s on node %d", vmName, nodes[0])
	}

	FQDNs := make(map[string]string, 0)
	// TODO: deploy each gateway in a separate goroutine
	for _, port := range ports {
		subdomain := randString(10)
		backend := fmt.Sprintf("%s:%s", portlessBackend, port)
		d.logger.Info().Msgf("deploying gateway '%s' for port %s", subdomain, port)
		err := d.deployGateway(ctx, backend, projectName, subdomain, nodes[0])
		if err != nil {
			return map[string]string{}, errors.Wrapf(err, "could not deploy gateway %s on node %d", subdomain, nodes[0])
		}
		FQDN, err := d.getFQDN(ctx, subdomain, nodes[0])
		if err != nil {
			return map[string]string{}, errors.Wrapf(err, "could not get FQDN for gateway %s on node %d", subdomain, nodes[0])
		}
		FQDNs[port] = FQDN
	}

	return FQDNs, nil
}

func (d *deployer) Destroy(projectName string) error {
	d.logger.Info().Msgf("getting contracts for project %s", projectName)
	contracts, err := d.tfPluginClient.ContractsGetter.ListContractsOfProjectName(projectName)
	if err != nil {
		return errors.Wrapf(err, "could not load contracts for project %s", projectName)
	}
	contractsSlice := make([]gridDeployer.Contract, 0)
	contractsSlice = append(contractsSlice, contracts.NameContracts...)
	contractsSlice = append(contractsSlice, contracts.NodeContracts...)
	// TODO: cancel each contract in a separate goroutine
	for _, contract := range contractsSlice {
		contractID, err := strconv.ParseUint(contract.ContractID, 0, 64)
		if err != nil {
			return errors.Wrapf(err, "could not parse contract %s into uint64", contract)
		}
		d.logger.Info().Msgf("canceling contract %d", contractID)

		err = d.tfPluginClient.SubstrateConn.CancelContract(d.tfPluginClient.Identity, contractID)
		if err != nil {
			return errors.Wrapf(err, "could not cancel contract %d", contractID)
		}
	}
	return nil
}
