package tfpluginclient

import (
	"context"

	gridDeployer "github.com/threefoldtech/grid3-go/deployer"
	"github.com/threefoldtech/grid3-go/graphql"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/grid_proxy_server/pkg/client"
	"github.com/threefoldtech/substrate-client"
)

// TFPluginClientInterface interface for tfPluginClient
type TFPluginClientInterface interface {
	DeployNetwork(ctx context.Context, znet *workloads.ZNet) error
	DeployDeployment(ctx context.Context, dl *workloads.Deployment) error
	DeployGatewayName(ctx context.Context, gw *workloads.GatewayNameProxy) error
	LoadVMFromGrid(nodeID uint32, name string, deploymentName string) (workloads.VM, error)
	LoadGatewayNameFromGrid(nodeID uint32, name string, deploymentName string) (workloads.GatewayNameProxy, error)
	ListContractsOfProjectName(projectName string) (graphql.Contracts, error)
	CancelContract(identity substrate.Identity, contractID uint64) error
	GetGridNetwork() string
	GetGridProxyClient() client.Client
	GetIdentity() substrate.Identity
}

// NewTFPluginClient returns new tfPluginClient given mnemonics and grid network
func NewTFPluginClient(mnemonics, network string) (tfPluginClient, error) {
	t, err := gridDeployer.NewTFPluginClient(mnemonics, "sr25519", network, "", "", "", true, false)
	if err != nil {
		return tfPluginClient{}, err
	}
	return tfPluginClient{
		&t,
	}, nil
}

// tfPluginClient wraps grid3-go tfPluginClient
type tfPluginClient struct {
	tfPluginClient *gridDeployer.TFPluginClient
}

// DeployNetwork deploys a network deployment to Threefold grid
func (t *tfPluginClient) DeployNetwork(ctx context.Context, znet *workloads.ZNet) error {
	return t.tfPluginClient.NetworkDeployer.Deploy(ctx, znet)
}

// DeployDeployment deploys a deployment to Threefold grid
func (t *tfPluginClient) DeployDeployment(ctx context.Context, dl *workloads.Deployment) error {
	return t.tfPluginClient.DeploymentDeployer.Deploy(ctx, dl)
}

// DeployNameGateway deploys a GatewayName deployment to Threefold grid
func (t *tfPluginClient) DeployGatewayName(ctx context.Context, gw *workloads.GatewayNameProxy) error {
	return t.tfPluginClient.GatewayNameDeployer.Deploy(ctx, gw)
}

// LoadVMFromGrid loads a VM from Threefold grid
func (t *tfPluginClient) LoadVMFromGrid(nodeID uint32, name string, deploymentName string) (workloads.VM, error) {
	return t.tfPluginClient.State.LoadVMFromGrid(nodeID, name, deploymentName)
}

// LoadGatewayNameFromGrid loads a GatewayName from Threefold grid
func (t *tfPluginClient) LoadGatewayNameFromGrid(nodeID uint32, name string, deploymentName string) (workloads.GatewayNameProxy, error) {
	return t.tfPluginClient.State.LoadGatewayNameFromGrid(nodeID, name, deploymentName)
}

// ListContractsOfProjectName returns contracts for a project name from Threefold grid
func (t *tfPluginClient) ListContractsOfProjectName(projectName string) (graphql.Contracts, error) {
	return t.tfPluginClient.ContractsGetter.ListContractsOfProjectName(projectName)
}

// CancelContract cancels a contract on Threefold grid
func (t *tfPluginClient) CancelContract(identity substrate.Identity, contractID uint64) error {
	return t.tfPluginClient.SubstrateConn.CancelContract(identity, contractID)
}

// GetGridNetwork returns the current grid network
func (t *tfPluginClient) GetGridNetwork() string {
	return t.tfPluginClient.Network
}

// GetGridProxyClient returns grid proxy client
func (t *tfPluginClient) GetGridProxyClient() client.Client {
	return t.tfPluginClient.GridProxyClient
}

// GetIdentity returns the identity of the user
func (t *tfPluginClient) GetIdentity() substrate.Identity {
	return t.tfPluginClient.Identity
}