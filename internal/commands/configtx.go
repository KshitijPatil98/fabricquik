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

func Configtx() {

	configtxCmd := flag.NewFlagSet("configtx", flag.ExitOnError)

	endorsement := configtxCmd.String("endorsement", "", "[REQUIRED] The type of the endorsement policy. Same policy is applied to lifecycle endorsement and admin endorsement as well. Available options are: 1. Majority\n2. Any\n3. And\n4. Or\nIf And or Or type is selected please specify the comma seperated orgs with the orgs flag.")
	orgs := configtxCmd.String("orgs", "", "[OPTIONAL] The comma seperated orgnames. Only set this flag if the type of endorsement is set to and/or.")

	if len(os.Args) < 3 {
		fmt.Println("Error: No flags supplied")
		fmt.Println(`Run "fabriquik configtx --help" for usage`)
		return

	}
	var networkConfig models.Network
	err := configtxCmd.Parse(os.Args[2:])
	if err != nil {
		fmt.Printf("Error occured while parsing: %v", err)
		return
	}
	if *endorsement == "" {
		fmt.Println(`Error: --endorsement flag is not set.`)
		fmt.Println("Please set the flag and try again")
		fmt.Println("Example usage :   fabriquik configtx --endorsement ANY")
		return

	}
	*endorsement = strings.ToUpper(*endorsement)
	if *endorsement != "MAJORITY" && *endorsement != "ANY" && *endorsement != "ALL" && *endorsement != "AND" && *endorsement != "OR" {
		fmt.Println(`Error: Invalid value for --endorsement flag.`)
		fmt.Println("Available options are: \n1. MAJORITY\n2. ANY\n3. ALL\n4. AND\n5. OR")
		return

	}

	if (*endorsement == "AND" || *endorsement == "OR") && *orgs == "" {
		fmt.Println(`Error: --orgs flag is not set. If the endorsement type is "AND" or "OR", the orgs flag need to be set,`)
		fmt.Println("Example usage :   fabriquik configtx --endorsement AND --orgs incalus,bolthouse")
		return

	}
	err = utils.ReadJson(&config.ConfigFilePath, &networkConfig)
	policyMap := map[string]string{}

	if *endorsement == "MAJORITY" || *endorsement == "ANY" || *endorsement == "ALL" {

		policyMap["type"] = "ImplicitMeta"
		policyMap["admin"] = fmt.Sprintf("%v Admins", *endorsement)
		policyMap["lifecycle"] = fmt.Sprintf("%v Endorsement", *endorsement)
		policyMap["endorsement"] = fmt.Sprintf("%v Endorsement", *endorsement)

	}

	if *endorsement == "AND" || *endorsement == "OR" {

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

		policyMap["type"] = "Signature"

		var adminPolicy string
		var normalPolicy string
		for index, org := range inputOrgs {

			orgfcMsp := fmt.Sprintf("%vMSP", strings.ToUpper(string(org[0]))+strings.ToLower(org[1:]))
			if index == len(inputOrgs)-1 {
				adminPolicy = adminPolicy + fmt.Sprintf(`'%v.admin'`, orgfcMsp)
				normalPolicy = normalPolicy + fmt.Sprintf(`'%v.peer'`, orgfcMsp)
			} else {
				adminPolicy = adminPolicy + fmt.Sprintf(`'%v.admin',`, orgfcMsp)
				normalPolicy = normalPolicy + fmt.Sprintf(`'%v.peer',`, orgfcMsp)
			}
		}
		adminPolicy = fmt.Sprintf("%v(%v)", *endorsement, adminPolicy)
		normalPolicy = fmt.Sprintf("%v(%v)", *endorsement, normalPolicy)
		policyMap["admin"] = fmt.Sprintf("%v", adminPolicy)
		policyMap["lifecycle"] = fmt.Sprintf("%v", normalPolicy)
		policyMap["endorsement"] = fmt.Sprintf("%v", normalPolicy)

	}

	err = utils.CreateConfigtx(&networkConfig, &policyMap)
	if err != nil {
		return
	}

}
