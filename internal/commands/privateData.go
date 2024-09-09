package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/KshitijPatil98/fabriquik/internal/config"
	"github.com/KshitijPatil98/fabriquik/internal/models"
	"github.com/KshitijPatil98/fabriquik/internal/utils"
)

func PrivateData() {
	privateDataCmd := flag.NewFlagSet("privatedata", flag.ExitOnError)

	endorsement := privateDataCmd.String("endorsement", "", "[REQUIRED] The type of the endorsement policy. Available options are: 1. And 2. Or")
	orgs := privateDataCmd.String("orgs", "", "[REQUIRED] The comma seperated orgnames. Only set this flag if the type of endorsement is set to AND/OR.")

	if len(os.Args) < 3 {
		fmt.Println("Error: No flags supplied")
		fmt.Println(`Run "fabriquik privatedata --help" for usage`)
		return

	}
	err := privateDataCmd.Parse(os.Args[2:])
	if err != nil {
		fmt.Printf("Error occured while parsing: %v", err)
		return
	}
	if *endorsement == "" {
		fmt.Println(`Error: --endorsement flag is not set.`)
		fmt.Println("Please set the flag and try again")
		fmt.Println("Example usage :   fabriquik privatedata  --endorsement OR --orgs incalus")
		return

	}
	if *orgs == "" {
		fmt.Println(`Error: --orgs flag is not set.`)
		fmt.Println("Please set the flag and try again")
		fmt.Println("Example usage :   fabriquik privatedata  --endorsement OR --orgs incalus")
		return

	}
	if *endorsement != "AND" && *endorsement != "OR" {
		fmt.Println(`Error: Invalid value for --endorsement flag.`)
		fmt.Println("Available options are: \n1. AND 2. OR ")
		return

	}
	var networkConfig models.Network
	err = utils.ReadJson(&config.ConfigFilePath, &networkConfig)
	if err != nil {
		return
	}
	// Slice to hold the keys

	orgMap := map[string]bool{}
	// Extracting the keys

	for org := range networkConfig.Orgs {
		if strings.ToLower(org) == "orderer" {
			continue
		}
		orgMap[org] = true
	}

	// Trim any leading or trailing whitespace
	*orgs = strings.TrimSpace(*orgs)

	// Check if the string contains at least one comma
	if !strings.Contains(*orgs, ",") {
		if !(len(strings.Split(*orgs, ",")) == 1) {
			fmt.Println("The value supplied for --orgs flag is invalid. Please supply a comman seperate value.")
			fmt.Println("Example usage :   fabriquik configtx --endorsement And --orgs incalus,bolthouse")
			return
		}

	}
	inputOrgs := strings.Split(*orgs, ",")

	for _, inputOrg := range inputOrgs {
		if !orgMap[inputOrg] {
			fmt.Printf(`The org "%v" was not onboarded. Please make sure you are specifying a correct org.`, inputOrg)
			fmt.Println()
			return
		}
	}
	var policy string
	for index, org := range inputOrgs {

		orgfcMsp := fmt.Sprintf("%vMSP", strings.ToUpper(string(org[0]))+strings.ToLower(org[1:]))
		if index == len(inputOrgs)-1 {
			policy = policy + fmt.Sprintf(`'%v.member'`, orgfcMsp)
		} else {
			policy = policy + fmt.Sprintf(`'%v.member',`, orgfcMsp)
		}
	}
	policy = fmt.Sprintf("%v(%v)", *endorsement, policy)

	err = utils.CreatePrivateData(&networkConfig, &policy, &inputOrgs)
	if err != nil {
		return
	}

}
