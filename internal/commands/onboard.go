package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/KshitijPatil98/fabriquik/internal/models"
	"github.com/KshitijPatil98/fabriquik/internal/utils"
)

func Onboard() {
	onboardCmd := flag.NewFlagSet("onboard", flag.ExitOnError)

	configFilePath := onboardCmd.String("config", "", "[REQUIRED] The name of the json file which consist of all the configiuration for the organisation. This confoguration file has to be in the root network directory. The root network directory is the directory where all the binaries were downloaded. It was specified using --path flag during setup. If generateconfig flag is set, then a sample configuration with name as specified using --config flag  will be created.")
	orgType := onboardCmd.String("type", "", "[REQUIRED] The type of organisation. Available options are\n 1.Owner 2.Org 3.Orderer")

	generateConfig := onboardCmd.Bool("generateconfig", false, "[OPTIONAL] Generates a sample configuration file with name specified using --config flag")

	if len(os.Args) < 3 {
		fmt.Println("Error: No flags supplied")
		fmt.Println(`Run "fabriquik onboard --help" for usage`)
		return

	}
	err := onboardCmd.Parse(os.Args[2:])
	if err != nil {
		fmt.Printf("The following error occured while parsing: %v", err)
		return
	}
	if *configFilePath == "" {
		fmt.Println(`Error: --config is not set.`)
		fmt.Println("Please set the flag and try again")
		fmt.Println("Example usage :   fabriquik onboard --config=/home/ubuntu/org1_config.json --type owner")
		fmt.Println("              OR                                ")
		fmt.Println("               fabriquik onboard --config=/home/ubuntu/org1_config.json --type owner--generateconfig")
		fmt.Println("\nRun the second command if you dont have a configuration file. The second command will generate the configuration file at the location specified using config flag. You can then edit the json file as per your network and then run the first command")
		return
	}

	if *orgType == "" {
		fmt.Println(`Error: --type is not set.`)
		fmt.Println("Please set the flag and try again")
		fmt.Println("Example usage :   fabriquik onboard --config=/home/ubuntu/org1_config.json --type owner")
		fmt.Println("              OR                                ")
		fmt.Println("               fabriquik onboard --config=/home/ubuntu/org1_config.json --type owner--generateconfig")
		fmt.Println("\nRun the second command if you dont have a configuration file. The second command will generate the configuration file at the location specified using config flag. You can then edit the json file as per your network and then run the first command")
		return
	}

	if *orgType != "owner" && *orgType != "org" && *orgType != "orderer" {
		fmt.Println("Invalid option for type.\nThe available options for type are: 1)owner 2)org 3)orderer")
		return

	}

	if *generateConfig && (*orgType == "owner" || *orgType == "org") {
		config := models.Org_Config{
			OrgName: "Incalus",
			Ca: models.Ca{
				TlscaPort: "7055",
				OrgcaPort: "8055",
			},
			Peers: []models.Peer{
				{
					PeerLis: "9055",
					PeerCc:  "1055",
					PeerOp:  "9445",
					Couchdb: "6985",
				},
				{
					PeerLis: "9065",
					PeerCc:  "1065",
					PeerOp:  "9455",
					Couchdb: "6995",
				},
			},
			Orderer: models.Peer_Orderer{
				OrdererName: "orderer0",
				OrdererLis:  "9056",
			},
		}
		err := utils.OutputJson(&config, configFilePath)
		if err != nil {
			return
		}
		return
	}

	if *generateConfig && (*orgType == "orderer") {
		// config := models.Org_Config{
		// 	OrgName: "Incalus",
		// 	Path:    "/home/ubuntu/Hyperledger",
		// 	Ca: models.Ca{
		// 		TlscaPort: "7055",
		// 		OrgcaPort: "8055",
		// 	},
		// 	Peers: []models.Peer{
		// 		{
		// 			PeerLis: "9055",
		// 			PeerCc:  "1055",
		// 			PeerOp:  "9445",
		// 			Couchdb: "6985",
		// 		},
		// 		{
		// 			PeerLis: "9065",
		// 			PeerCc:  "1065",
		// 			PeerOp:  "9455",
		// 			Couchdb: "6995",
		// 		},
		// 	},
		// 	Orderer: models.Peer_Orderer{
		// 		OrdererName: "orderer0",
		// 		OrdererLis:  "9056",
		// 	},
		// }
		// err := utils.OutputJson(&config, configFilePath)
		// if err != nil {
		// 	return
		// }
		fmt.Println("Yet to be developed")
		return
	}

	var configPath models.ConfigPath
	localPath := "../../internal/configurations/configInfo.json"
	err = utils.ReadJson(&localPath, &configPath)
	if err != nil {
		return
	}
	var orgConfig models.Org_Config
	var ordererConfig models.Orderer_Config
	var networkOrgConfig models.Network_Org_Config
	//var networkOrdererConfig models.Network_Orderer_Config

	if *orgType == "owner" || *orgType == "org" {
		err = utils.ReadJson(configFilePath, &orgConfig)
		if err != nil {
			return
		}
		orgConfig.OrgType = *orgType
		orgConfig.OrgName = strings.ToLower(orgConfig.OrgName)
		orgConfig.OrgType = strings.ToLower(orgConfig.OrgType)
		orgConfig.NetworkDirectory = configPath.NetworkDirectory
	}

	if *orgType == "orderer" {
		err = utils.ReadJson(configFilePath, &ordererConfig)
		if err != nil {
			return
		}
		ordererConfig.OrgType = *orgType
		ordererConfig.OrgName = strings.ToLower(ordererConfig.OrgName)
		ordererConfig.OrgType = strings.ToLower(ordererConfig.OrgType)
		ordererConfig.NetworkDirectory = configPath.NetworkDirectory

	}

	var present bool
	present, err = utils.CheckOnboardFolders(&configPath.NetworkDirectory)
	if err != nil {
		return
	}
	if !present {
		fmt.Println("Please make sure that the folders network_files ,config and the file network_config.json is present at the supplied path.\nIf not, run the fabriquik setup command and try again")
		return
	}

	var network models.Network
	err = utils.ReadJson(&configPath.NetworkConfigPath, &network)
	if err != nil {
		return
	}

	if *orgType == "owner" || *orgType == "org" {
		networkOrgConfig = models.Network_Org_Config{
			Config:  orgConfig,
			Network: network,
		}
		// err = utils.CreateOrgDirectories(o)
		// if err != nil {
		// 	return
		// }

		err = utils.CreateOrgFiles(&networkOrgConfig)
		if err != nil {
			return
		}

	}

	if *orgType == "orderer" {
		// networkOrdererConfig = models.Network_Orderer_Config{
		// 	Config:  ordererConfig,
		// 	Network: network,
		// }
	}

}
