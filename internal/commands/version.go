package commands

import "fmt"

const (
	version = "1.0.0"
)

func Version() {
	fmt.Println("Fabriquik is at version", version)
}
