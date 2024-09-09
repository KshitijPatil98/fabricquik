package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/KshitijPatil98/fabriquik/internal/config"
	"github.com/KshitijPatil98/fabriquik/internal/models"
)

func CheckSetupFolders(networkConfig *models.Network) (bool, error) {

	command := exec.Command("bash", "-c", "ls")
	command.Dir = networkConfig.NetworkDirectory
	// Run the command and check for errors
	outputByte, err := command.CombinedOutput()

	if err != nil {
		fmt.Printf("Following error occured while runnung the ls command: %v", err)
		return false, err
	}
	outputStr := string(outputByte)

	folderNames := strings.Split(outputStr, "\n")
	chaincodePresent := false
	for _, folderName := range folderNames {

		if folderName == networkConfig.ChaincodeName {
			chaincodePresent = true
		}

	}

	return chaincodePresent, nil
}

func CheckOnboardFolders(path, fileName *string) (bool, error) {

	command := exec.Command("bash", "-c", "ls")
	command.Dir = *path
	// Run the command and check for errors
	outputByte, err := command.CombinedOutput()

	if err != nil {
		fmt.Printf("Following error occured while runnung the ls command: %v", err)
		return false, err
	}
	outputStr := string(outputByte)

	folderNames := strings.Split(outputStr, "\n")
	configFolderPresent := false
	networkFilesFolderPresent := false
	orgConfigJsonPresent := false
	configFilePresent := false
	for _, folderName := range folderNames {

		if folderName == "config" {
			configFolderPresent = true
		}
		if folderName == "network_files" {
			networkFilesFolderPresent = true
		}

		if folderName == *fileName {
			orgConfigJsonPresent = true
		}

	}
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting the home directory of the user")
		return false, err
	}
	configFileFolder := filepath.Join(homedir, ".fabriquik")
	command = exec.Command("bash", "-c", "ls")
	command.Dir = configFileFolder

	// Run the command and check for errors
	outputByte, err = command.CombinedOutput()

	if err != nil {
		fmt.Printf("Following error occured while runnung the ls command: %v", err)
		fmt.Println()
		return false, err
	}
	outputStr = string(outputByte)

	folderNames = strings.Split(outputStr, "\n")

	for _, fileName := range folderNames {

		if fileName == "config.json" {
			configFilePresent = true
		}

	}
	if configFolderPresent && networkFilesFolderPresent && configFilePresent && orgConfigJsonPresent {
		return true, nil
	}
	if !configFolderPresent {
		fmt.Println(`The config folder is missing in your network root directory. This folder gets created when you run the "fabriquik setup" command. Please run "fabriquik setup --help" and try again`)
	}
	if !networkFilesFolderPresent {
		fmt.Println(`The network_files folder is missing in your network root directory. This folder gets created when you run the "fabriquik setup" command. Please run "fabriquik setup --help" and try again`)
	}
	if !orgConfigJsonPresent {
		fmt.Printf(`The configuration file %v does not exist in your network root directory. You can generate a sample file using "fabriquik onboard". Please run "fabriquik onboard --help" and try again`, *fileName)
		fmt.Println()
	}

	if !configFilePresent {
		fmt.Printf(`The config.json file is not present at %v. This file is created when you run "fabriquik setup" command and stored at %v. Please run "fabriquik setup --help" and try again`, config.ConfigFilePath, config.ConfigFilePath)
		fmt.Println()
	}

	return false, nil
}

func CreateProjectDirectories(networkConfig *models.Network) error {

	if networkConfig.NetworkType == "basic" {
		directories := []string{
			"network_files/organizations/ordererOrganizations",
			"network_files/organizations/peerOrganizations",
			"network_files/organizations/orgca_certs",
			"network_files/organizations/tlsca_certs",
			fmt.Sprintf("network_files/channel_files/channel_artifacts/%v", networkConfig.ChannelName),
			fmt.Sprintf("network_files/channel_files/configtx_files/%v", networkConfig.ChannelName),
			fmt.Sprintf("network_files/privatedata_files/%v", networkConfig.ChannelName),
			"network_files/compose_files/peer_couch_orderer",
			"network_files/compose_files/orgca",
			"network_files/compose_files/tlsca",
			"network_files/connection_files",
			"network_files/script_files/ccp",
			"network_files/script_files/channel",
			"network_files/script_files/orgcerts",
			"network_files/script_files/tlscerts",
		}
		for _, dir := range directories {
			fullPath := filepath.Join(networkConfig.NetworkDirectory, dir)
			err := os.MkdirAll(fullPath, os.ModePerm)
			if err != nil {
				fmt.Printf("Error occured while creating directory %s : %v", fullPath, err)
				fmt.Printf("The partially created directories will now be deleted")
				err := os.RemoveAll(networkConfig.NetworkDirectory)
				if err != nil {
					fmt.Printf("Error occured while cleaning partially created directories : %v. Please manually clean the directores", err)
					return err
				}
				fmt.Printf("Partially created directories are cleaned. Please rectify the error and try again")
				return err
			}

		}
		fmt.Println("\nAll the directories have been created successfully. You can go ahead and onboard the owner.")
		return nil

	}
	return nil

}

