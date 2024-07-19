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

	"github.com/KshitijPatil98/fabriquik/internal/models"
)

func CheckSetupFolders(networkConfig *models.Network) (bool, bool, error) {

	command := exec.Command("bash", "-c", "ls")
	command.Dir = networkConfig.NetworkDirectory
	// Run the command and check for errors
	outputByte, err := command.CombinedOutput()

	if err != nil {
		fmt.Printf("Following error occured while runnung the ls command: %v", err)
		return false, false, err
	}
	outputStr := string(outputByte)

	folderNames := strings.Split(outputStr, "\n")
	configPresent := false
	chaincodePresent := false
	for _, folderName := range folderNames {

		if folderName == "config" {
			configPresent = true
		}
		if folderName == networkConfig.ChaincodeName {
			chaincodePresent = true
		}

	}

	return configPresent, chaincodePresent, nil
}

func CheckOnboardFolders(path *string) (bool, error) {

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
	configPresent := false
	networkFilesPresent := false
	configJsonPresent := false
	for _, folderName := range folderNames {

		if folderName == "config" {
			configPresent = true
		}
		if folderName == "network_files" {
			networkFilesPresent = true
		}
		if folderName == "networkOrgConfig.json" {
			configJsonPresent = true
		}

	}
	if configPresent && networkFilesPresent && configJsonPresent {
		return true, nil
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
			"network_files/compose_files/peer_compose_files",
			"network_files/compose_files/orderer_compose_files",
			"network_files/compose_files/orgca_compose_files",
			"network_files/compose_files/tlsca_compose_files",
			"network_files/connection_files",
			"network_files/explorer_files",
			"network_files/script_files/ccp_scripts",
			"network_files/script_files/channel_scripts",
			"network_files/script_files/orgcerts_scripts",
			"network_files/script_files/tlscerts_scripts",
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
		fmt.Println("\nAll the directories have been created successfully. You can go ahead and onboard the owner")
		return nil

	}
	return nil

}

func CreateOrgDirectories(gen_config interface{}) error {

	var networkOrgConfig models.Network_Org_Config
	var networkOrdererConfig models.Network_Orderer_Config
	var ok bool
	networkOrgConfig, ok = gen_config.(models.Network_Org_Config)
	if !ok {

		networkOrdererConfig, ok = gen_config.(models.Network_Orderer_Config)
		if !ok {
			fmt.Println("Invalid model type passed.")
			return errors.New("INVALID MODEL TYPE PASSED")
		}

	}

	if reflect.ValueOf(networkOrgConfig).IsZero() {
		fmt.Println(networkOrdererConfig)
		return nil
	}

	if networkOrgConfig.Config.OrgType == "owner" {
		directories := []string{
			fmt.Sprintf("network_files/organizations/tlsca_certs/%v", networkOrgConfig.Config.OrgName),
			fmt.Sprintf("network_files/organizations/tlsca_certs/%v/users/tlsca_admin", networkOrgConfig.Config.OrgName),
			fmt.Sprintf("network_files/organizations/orgca_certs/%v", networkOrgConfig.Config.OrgName),
			fmt.Sprintf("network_files/organizations/orgca_certs/%v/tls", networkOrgConfig.Config.OrgName),

			fmt.Sprintf("network_files/compose_files/tlsca_compose_files/%v", networkOrgConfig.Config.OrgName),
			fmt.Sprintf("network_files/compose_files/orgca_compose_files/%v", networkOrgConfig.Config.OrgName),
			fmt.Sprintf("network_files/compose_files/peer_compose_files/%v", networkOrgConfig.Config.OrgName),

			fmt.Sprintf("network_files/script_files/tlscerts_scripts/%v", networkOrgConfig.Config.OrgName),
			fmt.Sprintf("network_files/script_files/orgcerts_scripts/%v", networkOrgConfig.Config.OrgName),
			fmt.Sprintf("network_files/script_files/channel_scripts/%v/%v", networkOrgConfig.Config.OrgName, networkOrgConfig.Network.ChaincodeName),
			fmt.Sprintf("network_files/script_files/ccp_scripts/%v", networkOrgConfig.Config.OrgName),

			fmt.Sprintf("network_files/connection_files/%v", networkOrgConfig.Config.OrgName),
		}
		for _, dir := range directories {
			fullPath := filepath.Join(networkOrgConfig.Config.NetworkDirectory, dir)
			err := os.MkdirAll(fullPath, os.ModePerm)
			if err != nil {
				fmt.Printf("Error occured while creating directory %s : %v", fullPath, err)
				fmt.Printf("The partially created directories will now be deleted")
				err := os.RemoveAll(networkOrgConfig.Config.NetworkDirectory)
				if err != nil {
					fmt.Printf("Error occured while cleaning partially created directories : %v. Please manually clean the directores", err)
					return err
				}
				fmt.Printf("Partially created directories are cleaned. Please rectify the error and try again")
				return err
			}

		}
		fmt.Println("All the org directories have been  successfully created. You can go ahead and generate the files for the owner")
		return nil

	} else {
		return nil
	}
}

