// Package deployer for project deployment
package deployer

import (
	"fmt"
	"net"
	"strings"

	"github.com/pkg/errors"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/grid_proxy_server/pkg/client"
	"github.com/threefoldtech/grid_proxy_server/pkg/types"
	"github.com/threefoldtech/zos/pkg/gridtypes"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
)

var (
	vmFlist      = "https://hub.grid.tf/aelawady.3bot/abdulrahmanelawady-gridify-test-latest.flist"
	vmCPU        = 2
	vmMemory     = 2 * 1024
	vmRootfsSize = 5 * 1024
	vmEntrypoint = "/init.sh"
	vmPublicIP   = true
	vmPlanetary  = true
)

func constructNodeFilter() types.NodeFilter {
	nodeStatus := "up"
	freeMRU := uint64(2)
	freeHRU := uint64(5)
	freeIPs := uint64(1)
	domain := true

	filter := types.NodeFilter{
		FarmIDs: []uint64{1},
		Status:  &nodeStatus,
		FreeMRU: &freeMRU,
		FreeHRU: &freeHRU,
		FreeIPs: &freeIPs,
		Domain:  &domain,
	}
	return filter
}

func findNode(gridProxyClient client.Client) (uint32, error) {
	filter := constructNodeFilter()
	nodes, _, err := gridProxyClient.Nodes(filter, types.Limit{})
	if err != nil {
		return 0, err
	}
	if len(nodes) == 0 {
		return 0, errors.New("no node with free resources available")
	}

	node := uint32(nodes[0].NodeID)
	return node, nil
}

func constructNetwork(projectName string, node uint32) workloads.ZNet {
	networkName := randString(10)
	network := workloads.ZNet{
		Name:  networkName,
		Nodes: []uint32{node},
		IPRange: gridtypes.NewIPNet(net.IPNet{
			IP:   net.IPv4(10, 20, 0, 0),
			Mask: net.CIDRMask(16, 32),
		}),
		SolutionType: projectName,
	}
	return network
}

func constructDeployment(networkName, projectName, repoURL string, node uint32) workloads.Deployment {
	vmName := randString(10)
	vm := workloads.VM{
		Name:       vmName,
		Flist:      vmFlist,
		CPU:        vmCPU,
		PublicIP:   vmPublicIP,
		Planetary:  vmPlanetary,
		Memory:     vmMemory,
		RootfsSize: vmRootfsSize,
		Entrypoint: vmEntrypoint,
		EnvVars: map[string]string{
			"REPO_URL": repoURL,
		},
		NetworkName: networkName,
	}

	dl := workloads.NewDeployment(vm.Name, node, projectName, nil, networkName, nil, nil, []workloads.VM{vm}, nil)
	return dl
}

func constructGateway(backend, projectName string, node uint32) workloads.GatewayNameProxy {
	subdomain := randString(10)
	gateway := workloads.GatewayNameProxy{
		NodeID:       node,
		Name:         subdomain,
		Backends:     []zos.Backend{zos.Backend(backend)},
		SolutionType: projectName,
	}
	return gateway
}

func constructPortlessBackend(ip string) string {
	publicIP := strings.Split(ip, "/")[0]
	backend := fmt.Sprintf("http://%s", publicIP)
	return backend
}
