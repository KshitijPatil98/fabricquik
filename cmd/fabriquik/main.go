package main

import (
	"fmt"
	"os"

	"github.com/KshitijPatil98/fabriquik/internal/commands"
	"github.com/KshitijPatil98/fabriquik/internal/config"
)

func main() {

	config.Initialize()

	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  fabriquik version")
		fmt.Println("Available Commands")
		fmt.Println("  version  Prints the current version")
		fmt.Println("  setup    Downloads fabric binaries, images and sets up a base project structure")
		fmt.Println("  onboard  Creates directories and files for boostrapping components of an organisation")
		fmt.Println("  configtx Creates a configtx file based on the configuration supplied by user")
		fmt.Println("  configtx Creates a private data collection config file based on the configuration supplied")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "help":
		fmt.Println("Usage:")
		fmt.Println("  fabriquik version")
		fmt.Println("Available Commands")
		fmt.Println("  version  Prints the current version")
		fmt.Println("  setup    Downloads fabric binaries, images and sets up a base project structure")
		fmt.Println("  onboard  Creates directories and files for boostrapping components of an organisation")
		fmt.Println("  configtx Creates a configtx file based on the configuration supplied by user")
		fmt.Println("  configtx Creates a private data collection config file based on the configuration supplied")
	case "version":
		commands.Version()

	case "setup":
		commands.Setup()

	case "onboard":
		commands.Onboard()

	case "bootstrap":
		commands.Bootstrap()

	case "configtx":
		commands.Configtx()

	case "privatedata":
		commands.PrivateData()

	case "stop":
		commands.Stop()

	default:
		fmt.Println("Unknown command", os.Args[1])
		os.Exit(1)

	}

}