func CreateOrgDirectories(gen_config interface{}) error {

	var networkOrgConfig *models.Network_Org_Config
	var networkOrdererConfig *models.Network_Orderer_Config

	inputType := reflect.TypeOf(gen_config)
	typeNetworkOrgConfig := reflect.TypeOf(networkOrgConfig)
	typeNetworkOrdererConfig := reflect.TypeOf(networkOrdererConfig)

	var directories []string
	//We dont use *networkOgConfig(ie we dont dereference) because in either cases one of the variables will be nil and deferencing a nil pointer will crash.
	if inputType == typeNetworkOrgConfig {
		networkOrgConfig, ok := gen_config.(*models.Network_Org_Config)
		if !ok {
			fmt.Println("Type assertion to models.Network_Org_Config failed while creating org directories")
			return errors.New("type assertion to models.Network_Org_Config failed")
		}
		directories = []string{
			fmt.Sprintf("network_files/organizations/tlsca_certs/%v", networkOrgConfig.Config.OrgName),
			fmt.Sprintf("network_files/organizations/tlsca_certs/%v/users/tlsca_admin", networkOrgConfig.Config.OrgName),
			fmt.Sprintf("network_files/organizations/orgca_certs/%v", networkOrgConfig.Config.OrgName),
			fmt.Sprintf("network_files/organizations/orgca_certs/%v/tls", networkOrgConfig.Config.OrgName),

			fmt.Sprintf("network_files/compose_files/tlsca/%v", networkOrgConfig.Config.OrgName),
			fmt.Sprintf("network_files/compose_files/orgca/%v", networkOrgConfig.Config.OrgName),
			fmt.Sprintf("network_files/compose_files/peer_couch_orderer/%v", networkOrgConfig.Config.OrgName),

			fmt.Sprintf("network_files/script_files/tlscerts/%v", networkOrgConfig.Config.OrgName),
			fmt.Sprintf("network_files/script_files/orgcerts/%v", networkOrgConfig.Config.OrgName),
			fmt.Sprintf("network_files/script_files/channel/%v/%v", networkOrgConfig.Config.OrgName, networkOrgConfig.Network.ChannelName),
			fmt.Sprintf("network_files/script_files/ccp/%v", networkOrgConfig.Config.OrgName),
			fmt.Sprintf("network_files/connection_files/%v", networkOrgConfig.Config.OrgName),
		}

		for _, dir := range directories {
			fullPath := filepath.Join(networkOrgConfig.Network.NetworkDirectory, dir)
			err := os.MkdirAll(fullPath, os.ModePerm)
			if err != nil {
				fmt.Printf("Error occured while creating directory %s : %v", fullPath, err)
				fmt.Printf("The partially created directories will now be deleted")
				err := os.RemoveAll(networkOrgConfig.Network.NetworkDirectory)
				if err != nil {
					fmt.Printf("Error occured while cleaning partially created directories : %v. Please manually clean the directores", err)
					return err
				}
				fmt.Printf("Partially created directories are cleaned. Please rectify the error and try again")
				return err
			}

		}
		fmt.Print("All the org directories have been successfully created.")
		return nil

	}
	//We dont use *networkOgConfig(ie we dont dereference) because in either cases one of the variables will be nil and deferencing a nil pointer will crash.
	if inputType == typeNetworkOrdererConfig {
		networkOrdererConfig, ok := gen_config.(*models.Network_Orderer_Config)
		if !ok {
			fmt.Println("Type assertion to models.Network_Orderer_Config failed while creating orderer directories")
			return errors.New("type assertion to models.Network_Orderer_Config failed")
		}

		directories = []string{
			fmt.Sprintf("network_files/organizations/tlsca_certs/%v", networkOrdererConfig.Config.OrgName),
			fmt.Sprintf("network_files/organizations/tlsca_certs/%v/users/tlsca_admin", networkOrdererConfig.Config.OrgName),
			fmt.Sprintf("network_files/organizations/orgca_certs/%v", networkOrdererConfig.Config.OrgName),
			fmt.Sprintf("network_files/organizations/orgca_certs/%v/tls", networkOrdererConfig.Config.OrgName),

			fmt.Sprintf("network_files/compose_files/tlsca/%v", networkOrdererConfig.Config.OrgName),
			fmt.Sprintf("network_files/compose_files/orgca/%v", networkOrdererConfig.Config.OrgName),
			fmt.Sprintf("network_files/compose_files/peer_couch_orderer/%v", networkOrdererConfig.Config.OrgName),

			fmt.Sprintf("network_files/script_files/tlscerts/%v", networkOrdererConfig.Config.OrgName),
			fmt.Sprintf("network_files/script_files/orgcerts/%v", networkOrdererConfig.Config.OrgName),
			fmt.Sprintf("network_files/script_files/channel/%v/%v", networkOrdererConfig.Config.OrgName, networkOrdererConfig.Network.ChannelName),
		}
		for _, dir := range directories {
			fullPath := filepath.Join(networkOrdererConfig.Network.NetworkDirectory, dir)
			err := os.MkdirAll(fullPath, os.ModePerm)
			if err != nil {
				fmt.Printf("Error occured while creating directory %s : %v", fullPath, err)
				fmt.Printf("The partially created directories will now be deleted")
				err := os.RemoveAll(networkOrdererConfig.Network.NetworkDirectory)
				if err != nil {
					fmt.Printf("Error occured while cleaning partially created directories : %v. Please manually clean the directores", err)
					return err
				}
				fmt.Printf("Partially created directories are cleaned. Please rectify the error and try again")
				return err
			}

		}
		fmt.Print("All the orderer directories have been successfully created.")
		return nil
	}

	return errors.New("INVALID STRUCT PASSED")
}

