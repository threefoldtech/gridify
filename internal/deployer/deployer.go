package deployer

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

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

func (d *deployer) Deploy(ctx context.Context, repoURL, projectName string, ports []string) ([]string, error) {
	rand.Seed(time.Now().UnixNano())
	randomString := randString(10)
	networkName := fmt.Sprintf("gnet%s", randomString)
	vmName := fmt.Sprintf("gvm%s", randomString)
	FQDNs := make([]string, 0)
	network := workloads.ZNet{
		Name:        networkName,
		Description: "network for testing",
		Nodes:       []uint32{2},
		IPRange: gridtypes.NewIPNet(net.IPNet{
			IP:   net.IPv4(10, 20, 0, 0),
			Mask: net.CIDRMask(16, 32),
		}),
		AddWGAccess:  false,
		SolutionType: projectName,
	}

	vm := workloads.VM{
		Name:       vmName,
		Flist:      "https://hub.grid.tf/aelawady.3bot/abdulrahmanelawady-gridify-test-latest.flist",
		CPU:        2,
		PublicIP:   true,
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
		return FQDNs, err
	}

	dl := workloads.NewDeployment(vm.Name, 2, projectName, nil, network.Name, nil, nil, []workloads.VM{vm}, nil)
	err = d.tfPluginClient.DeploymentDeployer.Deploy(ctx, &dl)

	if err != nil {
		return FQDNs, err
	}
	resVM, err := d.tfPluginClient.State.LoadVMFromGrid(2, vm.Name)
	if err != nil {
		return FQDNs, err
	}
	publicIP := strings.Split(resVM.ComputedIP, "/")[0]

	for _, port := range ports {
		domainName := randString(10)
		backend := fmt.Sprintf("http://%s:%s", publicIP, port)
		gw := workloads.GatewayNameProxy{
			NodeID:         2,
			Name:           domainName,
			TLSPassthrough: false,
			Backends:       []zos.Backend{zos.Backend(backend)},
			SolutionType:   projectName,
		}

		err = d.tfPluginClient.GatewayNameDeployer.Deploy(ctx, &gw)
		if err != nil {
			return FQDNs, nil
		}
		gwRes, err := d.tfPluginClient.State.LoadGatewayNameFromGrid(2, gw.Name)
		if err != nil {
			return FQDNs, err
		}
		FQDNs = append(FQDNs, gwRes.FQDN)
	}

	return FQDNs, nil
}

func (d *deployer) Destroy(projectName string) error {
	contracts, err := d.tfPluginClient.ContractsGetter.ListContractsOfProjectName(projectName)
	if err != nil {
		return err
	}
	contractsSlice := make([]gridDeployer.Contract, 0)
	contractsSlice = append(contractsSlice, contracts.NameContracts...)
	contractsSlice = append(contractsSlice, contracts.NodeContracts...)
	for _, contract := range contractsSlice {
		contractID, err := strconv.ParseUint(contract.ContractID, 0, 64)
		if err != nil {
			return err
		}
		err = d.tfPluginClient.SubstrateConn.CancelContract(d.tfPluginClient.Identity, contractID)
		if err != nil {
			return err
		}
	}
	return nil
}
