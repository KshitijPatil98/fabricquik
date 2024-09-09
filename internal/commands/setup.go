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

func Setup() {
	setupCmd := flag.NewFlagSet("setup", flag.ExitOnError)

	fabricVersion := setupCmd.String("fabric", "", "[REQUIRED] The version of fabric binaries and docker images of all components execept ca and couchdb")
	caVersion := setupCmd.String("ca", "", "[REQUIRED] The version of fabric-ca binaries and docker images")
	path := setupCmd.String("path", "", "[REQUIRED] The absolute path of directory where all the binaries will be downloaded.")
	networkType := setupCmd.String("type", "", "[REQUIRED] The type of network you want to setup. Available options are as follows\n 1) basic: Single channel with multiple organisations.")
	channelName := setupCmd.String("channel", "mychannel", "[REQUIRED] The name of the channel. Set the flag only if the type of network is basic or pdc. Default value will be used if flag not set in the case of basic and pdc type of network")
	chaincodeName := setupCmd.String("chaincode", "", "[REQUIRED] The name of the chaincode folder. Make sure this folder is inside the directory specified by path")
	networkName := setupCmd.String("network", "", "[REQUIRED] The name of the network. This will be used as the domain name and network name within our network")
	docker := setupCmd.Bool("docker", false, "[OPTIONAL] If set docker images will also be downloaded along with binaries.")
	couchVersion := setupCmd.String("couch", "", "[OPTIONAL] The version of couchdb images.")

	if len(os.Args) < 3 {
		fmt.Println("Error: No flags supplied")
		fmt.Println(`Run "fabriquik setup --help" for usage`)
		return

	}
	err := setupCmd.Parse(os.Args[2:])
	if err != nil {
		fmt.Printf("The following error occured while parsing: %v", err)
		return
	}

	if *fabricVersion == "" || *caVersion == "" || *path == "" || *networkType == "" || *networkName == "" || *channelName == "" || *chaincodeName == "" {
		fmt.Println(`Error: Either of the compulsory flags --fabric, --ca, --channel, --chaincode, --network, --path and --type are not set.`)
		fmt.Println("Please set all the compulsory flags and try again")
		fmt.Println("Example usage : fabriquik setup --fabri 2.5.4 --ca 1.5.7 --channel mychannel --network fabriquik  --type basic --path /home/ubuntu/Hyperledger")
		return
	}
	*networkType = strings.ToLower(*networkType)
	*networkName = strings.ToLower(*networkName)
	*channelName = strings.ToLower(*channelName)
	*chaincodeName = strings.ToLower(*chaincodeName)

	chaincodePath := filepath.Join(*path, *chaincodeName)
	chaincodePkgPath := filepath.Join(*path, (*chaincodeName)+".tar.gz")

	networkConfig := models.Network{
		NetworkDirectory: *path,
		NetworkType:      *networkType,
		ChannelName:      *channelName,
		NetworkName:      *networkName,
		ChaincodeName:    *chaincodeName,
		ChaincodePath:    chaincodePath,
		ChaincodePkgPath: chaincodePkgPath,
	}

	if networkConfig.NetworkType != "basic" {
		fmt.Println("The network type which is supported now is basic. Please select an appropriate option and try again")
		return

	}
	chaincodePresent, err := utils.CheckSetupFolders(&networkConfig)
	if err != nil {
		return
	}

	if !chaincodePresent {
		fmt.Printf("No chaincode folder present. Please make sure a folder with name %v is present at path %v and try again.", *chaincodeName, *path)
		return
	}
	fmt.Println("This might take few minutes depending on your internet connection")

	if *docker {
		//log.Println("Docker flag is set. Proceeding to download the binaries and images")
		cmdBin := fmt.Sprintf("curl -sSL https://bit.ly/2ysbOFE | bash -s -- %v %v -d -s", fabricVersion, caVersion)
		cmdDocker := fmt.Sprintf("curl -sSL https://bit.ly/2ysbOFE | bash -s -- %v %v -s -b", fabricVersion, caVersion)
		cmdDelete := "rm -rf *"
		command := exec.Command("bash", "-c", cmdBin)
		command.Dir = *path
		// Run the command and check for errors
		err := command.Run()
		if err != nil {
			fmt.Printf("The following error was encountered while downloading binaries : %v", err)
			return
		}
		command = exec.Command("bash", "-c", cmdDocker)
		err = command.Run()
		if err != nil {
			fmt.Printf("The following error was encountered while downloading images : %v", err)
			command = exec.Command("bash", "-c", cmdDelete)
			err = command.Run()
			if err != nil {
				fmt.Printf("The following error is encountered while deleting the binaries : %v.Please manually delete the binaries and try again", err)
				return
			}
			fmt.Println("The downloaded binaries have been cleaned as a result of incomplete execution. Please rectify the issue and run the command again.")
			return
		}
		fmt.Println("\nThe binaries and images have been successfully downloaded in the mentioned folder.")
	} else {
		//log.Println("Docker flag not set. Proceeding to download just the binaries")
		cmdBin := fmt.Sprintf("curl -sSL https://bit.ly/2ysbOFE | bash -s -- %v %v -d -s", *fabricVersion, *caVersion)
		command := exec.Command("bash", "-c", cmdBin)
		command.Dir = *path
		err := command.Run()
		if err != nil {
			fmt.Printf("The following error was encountered while downloading binaries : %v\n", err)
			return
		}
		fmt.Println("The binaries have been successfully downloaded in the mentioned folder.")
	}

	if *couchVersion != "" {
		//log.Println("Downloading docker images for couchdb")
		cmdCouch := fmt.Sprintf("docker pull couchdb:%v", *couchVersion)
		command := exec.Command("bash", "-c", cmdCouch)
		err := command.Run()
		if err != nil {
			fmt.Printf("The following error was encountered while downloading couchdb images : %v\n", err)
			return
		}
		fmt.Println("The couchdb docker image has been successfully downloaded ")

	}

	fmt.Println("\nPlease go ahead and set the path variable inside bashrc to point correctly at the binaries we just downloaded.")

	err = utils.CreateProjectDirectories(&networkConfig)
	if err != nil {
		return
	}

	templateFilePath := "../../internal/templates/generic/bootstrap.sh"
	outputFilePath := filepath.Join(networkConfig.NetworkDirectory, "network_files/bootstrap.sh")

	data, err := utils.ReadFile(templateFilePath)
	if err != nil {
		return
	}
	err = utils.WriteFile(outputFilePath, data)
	if err != nil {
		return
	}
	err = os.Chmod(outputFilePath, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	templateFilePath = "../../internal/templates/generic/stop.sh"
	outputFilePath = filepath.Join(networkConfig.NetworkDirectory, "network_files/stop.sh")

	data, err = utils.ReadFile(templateFilePath)
	if err != nil {
		return
	}
	err = utils.WriteFile(outputFilePath, data)
	if err != nil {
		return
	}
	// Equivalent to `chmod +x` in shell
	err = os.Chmod(outputFilePath, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}

	//The config file will be created and stored in the /home/.fabriquik folder.
	err = utils.OutputJson(&networkConfig, &config.ConfigFilePath)
	if err != nil {
		return
	}

}
