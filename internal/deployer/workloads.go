// Package deployer for project deployment
package deployer

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
)

func (d *Deployer) deployNetwork(ctx context.Context, networkName, projectName string, node uint32) error {
	network := workloads.ZNet{
		Name:  networkName,
		Nodes: []uint32{node},
		IPRange: gridtypes.NewIPNet(net.IPNet{
			IP:   net.IPv4(10, 20, 0, 0),
			Mask: net.CIDRMask(16, 32),
		}),
		SolutionType: projectName,
	}
	return d.tfPluginClient.NetworkDeployer.Deploy(ctx, &network)

}
func (d *Deployer) deployVM(ctx context.Context, networkName, projectName, repoURL, vmName string, node uint32) error {

	vm := workloads.VM{
		Name:       vmName,
		Flist:      "https://hub.grid.tf/aelawady.3bot/abdulrahmanelawady-gridify-test-latest.flist",
		CPU:        2,
		PublicIP:   true,
		Planetary:  true,
		Memory:     2 * 1024,
		RootfsSize: 5 * 1024,
		Entrypoint: "/init.sh",
		EnvVars: map[string]string{
			"REPO_URL": repoURL,
		},
		NetworkName: networkName,
	}

	dl := workloads.NewDeployment(vm.Name, 2, projectName, nil, networkName, nil, nil, []workloads.VM{vm}, nil)
	return d.tfPluginClient.DeploymentDeployer.Deploy(ctx, &dl)

}

func (d *Deployer) constructPortlessBackend(ctx context.Context, node uint32, vmName string) (string, error) {
	resVM, err := d.tfPluginClient.State.LoadVMFromGrid(node, vmName)
	if err != nil {
		return "", err
	}
	publicIP := strings.Split(resVM.ComputedIP, "/")[0]
	backend := fmt.Sprintf("http://%s", publicIP)
	return backend, nil
}
func (d *Deployer) deployGateway(ctx context.Context, backend, projectName, subdomain string, node uint32) error {
	gw := workloads.GatewayNameProxy{
		NodeID:       node,
		Name:         subdomain,
		Backends:     []zos.Backend{zos.Backend(backend)},
		SolutionType: projectName,
	}

	return d.tfPluginClient.GatewayNameDeployer.Deploy(ctx, &gw)
}

func (d *Deployer) getFQDN(ctx context.Context, subdomain string, node uint32) (string, error) {
	gwRes, err := d.tfPluginClient.State.LoadGatewayNameFromGrid(node, subdomain)
	return gwRes.FQDN, err

}
