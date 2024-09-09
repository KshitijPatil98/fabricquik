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

func Bootstrap() {

	bootstrapCmd := flag.NewFlagSet("bootstrap", flag.ExitOnError)

	orgName := bootstrapCmd.String("org", "", "[REQUIRED] The name of the organisation you want to boostrap.")

	if len(os.Args) < 3 {
		fmt.Println("Error: No flags supplied")
		fmt.Println(`Run "fabriquik bootstrap --help" for usage`)
		return

	}
	err := bootstrapCmd.Parse(os.Args[2:])
	if err != nil {
		fmt.Printf("Error occured while parsing: %v", err)
		return
	}
	if *orgName == "" {
		fmt.Println(`Error: --org flag is not set.`)
		fmt.Println("Please set the flag and try again")
		fmt.Println("Example usage :   fabriquik bootstrap --org Incalus")
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
		fmt.Println("The specified org has not been onboarded properly. Please onboard the org properly and try again")
		return
	}
	bootstrapFilePath := "./bootstrap.sh"
	bootstrapFileFullPath := filepath.Join(networkConfig.NetworkDirectory, "network_files", "bootstrap.sh")
	fmt.Println(bootstrapFileFullPath)
	_, err = os.Stat(bootstrapFileFullPath)
	if err != nil {
		fmt.Println("The bootstrap.sh file is missing. The file is generated duing the setup. Please make sure the file is present and try again.")
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

	cmd := exec.Command("bash", bootstrapFilePath, orgType, *orgName)
	cmd.Dir = filepath.Join(networkConfig.NetworkDirectory, "network_files")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command and capture output or errors
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Command execution failed with error: %v", err)
	}

}
