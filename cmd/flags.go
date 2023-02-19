package cmd

import (
	"flag"
	"fmt"
	"os"
)

func Run() {
	var showLogs bool

	// set deploy command
	deployCMD := flag.NewFlagSet("deploy", flag.ExitOnError)
	ports := deployCMD.String("ports", "", "ports to forward the FQDN to")
	deployCMD.BoolVar(&showLogs, "show-logs", true, "show command logs")

	deployCMD.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gridify deploy [options]\n\n")
		deployCMD.PrintDefaults()
	}

	// set destroy command
	destroyCMD := flag.NewFlagSet("destroy", flag.ExitOnError)
	destroyCMD.BoolVar(&showLogs, "show-logs", true, "show command logs")
	destroyCMD.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gridify destroy\n\n")
	}

	// set login command
	loginCMD := flag.NewFlagSet("login", flag.ExitOnError)
	mnemonics := loginCMD.String("mnemonics", "", "mnemonics used for authentication")
	network := loginCMD.String("network", "", "grid network, one of: dev, test, qa and main")
	loginCMD.BoolVar(&showLogs, "show-logs", true, "show command logs")
	loginCMD.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gridify login [options]\n\n")
		loginCMD.PrintDefaults()
	}

	flags := []*flag.FlagSet{deployCMD, destroyCMD, loginCMD}

	if len(os.Args) < 2 {
		printUsage(flags)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "deploy":
		deployCMD.Parse(os.Args[2:])
		if *ports == "" {
			fmt.Fprintf(os.Stderr, "fatal: ports not provided")
			os.Exit(1)
		}
		err := Deploy(*ports, showLogs)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	case "destroy":
		destroyCMD.Parse(os.Args[2:])
		err := Destroy(showLogs)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	case "login":
		loginCMD.Parse(os.Args[2:])
		if *mnemonics == "" {
			fmt.Fprintf(os.Stderr, "fatal: mnemonics not provided")
			os.Exit(1)
		}
		if *network == "" {
			fmt.Fprintf(os.Stderr, "fatal: network not provided")
			os.Exit(1)
		}
		err := Login(*mnemonics, *network, showLogs)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	default:
		printUsage(flags)
		os.Exit(1)
	}
}

func printUsage(flags []*flag.FlagSet) {
	fmt.Fprintf(os.Stderr, "Usage: gridify COMMAND [options]\n\n")
	fmt.Fprintf(os.Stderr, "Commands:\n")
	for _, flag := range flags {
		fmt.Fprintf(os.Stderr, "\t%s\n", flag.Name())
	}
	fmt.Fprintln(os.Stderr)
}
