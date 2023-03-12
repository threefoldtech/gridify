// Package deployer for project deployment
package deployer

import (
	"fmt"
	"net"
	"strings"

	"github.com/rawdaGastan/gridify/internal/tfplugin"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/grid_proxy_server/pkg/types"
	"github.com/threefoldtech/zos/pkg/gridtypes"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
)

var (
	vmFlist                 = "https://hub.grid.tf/aelawady.3bot/abdulrahmanelawady-gridify-test-latest.flist"
	eco                     = "eco"
	standard                = "standard"
	performance             = "performance"
	vmSpecs                 = []string{eco, standard, performance}
	ecoVMCPU                = 1
	ecoVMMemory             = 2 // GB
	ecoVMRootfsSize         = 5 // GB
	standardVMCPU           = 2
	standardVMMemory        = 4  // GB
	standardVMRootfsSize    = 10 // GB
	performanceVMCPU        = 4
	performanceVMMemory     = 8  // GB
	performanceVMRootfsSize = 15 // GB
	vmEntrypoint            = "/init.sh"
	vmPublicIP              = true
	vmPlanetary             = true
)

func isValidVMSpec(spec string) bool {
	for _, s := range vmSpecs {
		if s == spec {
			return true
		}
	}
	return false
}

func buildNodeFilter(vmSpec string) types.NodeFilter {
	nodeStatus := "up"
	freeMRU := uint64(0)
	freeSRU := uint64(0)
	freeIPs := uint64(0)
	if vmPublicIP {
		freeIPs = 1
	}
	domain := true

	switch vmSpec {
	case eco:
		freeMRU = uint64(ecoVMMemory)
		freeSRU = uint64(ecoVMRootfsSize)
	case standard:
		freeMRU = uint64(standardVMMemory)
		freeSRU = uint64(standardVMRootfsSize)
	case performance:
		freeMRU = uint64(performanceVMMemory)
		freeSRU = uint64(performanceVMRootfsSize)
	}
	filter := types.NodeFilter{
		FarmIDs: []uint64{1},
		Status:  &nodeStatus,
		FreeMRU: &freeMRU,
		FreeSRU: &freeSRU,
		FreeIPs: &freeIPs,
		Domain:  &domain,
	}
	return filter
}

func findNode(vmSpec string, tfPluginClient tfplugin.TFPluginClientInterface) (uint32, error) {
	filter := buildNodeFilter(vmSpec)
	nodes, _, err := tfPluginClient.FilterNodes(filter, types.Limit{})
	if err != nil {
		return 0, err
	}
	if len(nodes) == 0 {
		return 0, fmt.Errorf(
			"no node with free resources available using node filter: farmIDs: %v, mru: %d, sru: %d, freeips: %d, domain: %t",
			filter.FarmIDs,
			*filter.FreeMRU,
			*filter.FreeSRU,
			*filter.FreeIPs,
			*filter.Domain,
		)
	}

	node := uint32(nodes[0].NodeID)
	return node, nil
}

func buildNetwork(projectName string, node uint32) workloads.ZNet {
	networkName := randName(10)
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

func buildDeployment(vmSpec, networkName, projectName, repoURL string, node uint32) workloads.Deployment {
	vmName := randName(10)
	vm := workloads.VM{
		Name:       vmName,
		Flist:      vmFlist,
		PublicIP:   vmPublicIP,
		Planetary:  vmPlanetary,
		Entrypoint: vmEntrypoint,
		EnvVars: map[string]string{
			"REPO_URL": repoURL,
		},
		NetworkName: networkName,
	}

	switch vmSpec {
	case eco:
		vm.CPU = ecoVMCPU
		vm.Memory = ecoVMMemory * 1024
		vm.RootfsSize = ecoVMRootfsSize * 1024
	case standard:
		vm.CPU = standardVMCPU
		vm.Memory = standardVMMemory * 1024
		vm.RootfsSize = standardVMRootfsSize * 1024
	case performance:
		vm.CPU = performanceVMCPU
		vm.Memory = performanceVMMemory * 1024
		vm.RootfsSize = performanceVMRootfsSize * 1024
	}

	dl := workloads.NewDeployment(vm.Name, node, projectName, nil, networkName, nil, nil, []workloads.VM{vm}, nil)
	return dl
}

func buildGateway(backend, projectName string, node uint32) workloads.GatewayNameProxy {
	subdomain := randName(10)
	gateway := workloads.GatewayNameProxy{
		NodeID:       node,
		Name:         subdomain,
		Backends:     []zos.Backend{zos.Backend(backend)},
		SolutionType: projectName,
	}
	return gateway
}

func buildPortlessBackend(ip string) string {
	publicIP := strings.Split(ip, "/")[0]
	backend := fmt.Sprintf("http://%s", publicIP)
	return backend
}