func CreateOrgFiles(gen_config interface{}) error {

	// templateMap := map[string]string{

	// 	"../templates/org/channel/join_peers.sh":           fmt.Sprintf("network_files/script_files/channel_scripts/%v/%v/join_peers.sh", *orgName, *channelName),
	// 	"../templates/org/channel/setAnchorPeer.sh":        fmt.Sprintf("network_files/script_files/channel_scripts/%v/%v/setAnchorPeer.sh", *orgName, *channelName),
	// 	"../templates/org/channel/ccutils.sh":              fmt.Sprintf("network_files/script_files/channel_scripts/%v/%v/ccutils.sh", *orgName, *channelName),
	// 	"../templates/org/channel/pkg_ins_app_com_que.sh":  fmt.Sprintf("network_files/script_files/channel_scripts/%v/%v/pkg_ins_app_com_que.sh", *orgName, *channelName),

	// 	"../templates/org/connection/ccp-generate.sh":          fmt.Sprintf("network_files/script_files/ccp_scripts/%v/ccp-generate.sh", *orgName),
	// 	"../templates/org/connection/ccp-template.json":        fmt.Sprintf("network_files/script_files/ccp_scripts/%v/ccp-template.json", *orgName),
	// 	"../templates/org/connection/copy_connection_files.sh": fmt.Sprintf("network_files/script_files/ccp_scripts/%v/copy_connection_files.sh", *orgName),
	// }

	var networkOrgConfig models.Network_Org_Config
	var networkOrdererConfig models.Network_Orderer_Config
	var ok bool
	networkOrgConfig, ok = gen_config.(models.Network_Org_Config)
	if !ok {

		networkOrdererConfig, ok = gen_config.(models.Network_Orderer_Config)
		if !ok {
			fmt.Println("Invalid model type passed.")
			return errors.New("INVALID MODEL TYPE PASSED")
		}

	}

	if reflect.ValueOf(networkOrgConfig).IsZero() {
		fmt.Println(networkOrdererConfig)
		return nil
	}

	templateFilePath := "../../internal/templates/org/tlsca/fabric-ca-server-config.yaml"
	outputFilePath := filepath.Join(networkOrgConfig.Config.NetworkDirectory, fmt.Sprintf("network_files/organizations/tlsca_certs/%v/fabric-ca-server-config.yaml", networkOrgConfig.Config.OrgName))

	var replacements = map[string]string{
		"${ORGNAME}":   networkOrgConfig.Config.OrgName,
		"${TLSCAPORT}": networkOrgConfig.Config.Ca.TlscaPort,
	}
	err := ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
	if err != nil {
		return err
	}

	templateFilePath = "../../internal/templates/org/tlsca/tlsca_compose.yaml"
	outputFilePath = filepath.Join(networkOrgConfig.Config.NetworkDirectory, fmt.Sprintf("network_files/compose_files/tlsca_compose_files/%v/tlsca_compose.yaml", networkOrgConfig.Config.OrgName))

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
	outputFilePath = filepath.Join(networkOrgConfig.Config.NetworkDirectory, fmt.Sprintf("network_files/script_files/tlscerts_scripts/%v/tlscerts.sh", networkOrgConfig.Config.OrgName))

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
	outputFilePath = filepath.Join(networkOrgConfig.Config.NetworkDirectory, fmt.Sprintf("network_files/organizations/orgca_certs/%v/fabric-ca-server-config.yaml", networkOrgConfig.Config.OrgName))

	replacements = map[string]string{
		"${ORGNAME}":   networkOrgConfig.Config.OrgName,
		"${ORGCAPORT}": networkOrgConfig.Config.Ca.OrgcaPort,
	}
	err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
	if err != nil {
		return err
	}

	templateFilePath = "../../internal/templates/org/orgca/orgca_compose.yaml"
	outputFilePath = filepath.Join(networkOrgConfig.Config.NetworkDirectory, fmt.Sprintf("network_files/compose_files/orgca_compose_files/%v/orgca_compose.yaml", networkOrgConfig.Config.OrgName))

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
	outputFilePath = filepath.Join(networkOrgConfig.Config.NetworkDirectory, fmt.Sprintf("network_files/script_files/orgcerts_scripts/%v/orgcerts.sh", networkOrgConfig.Config.OrgName))

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
	outputFilePath = filepath.Join(networkOrgConfig.Config.NetworkDirectory, fmt.Sprintf("network_files/compose_files/peer_compose_files/%v/peer.yaml", networkOrgConfig.Config.OrgName))

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
	outputFilePath = filepath.Join(networkOrgConfig.Config.NetworkDirectory, fmt.Sprintf("network_files/compose_files/peer_compose_files/%v/couch.yaml", networkOrgConfig.Config.OrgName))

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
	outputFilePath = filepath.Join(networkOrgConfig.Config.NetworkDirectory, fmt.Sprintf("network_files/compose_files/peer_compose_files/%v/cli.yaml", networkOrgConfig.Config.OrgName))

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
	outputFilePath = filepath.Join(networkOrgConfig.Config.NetworkDirectory, fmt.Sprintf("network_files/script_files/channel_scripts/%v/configUpdate.sh", networkOrgConfig.Config.OrgName))

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
	outputFilePath = filepath.Join(networkOrgConfig.Config.NetworkDirectory, fmt.Sprintf("network_files/script_files/channel_scripts/%v/envVar.sh", networkOrgConfig.Config.OrgName))

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
	outputFilePath = filepath.Join(networkOrgConfig.Config.NetworkDirectory, fmt.Sprintf("network_files/script_files/channel_scripts/%v/utils.sh", networkOrgConfig.Config.OrgName))

	data, err := ReadFile(templateFilePath)
	if err != nil {
		return err
	}
	err = WriteFile(outputFilePath, data)
	if err != nil {
		return err
	}

	templateFilePath = "../../internal/templates/org/channel/create_genesis_block.sh"
	outputFilePath = filepath.Join(networkOrgConfig.Config.NetworkDirectory, fmt.Sprintf("network_files/script_files/channel_scripts/%v/%v/create_genesis_block.sh", networkOrgConfig.Config.OrgName, networkOrgConfig.Network.ChannelName))

	replacements = map[string]string{
		"${CHANNELNAME}": networkOrgConfig.Network.ChannelName,
	}
	err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
	if err != nil {
		return err
	}

	templateFilePath = "../../internal/templates/org/channel/join_peers.sh"
	outputFilePath = filepath.Join(networkOrgConfig.Config.NetworkDirectory, fmt.Sprintf("network_files/script_files/channel_scripts/%v/%v/join_peers.sh", networkOrgConfig.Config.OrgName, networkOrgConfig.Network.ChannelName))

	replacements = map[string]string{
		"${ORGNAME}":     networkOrgConfig.Config.OrgName,
		"${CHANNELNAME}": networkOrgConfig.Network.ChannelName,
	}
	err = ProcessTemplates(&templateFilePath, &outputFilePath, &replacements)
	if err != nil {
		return err
	}

	templateFilePath = "../../internal/templates/org/channel/setAnchorPeer.sh"
	outputFilePath = filepath.Join(networkOrgConfig.Config.NetworkDirectory, fmt.Sprintf("network_files/script_files/channel_scripts/%v/%v/setAnchorPeer.sh", networkOrgConfig.Config.OrgName, networkOrgConfig.Network.ChannelName))

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
	outputFilePath = filepath.Join(networkOrgConfig.Config.NetworkDirectory, fmt.Sprintf("network_files/script_files/channel_scripts/%v/%v/ccutils.sh", networkOrgConfig.Config.OrgName, networkOrgConfig.Network.ChannelName))

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

	templateFilePath = "../../internal/templates/org/channel/pkg_ins_app_com_que.sh"
	outputFilePath = filepath.Join(networkOrgConfig.Config.NetworkDirectory, fmt.Sprintf("network_files/script_files/channel_scripts/%v/%v/pkg_ins_app_com_que.sh", networkOrgConfig.Config.OrgName, networkOrgConfig.Network.ChannelName))

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

func ReadJson(path *string, holder interface{}) error {

	byteValue, err := ReadFile(*path)
	if err != nil {
		return err
	}

	//In case of using interface, the umarshal will store the objects as map[string]interface{} inside the interface array
	err = json.Unmarshal(byteValue, holder)
	if err != nil {
		fmt.Printf("Following error occured while unmarshalling json : %v", err)
		return err
	}
	//fmt.Println(holder)
	//fmt.Println(reflect.TypeOf(holder))

	return nil

}

func OutputJson(data interface{}, path *string) error {

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Following error occured while converting struct to json: %v", err)
		return err
	}

	err = WriteFile(*path, jsonData)

	if err != nil {
		return err
	}

	fmt.Printf("\nA configuration file has been created and stored at %v. Please make sure the file is not deleted and is always kept parallel to the network_files and config folder", *path)
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
