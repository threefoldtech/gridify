package deployer

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	gridDeployer "github.com/threefoldtech/grid3-go/deployer"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes"
)

type deployer struct {
	tfPluginClient *gridDeployer.TFPluginClient
}

func NewDeployer() deployer {
	mnemonics := os.Getenv("MNEMONICS")
	log.Printf("mnemonics: %s", mnemonics)

	network := os.Getenv("NETWORK")
	log.Printf("network: %s", network)

	tfPluginClient, err := gridDeployer.NewTFPluginClient(mnemonics, "sr25519", network, "", "", true, false)
	if err != nil {
		return deployer{}
	}
	return deployer{tfPluginClient: &tfPluginClient}
}

func (b *deployer) Deploy(ctx context.Context, repoURL string) (string, error) {
	network := workloads.ZNet{
		Name:        "gridifyNetwork",
		Description: "network for testing",
		Nodes:       []uint32{3},
		IPRange: gridtypes.NewIPNet(net.IPNet{
			IP:   net.IPv4(10, 20, 0, 0),
			Mask: net.CIDRMask(16, 32),
		}),
		AddWGAccess: false,
	}

	vm := workloads.VM{
		Name:       "gridifyVM",
		Flist:      "https://hub.grid.tf/aelawady.3bot/abdulrahmanelawady-gridify-test-latest.flist",
		CPU:        2,
		PublicIP:   true,
		Planetary:  true,
		Memory:     1024,
		RootfsSize: 20 * 1024,
		Entrypoint: "/sbin/zinit init",
		EnvVars: map[string]string{
			"REPO_URL": repoURL,
		},
		NetworkName: network.Name,
	}

	err := b.tfPluginClient.NetworkDeployer.Deploy(ctx, &network)
	if err != nil {
		return "", err
	}

	dl := workloads.NewDeployment("gridifyVM", 3, "", nil, network.Name, nil, nil, []workloads.VM{vm}, nil)
	err = b.tfPluginClient.DeploymentDeployer.Deploy(ctx, &dl)

	if err != nil {
		return "", err
	}

	resVM, err := b.tfPluginClient.StateLoader.LoadVMFromGrid(3, "gridifyVM")
	if err != nil {
		return "", err
	}

	publicIP := strings.Split(resVM.ComputedIP, "/")[0]
	if !testConnection(publicIP, "22") {
		return "", fmt.Errorf("public ip %s is not reachable", publicIP)
	}

	return publicIP, nil
}

func testConnection(addr string, port string) bool {
	for t := time.Now(); time.Since(t) < 3*time.Minute; {
		con, err := net.DialTimeout("tcp", net.JoinHostPort(addr, port), time.Second*12)
		if err == nil {
			con.Close()
			return true
		}
	}
	return false
}
