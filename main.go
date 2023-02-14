package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rawdaGastan/gridify/internal/cmd"
)

func main() {
	var mnemonics string

	deployCMD := flag.NewFlagSet("deploy", flag.ExitOnError)

	deployCMD.StringVar(&mnemonics, "mnemonics", "", "mneomnics for authentication")
	port := deployCMD.String("port", "", "port to forward the FQDN to")

	destroyCMD := flag.NewFlagSet("destroy", flag.ExitOnError)
	destroyCMD.StringVar(&mnemonics, "mnemonics", "", "mneomnics for authentication")

	if len(os.Args) < 2 {
		fmt.Printf("expected '%s' or '%s'\n", deployCMD.Name(), destroyCMD.Name())
		os.Exit(1)
	}

	switch os.Args[1] {
	case "deploy":
		err := deployCMD.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if mnemonics == "" {
			fmt.Println("mnemonics not provided")
			os.Exit(1)
		}
		if *port == "" {
			fmt.Println("port not provided")
			os.Exit(1)
		}
		err = cmd.HandleDeploy(mnemonics, *port)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "destroy":
		err := destroyCMD.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if mnemonics == "" {
			fmt.Println("mnemonics not provided")
			os.Exit(1)
		}
		err = cmd.HandleDestroy(mnemonics)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		fmt.Printf("expected '%s' or '%s'\n", deployCMD.Name(), destroyCMD.Name())
		os.Exit(1)
	}

}
