package deployer

import (
	"context"
	"fmt"
	"net"
	"os"

	gridDeployer "github.com/threefoldtech/grid3-go/deployer"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
)

type deployer struct {
	tfPluginClient *gridDeployer.TFPluginClient
}

func NewDeployer(mnemonics string) (deployer, error) {

	network := os.Getenv("NETWORK")

	if network == "" {
		network = "qa"
	}
	tfPluginClient, err := gridDeployer.NewTFPluginClient(mnemonics, "sr25519", network, "", "", true, false)
	if err != nil {
		return deployer{}, err
	}
	return deployer{tfPluginClient: &tfPluginClient}, nil
}

func (d *deployer) Deploy(ctx context.Context, repoURL, port string) (string, error) {
	randomString := randString(10)
	networkName := fmt.Sprintf("gnet%s", randomString)
	vmName := fmt.Sprintf("gvm%s", randomString)
	network := workloads.ZNet{
		Name:        networkName,
		Description: "network for testing",
		Nodes:       []uint32{2},
		IPRange: gridtypes.NewIPNet(net.IPNet{
			IP:   net.IPv4(10, 20, 0, 0),
			Mask: net.CIDRMask(16, 32),
		}),
		AddWGAccess: false,
	}

	vm := workloads.VM{
		Name:       vmName,
		Flist:      "https://hub.grid.tf/aelawady.3bot/abdulrahmanelawady-gridify-test-latest.flist",
		CPU:        2,
		Planetary:  true,
		Memory:     1024,
		RootfsSize: 20 * 1024,
		Entrypoint: "/init.sh",
		EnvVars: map[string]string{
			"REPO_URL": repoURL,
		},
		NetworkName: network.Name,
	}

	err := d.tfPluginClient.NetworkDeployer.Deploy(ctx, &network)
	if err != nil {
		return "", err
	}

	dl := workloads.NewDeployment(vm.Name, 2, "", nil, network.Name, nil, nil, []workloads.VM{vm}, nil)
	err = d.tfPluginClient.DeploymentDeployer.Deploy(ctx, &dl)

	if err != nil {
		return "", err
	}
	resVM, err := d.tfPluginClient.StateLoader.LoadVMFromGrid(2, vm.Name)
	if err != nil {
		return "", err
	}

	backend := fmt.Sprintf("http://[%s]:%s", resVM.YggIP, port)
	gw := workloads.GatewayNameProxy{
		NodeID:         2,
		Name:           randomString,
		TLSPassthrough: false,
		Backends:       []zos.Backend{zos.Backend(backend)},
	}

	err = d.tfPluginClient.GatewayNameDeployer.Deploy(ctx, &gw)
	if err != nil {
		return "", nil
	}
	gwRes, err := d.tfPluginClient.StateLoader.LoadGatewayNameFromGrid(2, gw.Name)
	if err != nil {
		return "", err
	}

	return gwRes.FQDN, nil
}
