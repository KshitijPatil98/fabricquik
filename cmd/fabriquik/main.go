package main

import (
	"fmt"
	"os"

	"github.com/KshitijPatil98/fabriquik/internal/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  fabriquick version")
		fmt.Println("Available Commands")
		fmt.Println("  version  Prints the current version")
		fmt.Println("  setup    Downloads fabric binaries, images and sets up a base project structure")
		fmt.Println("  onboard  Creates directories and files for boostrapping components of an organisation")

		os.Exit(1)
	}

	switch os.Args[1] {
	case "version":
		commands.Version()

	case "setup":
		commands.Setup()

	case "onboard":
		commands.Onboard()

	default:
		fmt.Println("Unknown command", os.Args[1])
		os.Exit(1)

	}

}
