// Package cmd for handling command line arguments
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rawdaGastan/gridify/internal/config"
	"github.com/spf13/cobra"
	gridDeployer "github.com/threefoldtech/grid3-go/deployer"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login with mnemonics on a grid network",
	Run: func(cmd *cobra.Command, args []string) {
		scanner := bufio.NewReader(os.Stdin)

		fmt.Print("Please enter your mnemonics: ")
		mnemonics, err := scanner.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to read mnemonics")
			os.Exit(1)
		}
		mnemonics = strings.TrimSpace(mnemonics)

		fmt.Print("Please enter grid network: ")
		network, err := scanner.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to read network")
			os.Exit(1)
		}
		network = strings.TrimSpace(network)

		_, err = gridDeployer.NewTFPluginClient(mnemonics, "sr25519", network, "", "", "", true, false)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		err = config.SaveConfigData(mnemonics, network)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
