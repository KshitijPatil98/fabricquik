package commands

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/KshitijPatil98/fabriquik/internal/config"
	"github.com/KshitijPatil98/fabriquik/internal/models"
	"github.com/KshitijPatil98/fabriquik/internal/utils"
)

func Channel() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: fabriquik channel [create|join]")
		return
	}

	switch os.Args[2] {
	case "create":
		// Handle `fabricquik channel create --owner incalus`
		createChannel()

	case "join":
		// Handle `fabricquik channel join --org incalus`
		joinChannel()

	default:
		fmt.Println("Unknown subcommand for 'channel':", os.Args[2])
		fmt.Println("Example usage : fabriquik channel create --owner incalus\n\t\tfabriquik channel join --org incalus")

		return
	}
}

func createChannel() {

	createChannelCmd := flag.NewFlagSet("create", flag.ExitOnError)

	ownerName := createChannelCmd.String("owner", "", " [REQUIRED] The name of the owner org.")

	if len(os.Args) < 4 {
		fmt.Println("Error: No flags supplied")
		fmt.Println(`Example usage "fabriquik channel create --owner incalus" for usage`)
		return

	}
	err := createChannelCmd.Parse(os.Args[3:])
	if err != nil {
		fmt.Printf("Error occured while parsing: %v", err)
		return
	}

	if *ownerName == "" {
		fmt.Println(`Error: --owner flag is not set.`)
		fmt.Println("Please set the flag and try again")
		fmt.Println("Example usage : fabriquik channel create --owner incalus")
		return
	}

	var networkConfig models.Network

	err = utils.ReadJson(&config.ConfigFilePath, &networkConfig)
	if err != nil {
		return
	}
	orgNameMap := networkConfig.Orgs
	channelName := networkConfig.ChannelName
	configFilePath := orgNameMap[*ownerName]
	if configFilePath == "" {
		fmt.Println("The owner org does not exist. Please make sure the org is correctly onboarded")
		return
	}
	var orgConfig models.Org_Config
	err = utils.ReadJson(&configFilePath, &orgConfig)
	if err != nil {
		return
	}
	orgType := orgConfig.OrgType
	if orgType != "owner" {
		fmt.Println("The org you specified using the owner flag is not the owner. Please make sure the supplied org is owner")
		return
	}

	createChannelFilePath := networkConfig.NetworkDirectory + fmt.Sprintf(`/network_files/script_files/channel/%v/%v/create_genesis_block.sh`, *ownerName, channelName)

	_, err = os.Stat(createChannelFilePath)
	if err != nil {
		fmt.Println("The create_genesis_block.sh file is missing. The file is generated duing the bootstrap of the org. Please make sure the file is present and try again.")
		return
	}
	cmd := exec.Command("bash", createChannelFilePath)
	cmd.Dir = filepath.Join(networkConfig.NetworkDirectory, "network_files")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// // Run the command and capture output or errors
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Command execution failed with error: %v", err)
	}

}

func joinChannel() {

	joinChannelCmd := flag.NewFlagSet("join", flag.ExitOnError)

	orgName := joinChannelCmd.String("org", "", " [REQUIRED] The name of the org.")

	if len(os.Args) < 4 {
		fmt.Println("Error: No flags supplied")
		fmt.Println(`Example usage "fabriquik channel join --org incalus" for usage`)
		return

	}
	err := joinChannelCmd.Parse(os.Args[3:])
	if err != nil {
		fmt.Printf("Error occured while parsing: %v", err)
		return
	}
	if *orgName == "" {
		fmt.Println(`Error: --org flag is not set.`)
		fmt.Println("Please set the flag and try again")
		fmt.Println("Example usage : fabriquik channel join --help")
		return
	}

	var networkConfig models.Network

	err = utils.ReadJson(&config.ConfigFilePath, &networkConfig)
	if err != nil {
		return
	}
	orgNameMap := networkConfig.Orgs
	channelName := networkConfig.ChannelName

	configFilePath := orgNameMap[*orgName]
	if configFilePath == "" {
		fmt.Println("The org which you mentioned does not exist. Please make sure the org is correctly onboarded")
		return
	}
	var orgConfig models.Org_Config
	err = utils.ReadJson(&configFilePath, &orgConfig)
	if err != nil {
		return
	}
	orgType := orgConfig.OrgType

	var joinChannelFilePath string

	if orgType == "owner" || orgType == "org" {
		joinChannelFilePath = networkConfig.NetworkDirectory + fmt.Sprintf(`/network_files/script_files/channel/%v/%v/join_peers.sh`, *orgName, channelName)
		anchorPeerFilePath := networkConfig.NetworkDirectory + fmt.Sprintf(`/network_files/script_files/channel/%v/%v/setAnchorPeer.sh`, *orgName, channelName)

		err := os.Chmod(anchorPeerFilePath, 0755)
		if err != nil {
			fmt.Println("Error setting permissions:", err)
			return
		}

	} else {
		joinChannelFilePath = networkConfig.NetworkDirectory + fmt.Sprintf(`/network_files/script_files/channel/orderer/%v/osn_orderer_join.sh`, channelName)

	}

	_, err = os.Stat(joinChannelFilePath)
	if err != nil {
		fmt.Println("The join_peers.sh file is missing. The file is generated duing the bootstrap of the org. Please make sure the file is present and try again.")
		return
	}
	cmd := exec.Command("bash", joinChannelFilePath)
	cmd.Dir = filepath.Join(networkConfig.NetworkDirectory, "network_files")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// // Run the command and capture output or errors
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Command execution failed with error: %v", err)
	}

}
