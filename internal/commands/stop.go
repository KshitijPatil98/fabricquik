package commands

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/KshitijPatil98/fabriquik/internal/config"
	"github.com/KshitijPatil98/fabriquik/internal/models"
	"github.com/KshitijPatil98/fabriquik/internal/utils"
)

func Stop() {

	stopCmd := flag.NewFlagSet("stop", flag.ExitOnError)

	orgName := stopCmd.String("org", "", "[REQUIRED] The name of the organisation you want to stop.")

	if len(os.Args) < 3 {
		fmt.Println("Error: No flags supplied")
		fmt.Println(`Run "fabriquik stop --help" for usage`)
		return

	}
	err := stopCmd.Parse(os.Args[2:])
	if err != nil {
		fmt.Printf("Error occured while parsing: %v", err)
		return
	}
	if *orgName == "" {
		fmt.Println(`Error: --org flag is not set.`)
		fmt.Println("Please set the flag and try again")
		fmt.Println("Example usage :   fabriquik stop --org Incalus")
		return
	}
	*orgName = strings.ToLower(*orgName)

	var networkConfig models.Network
	err = utils.ReadJson(&config.ConfigFilePath, &networkConfig)
	if err != nil {
		return
	}
	orgNameMap := networkConfig.Orgs

	if orgNameMap[*orgName] == "" {
		fmt.Println("The specified org does not exist")
		return
	}
	stopFilePath := "./stop.sh"
	stopFileFullPath := filepath.Join(networkConfig.NetworkDirectory, "network_files", "stop.sh")

	_, err = os.Stat(stopFileFullPath)
	if err != nil {
		fmt.Println("The stop.sh file is missing. The file is generated duing the setup. Please make sure the file is present and try again.")
		return
	}
	orgType := ""
	if *orgName == "orderer" {
		orgType = "orderer"
	} else {
		orgConfigPath := orgNameMap[*orgName]
		var orgConfig models.Org_Config

		err = utils.ReadJson(&orgConfigPath, &orgConfig)
		if err != nil {
			return
		}
		orgType = strings.ToLower(orgConfig.OrgType)

	}

	cmd := exec.Command("bash", stopFilePath, *orgName, orgType, networkConfig.ChannelName, networkConfig.NetworkName)
	cmd.Dir = filepath.Join(networkConfig.NetworkDirectory, "network_files")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and capture output or errors
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Command execution failed with error: %v", err)
	}

}