func CreateOrgFiles(gen_config interface{}) error {

	var networkOrgConfig *models.Network_Org_Config
	var networkOrdererConfig *models.Network_Orderer_Config
	inputType := reflect.TypeOf(gen_config)
	typeNetworkOrgConfig := reflect.TypeOf(networkOrgConfig)
	typeNetworkOrdererConfig := reflect.TypeOf(networkOrdererConfig)

	if inputType == typeNetworkOrgConfig {
		networkOrgConfig, ok := gen_config.(*models.Network_Org_Config)
		if !ok {
			fmt.Println("Type assertion to models.Network_Org_Config failed")
			return errors.New("type assertion to models.Network_Org_Config failed")
		}
		templateFilePath := "../../internal/templates/org/tlsca/fabric-ca-server-config.yaml"
		outputFilePath := filepath.Join(networkOrgConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/organizations/tlsca_certs/%v/fabric-ca-server-config.yaml", networkOrgConfig.Config.OrgName))

		var replacements = map[string]string{
			"${ORGNAME}":   networkOrgConfig.Config.OrgName,
			"${TLSCAPORT}": networkOrgConfig.Config.Ca.TlscaPort,
		}
		err := ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/org/tlsca/tlsca_compose.yaml"
		outputFilePath = filepath.Join(networkOrgConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/compose_files/tlsca/%v/tlsca.yaml", networkOrgConfig.Config.OrgName))

		replacements = map[string]string{
			"${ORGNAME}":   networkOrgConfig.Config.OrgName,
			"${TLSCAPORT}": networkOrgConfig.Config.Ca.TlscaPort,
			"${NETWORK}":   networkOrgConfig.Network.NetworkName,
		}
		err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/org/tlsca/tlscerts.sh"
		outputFilePath = filepath.Join(networkOrgConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/script_files/tlscerts/%v/tlscerts.sh", networkOrgConfig.Config.OrgName))

		replacements = map[string]string{
			"${ORGNAME}":   networkOrgConfig.Config.OrgName,
			"${TLSCAPORT}": networkOrgConfig.Config.Ca.TlscaPort,
			"${NETWORK}":   networkOrgConfig.Network.NetworkName,
		}
		err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/org/orgca/fabric-ca-server-config.yaml"
		outputFilePath = filepath.Join(networkOrgConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/organizations/orgca_certs/%v/fabric-ca-server-config.yaml", networkOrgConfig.Config.OrgName))

		replacements = map[string]string{
			"${ORGNAME}":   networkOrgConfig.Config.OrgName,
			"${ORGCAPORT}": networkOrgConfig.Config.Ca.OrgcaPort,
		}
		err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/org/orgca/orgca_compose.yaml"
		outputFilePath = filepath.Join(networkOrgConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/compose_files/orgca/%v/orgca.yaml", networkOrgConfig.Config.OrgName))

		replacements = map[string]string{
			"${ORGNAME}":   networkOrgConfig.Config.OrgName,
			"${ORGCAPORT}": networkOrgConfig.Config.Ca.OrgcaPort,
			"${NETWORK}":   networkOrgConfig.Network.NetworkName,
		}
		err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/org/orgca/orgcerts.sh"
		outputFilePath = filepath.Join(networkOrgConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/script_files/orgcerts/%v/orgcerts.sh", networkOrgConfig.Config.OrgName))

		replacements = map[string]string{
			"${ORGNAME}":   networkOrgConfig.Config.OrgName,
			"${ORGCAPORT}": networkOrgConfig.Config.Ca.OrgcaPort,
			"${NETWORK}":   networkOrgConfig.Network.NetworkName,
		}
		err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/org/peer_couch/peer.yaml"
		outputFilePath = filepath.Join(networkOrgConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/compose_files/peer_couch_orderer/%v/peer.yaml", networkOrgConfig.Config.OrgName))

		replacements = map[string]string{
			"${ORGNAME}":   networkOrgConfig.Config.OrgName,
			"${ORGNAMEFC}": strings.ToUpper(string(networkOrgConfig.Config.OrgName[0])) + networkOrgConfig.Config.OrgName[1:],
			"${PEER0LIS}":  networkOrgConfig.Config.Peers[0].PeerLis,
			"${PEER0CC}":   networkOrgConfig.Config.Peers[0].PeerCc,
			"${PEER0OP}":   networkOrgConfig.Config.Peers[0].PeerOp,
			"${PEER1LIS}":  networkOrgConfig.Config.Peers[1].PeerLis,
			"${PEER1CC}":   networkOrgConfig.Config.Peers[1].PeerCc,
			"${PEER1OP}":   networkOrgConfig.Config.Peers[1].PeerOp,
			"${NETWORK}":   networkOrgConfig.Network.NetworkName,
		}
		err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/org/peer_couch/couch.yaml"
		outputFilePath = filepath.Join(networkOrgConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/compose_files/peer_couch_orderer/%v/couch.yaml", networkOrgConfig.Config.OrgName))

		replacements = map[string]string{
			"${ORGNAME}":      networkOrgConfig.Config.OrgName,
			"${NETWORK}":      networkOrgConfig.Network.NetworkName,
			"${COUCHDBPEER0}": networkOrgConfig.Config.Peers[0].Couchdb,
			"${COUCHDBPEER1}": networkOrgConfig.Config.Peers[1].Couchdb,
		}
		err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/org/peer_couch/cli.yaml"
		outputFilePath = filepath.Join(networkOrgConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/compose_files/peer_couch_orderer/%v/cli.yaml", networkOrgConfig.Config.OrgName))

		replacements = map[string]string{
			"${ORGNAME}": networkOrgConfig.Config.OrgName,
			"${NETWORK}": networkOrgConfig.Network.NetworkName,
			"${CCNAME}":  networkOrgConfig.Network.ChaincodeName,
			"${CCPATH}":  networkOrgConfig.Network.ChaincodePath,
		}
		err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/org/channel/configUpdate.sh"
		outputFilePath = filepath.Join(networkOrgConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/script_files/channel/%v/configUpdate.sh", networkOrgConfig.Config.OrgName))

		replacements = map[string]string{
			"${ORGNAME}":     networkOrgConfig.Config.OrgName,
			"${NETWORK}":     networkOrgConfig.Network.NetworkName,
			"${ORDERERNAME}": networkOrgConfig.Config.Orderer.OrdererName,
			"${ORDERERPORT}": networkOrgConfig.Config.Orderer.OrdererLis,
		}
		err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/org/channel/envVar.sh"
		outputFilePath = filepath.Join(networkOrgConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/script_files/channel/%v/envVar.sh", networkOrgConfig.Config.OrgName))

		replacements = map[string]string{
			"${ORGNAME}":   networkOrgConfig.Config.OrgName,
			"${NETWORK}":   networkOrgConfig.Network.NetworkName,
			"${ORGNAMEFC}": strings.ToUpper(string(networkOrgConfig.Config.OrgName[0])) + networkOrgConfig.Config.OrgName[1:],
			"${PEER0LIS}":  networkOrgConfig.Config.Peers[0].PeerLis,
			"${PEER1LIS}":  networkOrgConfig.Config.Peers[1].PeerLis,
		}
		err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/org/channel/utils.sh"
		outputFilePath = filepath.Join(networkOrgConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/script_files/channel/%v/utils.sh", networkOrgConfig.Config.OrgName))

		data, err := ReadFile(templateFilePath)
		if err != nil {
			return err
		}
		err = WriteFile(outputFilePath, data)
		if err != nil {
			return err
		}

		if networkOrgConfig.Config.OrgType == "owner" {
			templateFilePath = "../../internal/templates/org/channel/create_genesis_block.sh"
			outputFilePath = filepath.Join(networkOrgConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/script_files/channel/%v/%v/create_genesis_block.sh", networkOrgConfig.Config.OrgName, networkOrgConfig.Network.ChannelName))

			replacements = map[string]string{
				"${CHANNELNAME}": networkOrgConfig.Network.ChannelName,
			}
			err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
			if err != nil {
				return err
			}
		}
		templateFilePath = "../../internal/templates/org/channel/join_peers.sh"
		outputFilePath = filepath.Join(networkOrgConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/script_files/channel/%v/%v/join_peers.sh", networkOrgConfig.Config.OrgName, networkOrgConfig.Network.ChannelName))

		replacements = map[string]string{
			"${ORGNAME}":     networkOrgConfig.Config.OrgName,
			"${CHANNELNAME}": networkOrgConfig.Network.ChannelName,
		}
		err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/org/channel/setAnchorPeer.sh"
		outputFilePath = filepath.Join(networkOrgConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/script_files/channel/%v/%v/setAnchorPeer.sh", networkOrgConfig.Config.OrgName, networkOrgConfig.Network.ChannelName))

		replacements = map[string]string{
			"${ORGNAME}":     networkOrgConfig.Config.OrgName,
			"${NETWORK}":     networkOrgConfig.Network.NetworkName,
			"${ORDERERNAME}": networkOrgConfig.Config.Orderer.OrdererName,
			"${ORDERERPORT}": networkOrgConfig.Config.Orderer.OrdererLis,
			"${PEER0LIS}":    networkOrgConfig.Config.Peers[0].PeerLis,
			"${PEER1LIS}":    networkOrgConfig.Config.Peers[1].PeerLis,
		}
		err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/org/channel/ccutils.sh"
		outputFilePath = filepath.Join(networkOrgConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/script_files/channel/%v/%v/ccutils.sh", networkOrgConfig.Config.OrgName, networkOrgConfig.Network.ChannelName))

		replacements = map[string]string{
			"${ORGNAME}":     networkOrgConfig.Config.OrgName,
			"${ORGNAMEC}":    strings.ToUpper(networkOrgConfig.Config.OrgName),
			"${NETWORK}":     networkOrgConfig.Network.NetworkName,
			"${ORDERERNAME}": networkOrgConfig.Config.Orderer.OrdererName,
			"${ORDERERPORT}": networkOrgConfig.Config.Orderer.OrdererLis,
			"${PEER0LIS}":    networkOrgConfig.Config.Peers[0].PeerLis,
		}
		err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		if networkOrgConfig.Config.OrgType == "owner" {
			templateFilePath = "../../internal/templates/org/channel/pkg_ins_app_com_que.sh"
			outputFilePath = filepath.Join(networkOrgConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/script_files/channel/%v/%v/pkg_ins_app_com_que.sh", networkOrgConfig.Config.OrgName, networkOrgConfig.Network.ChannelName))
		}

		if networkOrgConfig.Config.OrgType == "org" {
			templateFilePath = "../../internal/templates/org/channel/ins_app_que.sh"
			outputFilePath = filepath.Join(networkOrgConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/script_files/channel/%v/%v/ins_app_que.sh", networkOrgConfig.Config.OrgName, networkOrgConfig.Network.ChannelName))
		}
		replacements = map[string]string{
			"${ORGNAME}":     networkOrgConfig.Config.OrgName,
			"${ORGNAMEC}":    strings.ToUpper(networkOrgConfig.Config.OrgName),
			"${NETWORK}":     networkOrgConfig.Network.NetworkName,
			"${CCNAME}":      networkOrgConfig.Network.ChaincodeName,
			"${CCPATH}":      networkOrgConfig.Network.ChaincodePath,
			"${CCPKGPATH}":   networkOrgConfig.Network.ChaincodePkgPath,
			"${CHANNELNAME}": networkOrgConfig.Network.ChannelName,
		}
		err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/org/connection/ccp-generate.sh"
		outputFilePath = filepath.Join(networkOrgConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/script_files/ccp/%v/ccp-generate.sh", networkOrgConfig.Config.OrgName))

		replacements = map[string]string{
			"${ORGNAME}":   networkOrgConfig.Config.OrgName,
			"${ORGNAMEFC}": strings.ToUpper(string(networkOrgConfig.Config.OrgName[0])) + networkOrgConfig.Config.OrgName[1:],
			"${NETWORK}":   networkOrgConfig.Network.NetworkName,
			"${PEER0LIS}":  networkOrgConfig.Config.Peers[0].PeerLis,
			"${ORGCAPORT}": networkOrgConfig.Config.Ca.OrgcaPort,
		}
		err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/org/connection/ccp-template.json"
		outputFilePath = filepath.Join(networkOrgConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/script_files/ccp/%v/ccp-template.json", networkOrgConfig.Config.OrgName))

		data, err = ReadFile(templateFilePath)
		if err != nil {
			return err
		}
		err = WriteFile(outputFilePath, data)
		if err != nil {
			return err
		}

		fmt.Println(" All the org files have been successfully created.")
		return nil

	}
	if inputType == typeNetworkOrdererConfig {
		networkOrdererConfig, ok := gen_config.(*models.Network_Orderer_Config)
		if !ok {
			fmt.Println("Type assertion to models.Network_Orderer_Config failed")
			return errors.New("type assertion to models.Network_Orderer_Config failed")
		}
		templateFilePath := "../../internal/templates/orderer/tlsca/fabric-ca-server-config.yaml"
		outputFilePath := filepath.Join(networkOrdererConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/organizations/tlsca_certs/%v/fabric-ca-server-config.yaml", networkOrdererConfig.Config.OrgName))

		var replacements = map[string]string{
			"${ORGNAME}":   networkOrdererConfig.Config.OrgName,
			"${TLSCAPORT}": networkOrdererConfig.Config.Ca.TlscaPort,
		}
		err := ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/orderer/tlsca/tlsca_compose.yaml"
		outputFilePath = filepath.Join(networkOrdererConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/compose_files/tlsca/%v/tlsca.yaml", networkOrdererConfig.Config.OrgName))

		replacements = map[string]string{
			"${ORGNAME}":   networkOrdererConfig.Config.OrgName,
			"${TLSCAPORT}": networkOrdererConfig.Config.Ca.TlscaPort,
			"${NETWORK}":   networkOrdererConfig.Network.NetworkName,
		}
		err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/orderer/tlsca/tlscerts.sh"
		outputFilePath = filepath.Join(networkOrdererConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/script_files/tlscerts/%v/tlscerts.sh", networkOrdererConfig.Config.OrgName))

		replacements = map[string]string{
			"${ORGNAME}":   networkOrdererConfig.Config.OrgName,
			"${TLSCAPORT}": networkOrdererConfig.Config.Ca.TlscaPort,
			"${NETWORK}":   networkOrdererConfig.Network.NetworkName,
		}
		err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/orderer/orgca/fabric-ca-server-config.yaml"
		outputFilePath = filepath.Join(networkOrdererConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/organizations/orgca_certs/%v/fabric-ca-server-config.yaml", networkOrdererConfig.Config.OrgName))

		replacements = map[string]string{
			"${ORGNAME}":   networkOrdererConfig.Config.OrgName,
			"${ORGCAPORT}": networkOrdererConfig.Config.Ca.OrgcaPort,
		}
		err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/orderer/orgca/orgca_compose.yaml"
		outputFilePath = filepath.Join(networkOrdererConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/compose_files/orgca/%v/orgca.yaml", networkOrdererConfig.Config.OrgName))

		replacements = map[string]string{
			"${ORGNAME}":   networkOrdererConfig.Config.OrgName,
			"${ORGCAPORT}": networkOrdererConfig.Config.Ca.OrgcaPort,
			"${NETWORK}":   networkOrdererConfig.Network.NetworkName,
		}
		err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/orderer/orgca/orgcerts.sh"
		outputFilePath = filepath.Join(networkOrdererConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/script_files/orgcerts/%v/orgcerts.sh", networkOrdererConfig.Config.OrgName))

		replacements = map[string]string{
			"${ORGNAME}":   networkOrdererConfig.Config.OrgName,
			"${ORGCAPORT}": networkOrdererConfig.Config.Ca.OrgcaPort,
			"${NETWORK}":   networkOrdererConfig.Network.NetworkName,
		}
		err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/orderer/orderer/orderer.yaml"
		outputFilePath = filepath.Join(networkOrdererConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/compose_files/peer_couch_orderer/%v/orderer.yaml", networkOrdererConfig.Config.OrgName))

		replacements = map[string]string{
			"${ORGNAME}":       networkOrdererConfig.Config.OrgName,
			"${ORGNAMEFC}":     strings.ToUpper(string(networkOrdererConfig.Config.OrgName[0])) + networkOrdererConfig.Config.OrgName[1:],
			"${NETWORK}":       networkOrdererConfig.Network.NetworkName,
			"${ORDERER0LIS}":   networkOrdererConfig.Config.Orderers[0].OrdererLis,
			"${ORDERER0ADM}":   networkOrdererConfig.Config.Orderers[0].OrdererAdm,
			"${ORDERER0ADMOP}": networkOrdererConfig.Config.Orderers[0].OrdererAdmOp,
			"${ORDERER1LIS}":   networkOrdererConfig.Config.Orderers[1].OrdererLis,
			"${ORDERER1ADM}":   networkOrdererConfig.Config.Orderers[1].OrdererAdm,
			"${ORDERER1ADMOP}": networkOrdererConfig.Config.Orderers[1].OrdererAdmOp,
			"${ORDERER2LIS}":   networkOrdererConfig.Config.Orderers[2].OrdererLis,
			"${ORDERER2ADM}":   networkOrdererConfig.Config.Orderers[2].OrdererAdm,
			"${ORDERER2ADMOP}": networkOrdererConfig.Config.Orderers[2].OrdererAdmOp,
		}
		err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/orderer/channel/osn_orderer_join.sh"
		outputFilePath = filepath.Join(networkOrdererConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/script_files/channel/%v/%v/osn_orderer_join.sh", networkOrdererConfig.Config.OrgName, networkOrdererConfig.Network.ChannelName))

		replacements = map[string]string{
			"${ORGNAME}":     networkOrdererConfig.Config.OrgName,
			"${ORGNAMEC}":    strings.ToUpper(networkOrdererConfig.Config.OrgName),
			"${NETWORK}":     networkOrdererConfig.Network.NetworkName,
			"${CHANNELNAME}": networkOrdererConfig.Network.ChannelName,
			"${ORDERER0ADM}": networkOrdererConfig.Config.Orderers[0].OrdererAdm,
			"${ORDERER1ADM}": networkOrdererConfig.Config.Orderers[1].OrdererAdm,
			"${ORDERER2ADM}": networkOrdererConfig.Config.Orderers[2].OrdererAdm,
		}
		err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}

		templateFilePath = "../../internal/templates/orderer/channel/utils.sh"
		outputFilePath = filepath.Join(networkOrdererConfig.Network.NetworkDirectory, fmt.Sprintf("network_files/script_files/channel/%v/utils.sh", networkOrdererConfig.Config.OrgName))

		data, err := ReadFile(templateFilePath)
		if err != nil {
			return err
		}
		err = WriteFile(outputFilePath, data)
		if err != nil {
			return err
		}

		fmt.Println(" All the orderer files have been successfully created.")
		return nil

	}

	fmt.Println("Model type issue")
	return errors.New("MODEL TYPE ISSUE")
}

func CreateConfigtx(networkConfig *models.Network, policyMap *map[string]string) error {

	outputFilePath := filepath.Join(networkConfig.NetworkDirectory, fmt.Sprintf("network_files/channel_files/configtx_files/%v/configtx.yaml", networkConfig.ChannelName))

	orgConfigs := []models.Org_Config{}

	//A map where key is orgname and the value is it struct. I am doing this because i wont have to create variable before hand"

	for org, orgConfigPath := range networkConfig.Orgs {
		if org == "orderer" {
			continue
		}
		var orgConfig models.Org_Config
		err := ReadJson(&orgConfigPath, &orgConfig)
		if err != nil {
			return err
		}
		orgConfigs = append(orgConfigs, orgConfig)
	}

	var ordererConfig models.Orderer_Config
	var ordererConfigPath = networkConfig.Orgs["orderer"]
	err := ReadJson(&ordererConfigPath, &ordererConfig)
	if err != nil {
		return err
	}

	_, err = os.Stat(outputFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("Error occured while checking if the configtx file exists:", err)
			return err
		}
	} else {
		// File exists, delete it
		err := os.Remove(outputFilePath)
		if err != nil {
			fmt.Println("Error deleting the existing configtx file:", err)
			return err
		}
	}

	templateFilePath := "../../internal/templates/generic/configtx/orderer_section.yaml"
	replacements := map[string]string{
		"${ORDERER0LIS}": ordererConfig.Orderers[0].OrdererLis,
		"${ORDERER1LIS}": ordererConfig.Orderers[1].OrdererLis,
		"${ORDERER2LIS}": ordererConfig.Orderers[2].OrdererLis,
		"${NETWORK}":     networkConfig.NetworkName,
	}
	err = ProcessTemplatesAppend(&templateFilePath, &outputFilePath, &replacements)
	if err != nil {
		return err
	}

	templateFilePath = "../../internal/templates/generic/configtx/org_section.yaml"
	for _, orgConfig := range orgConfigs {
		replacements = map[string]string{
			"${ORGNAME}":     strings.ToLower(orgConfig.OrgName),
			"${ORGNAMEFC}":   strings.ToUpper(string(orgConfig.OrgName[0])) + strings.ToLower(orgConfig.OrgName[1:]),
			"${NETWORK}":     networkConfig.NetworkName,
			"${ORDERERNAME}": orgConfig.Orderer.OrdererName,
			"${ORDERERPORT}": orgConfig.Orderer.OrdererLis,
		}
		err = ProcessTemplatesAppend(&templateFilePath, &outputFilePath, &replacements)
		if err != nil {
			return err
		}
	}

	policyMapDeref := *(policyMap)
	adminPolicy := policyMapDeref["admin"]
	lifecyclePolicy := policyMapDeref["lifecycle"]
	endrsementPolicy := policyMapDeref["endorsement"]
	templateFilePath = "../../internal/templates/generic/configtx/cap_app.yaml"

	replacements = map[string]string{
		"${ADMIN}":       adminPolicy,
		"${LIFECYCLE}":   lifecyclePolicy,
		"${ENDORSEMENT}": endrsementPolicy,
	}
	err = ProcessTemplatesAppend(&templateFilePath, &outputFilePath, &replacements)
	if err != nil {
		return err
	}

	templateFilePath = "../../internal/templates/generic/configtx/orderer_channel_defaults.yaml"

	replacements = map[string]string{
		"${NETWORK}":     networkConfig.NetworkName,
		"${ORDERER0LIS}": ordererConfig.Orderers[0].OrdererLis,
		"${ORDERER1LIS}": ordererConfig.Orderers[1].OrdererLis,
		"${ORDERER2LIS}": ordererConfig.Orderers[2].OrdererLis,
	}
	err = ProcessTemplatesAppend(&templateFilePath, &outputFilePath, &replacements)
	if err != nil {
		return err
	}

	templateFilePath = "../../internal/templates/generic/configtx/profiles.yaml"

	var profileString string
	spaces := "                "
	for index, orgConfig := range orgConfigs {

		if index == 0 {
			profileString = profileString + fmt.Sprintf("- *%v\n", strings.ToUpper(string(orgConfig.OrgName[0]))+strings.ToLower(orgConfig.OrgName[1:]))
		} else if index == len(orgConfigs)-1 {
			profileString = profileString + fmt.Sprintf("%s- *%v", spaces, strings.ToUpper(string(orgConfig.OrgName[0]))+strings.ToLower(orgConfig.OrgName[1:]))
		} else {
			profileString = profileString + fmt.Sprintf("%s- *%v\n", spaces, strings.ToUpper(string(orgConfig.OrgName[0]))+strings.ToLower(orgConfig.OrgName[1:]))
		}
	}
	replacements = map[string]string{
		"${ORGS}": profileString,
	}
	err = ProcessTemplatesAppend(&templateFilePath, &outputFilePath, &replacements)
	if err != nil {
		return err
	}

	fmt.Println("The configtx file is created successfully.")

	return nil

}
func CreatePrivateData(networkConfig *models.Network, policy *string, inputOrgs *[]string) error {

	outputFilePath := filepath.Join(networkConfig.NetworkDirectory, fmt.Sprintf("network_files/privatedata_files/%v/collection_config.json", networkConfig.ChannelName))

	pdName := ""
	for _, org := range *inputOrgs {
		pdName = pdName + strings.ToLower(org)

	}

	//A map where key is orgname and the value is it struct. I am doing this because i wont have to create variable before hand"

	_, err := os.Stat(outputFilePath)

	pdDetails := models.PdDetails{}
	pdDetailsSlice := []models.PdDetails{}
	if os.IsNotExist(err) {
		pdDetails.Name = pdName
		pdDetails.Policy = *policy
		pdDetails.RequiredPeerCount = 1
		pdDetails.MaxPeerCount = 2
		pdDetails.BlockToLive = 0
		pdDetails.MemberOnlyRead = false
		pdDetails.MemberOnlyWrite = true
		pdDetails.EndorsementPolicy = models.PdEndorsementPolicy{
			SignaturePolicy: *policy,
		}
		pdDetailsSlice = append(pdDetailsSlice, pdDetails)
		err = OutputJson(&pdDetailsSlice, &outputFilePath)
		if err != nil {
			return err
		}

	} else if err == nil {
		err = ReadJson(&outputFilePath, &pdDetailsSlice)
		if err != nil {
			return err
		}
		pdDetails.Name = pdName
		pdDetails.Policy = *policy
		pdDetails.RequiredPeerCount = 1
		pdDetails.MaxPeerCount = 2
		pdDetails.BlockToLive = 0
		pdDetails.MemberOnlyRead = false
		pdDetails.MemberOnlyWrite = true
		pdDetails.EndorsementPolicy = models.PdEndorsementPolicy{
			SignaturePolicy: *policy,
		}
		pdDetailsSlice = append(pdDetailsSlice, pdDetails)
		err = OutputJson(&pdDetailsSlice, &outputFilePath)
		if err != nil {
			return err
		}

	}
	return nil

}

func ProcessTemplates(templateFilePath, outputFilePath *string, replacements *map[string]string) error {

	byteData, err := ReadFile(*templateFilePath)
	if err != nil {
		return err
	}
	templateData := string(byteData)
	for src, dest := range *replacements {
		templateData = strings.ReplaceAll(templateData, src, dest)
	}

	templateDataByte := []byte(templateData)

	err = WriteFile(*outputFilePath, templateDataByte)
	if err != nil {
		return err
	}
	return nil
}

func ProcessTemplatesAppend(templateFilePath, outputFilePath *string, replacements *map[string]string) error {
	// Read the template file
	byteData, err := ReadFile(*templateFilePath)
	if err != nil {
		return err
	}
	templateData := string(byteData)

	// Perform replacements
	for src, dest := range *replacements {
		templateData = strings.ReplaceAll(templateData, src, dest)
	}

	// Open the output file in append mode
	file, err := os.OpenFile(*outputFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Convert templateData to a byte slice and write to the file
	_, err = file.WriteString(templateData)
	if err != nil {
		return err
	}

	return nil
}

func ReadJson(path *string, holder interface{}) error {
	byteValue, err := ReadFile(*path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, holder)
	if err != nil {
		fmt.Printf("Following error occured while unmarshalling json : %v", err)
		return err
	}
	//fmt.Println(holder)
	//fmt.Println(reflect.TypeOf(holder))

	return nil

}
func FindOrgType(path *string) (string, error) {

	byteValue, err := ReadFile(*path)
	if err != nil {
		fmt.Printf("Error reading file: %v", err)
		fmt.Println()
		return "", err
	}
	var data map[string]interface{}
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		fmt.Printf("Following error occured while unmarshalling json : %v", err)
		return "", err
	}
	orgTypeInterface, ok := data["orgType"]
	if !ok {
		fmt.Println("OrgType field does not exist")
		return "", errors.New("")
	}

	orgType, ok := orgTypeInterface.(string)
	orgType = strings.ToLower(orgType)
	if !ok {
		fmt.Println("Type assertion from interface to orgtype failed")
		return "", errors.New("")
	}
	//fmt.Println(holder)
	//fmt.Println(reflect.TypeOf(holder))

	return orgType, nil

}
func OutputJson(data interface{}, path *string) error {

	var networkConfig models.Network
	var orgConfig models.Org_Config
	var ordererConfig models.Orderer_Config

	networkConfigType := reflect.TypeOf(networkConfig)
	orgConfigType := reflect.TypeOf(orgConfig)
	ordererConfigType := reflect.TypeOf(ordererConfig)
	varType := reflect.TypeOf(data).Elem() //.Elem because we receive a pointer to struct.
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Following error occured while converting struct to json: %v", err)
		return err
	}

	err = WriteFile(*path, jsonData)

	if err != nil {
		return err
	}
	if varType == networkConfigType {
		fmt.Printf("\nA config.json file is created/edited and stored at %v. Please dont delete,move or edit the file.", *path)
		fmt.Println()
	}
	if varType == orgConfigType || varType == ordererConfigType {
		fmt.Printf("\n A configuration file for the organisation is created and stored at %v. Please dont delete,move or edit the file.", *path)
		fmt.Println()
	}

	return nil

}

func ReadFile(filepath string) ([]byte, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("\nFollowing error occured while reading the file %v : %v", filepath, err)
		return nil, err
	}

	return data, nil
}

func WriteFile(filepath string, data []byte) error {
	err := os.WriteFile(filepath, data, 0644)
	if err != nil {
		fmt.Printf("\nFollowing error occurred while writing to the file %v: %v", filepath, err)
		return err
	}

	return nil
}
