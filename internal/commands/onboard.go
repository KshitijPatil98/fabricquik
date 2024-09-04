package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/KshitijPatil98/fabriquik/internal/config"
	"github.com/KshitijPatil98/fabriquik/internal/models"
	"github.com/KshitijPatil98/fabriquik/internal/utils"
)

func Onboard() {
	onboardCmd := flag.NewFlagSet("onboard", flag.ExitOnError)

	configFileName := onboardCmd.String("config", "", "[REQUIRED] The name of the json file which consist of all the configiuration for the organisation. This configuration file has to be in the root network directory. The root network directory is the directory where all the binaries were downloaded. It was specified using --path flag during setup. If generateconfig flag is set, then a sample configuration with name as specified using --config flag  will be created.")
	//The generateOrgType will only be set and should only be used when the generateconfig flag is set.
	generateOrgType := onboardCmd.String("type", "", "[OPTIONAL] The type of organisation. Available options are\n 1.Owner 2.Org 3.Orderer.\n This flag is required if the generateConfig flag is set")

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
	if *configFileName == "" {
		fmt.Println(`Error: --config is not set.`)
		fmt.Println("Please set the flag and try again")
		fmt.Println("Example usage :   fabriquik onboard --config=org1_config.json")
		fmt.Println("              OR                                ")
		fmt.Println("                  fabriquik onboard --config=org1_config.json --type owner--generateconfig")
		fmt.Println("\nRun the second command if you dont have a configuration file. The second command will generate the configuration file at your root network folder with the name specified using config flag. You can then edit the json file as per your network and then run the first command")
		return
	}

	if *generateConfig && *generateOrgType == "" {
		fmt.Println(`Error: --type is not set.`)
		fmt.Println("Please set the flag and try again")
		fmt.Println("Example usage : fabriquik onboard --type owner --generateconfig")
		return
	}

	if *generateConfig && *generateOrgType != "" {
		*generateOrgType = strings.ToLower(*generateOrgType)
		if *generateOrgType != "owner" && *generateOrgType != "org" && *generateOrgType != "orderer" {
			fmt.Println("Invalid option for type. The available options for type are: 1)Owner 2)Org 3)Orderer")
			return

		}
	}

	var networkConfig models.Network

	err = utils.ReadJson(&config.ConfigFilePath, &networkConfig)
	if err != nil {
		return
	}

	orgConfigFilePath := filepath.Join(networkConfig.NetworkDirectory, *configFileName)

	if *generateConfig && *generateOrgType == "owner" {
		config := models.Org_Config{
			OrgName: "incalus",
			OrgType: "owner",
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
		err := utils.OutputJson(&config, &orgConfigFilePath)
		if err != nil {
			return
		}
		return
	}

	if *generateConfig && *generateOrgType == "org" {
		config := models.Org_Config{
			OrgName: "bolthouse",
			OrgType: "org",
			Ca: models.Ca{
				TlscaPort: "7054",
				OrgcaPort: "8054",
			},
			Peers: []models.Peer{
				{
					PeerLis: "9054",
					PeerCc:  "1054",
					PeerOp:  "9444",
					Couchdb: "6984",
				},
				{
					PeerLis: "9064",
					PeerCc:  "1064",
					PeerOp:  "9454",
					Couchdb: "6994",
				},
			},
			Orderer: models.Peer_Orderer{
				OrdererName: "orderer1",
				OrdererLis:  "9057",
			},
		}
		err := utils.OutputJson(&config, &orgConfigFilePath)
		if err != nil {
			return
		}
		return
	}

	if *generateConfig && *generateOrgType == "orderer" {
		config := models.Orderer_Config{
			OrgType: "orderer",
			Ca: models.Ca{
				TlscaPort: "7056",
				OrgcaPort: "8056",
			},
			Orderers: []models.Orderer{
				{
					OrdererLis:   "9056",
					OrdererAdm:   "9046",
					OrdererAdmOp: "9446",
				},
				{
					OrdererLis:   "9057",
					OrdererAdm:   "9047",
					OrdererAdmOp: "9447",
				},
				{
					OrdererLis:   "9058",
					OrdererAdm:   "9048",
					OrdererAdmOp: "9448",
				},
			},
		}
		err := utils.OutputJson(&config, &orgConfigFilePath)
		if err != nil {
			return
		}

		return
	}

	var present bool
	present, err = utils.CheckOnboardFolders(&networkConfig.NetworkDirectory, configFileName)
	if err != nil {
		return
	}
	if !present {
		return
	}

	var orgConfig models.Org_Config
	var ordererConfig models.Orderer_Config
	var networkOrgConfig models.Network_Org_Config
	var networkOrdererConfig models.Network_Orderer_Config

	orgType, err := utils.FindOrgType(&orgConfigFilePath)
	if err != nil {
		return
	}

	if orgType != "owner" && orgType != "org" && orgType != "orderer" {
		fmt.Println("The OrgType value is invalid. Available options are 1.Owner 2.Org 3.Orderer")
		return
	}
	if orgType == "owner" || orgType == "org" {
		err = utils.ReadJson(&orgConfigFilePath, &orgConfig)
		if err != nil {
			return
		}
		orgConfig.OrgName = strings.ToLower(orgConfig.OrgName)
	}

	if orgType == "orderer" {
		err = utils.ReadJson(&orgConfigFilePath, &ordererConfig)
		if err != nil {
			return
		}

		ordererConfig.OrgName = "orderer"

	}

	if orgType == "owner" || orgType == "org" {
		networkOrgConfig = models.Network_Org_Config{
			Config:  orgConfig,
			Network: networkConfig,
		}
		err = utils.CreateOrgDirectories(&networkOrgConfig)
		if err != nil {
			return
		}

		err = utils.CreateOrgFiles(&networkOrgConfig)
		if err != nil {
			return
		}

		if networkConfig.Orgs == nil {
			networkConfig.Orgs = make(map[string]string)
		}
		networkConfig.Orgs[orgConfig.OrgName] = orgConfigFilePath
		err = utils.OutputJson(&networkConfig, &config.ConfigFilePath)
		if err != nil {
			return
		}

	}
	if orgType == "orderer" {
		networkOrdererConfig = models.Network_Orderer_Config{
			Config:  ordererConfig,
			Network: networkConfig,
		}
		//fmt.Println("Inside onboard", networkOrdererConfig)
		err = utils.CreateOrgDirectories(&networkOrdererConfig)
		if err != nil {
			return
		}

		err = utils.CreateOrgFiles(&networkOrdererConfig)
		if err != nil {
			return
		}
		if networkConfig.Orgs == nil {
			networkConfig.Orgs = make(map[string]string)
		}
		networkConfig.Orgs[ordererConfig.OrgName] = orgConfigFilePath
		err = utils.OutputJson(&networkConfig, &config.ConfigFilePath)
		if err != nil {
			return
		}

	}

}
