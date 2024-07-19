#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

# import utils
. ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${OWNERNAME_ALL_LOWER}/utils.sh
. ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${OWNERNAME_ALL_LOWER}/${CHANNELNAME}/${ORGNAME_ALL_LOWER}/configUpdate.sh

export FABRIC_CFG_PATH=${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${OWNERNAME_ALL_LOWER}/${CHANNELNAME}/${ORGNAME_ALL_LOWER}

configtxgen -printOrg ${ORGNAME_FIRST_UPPER}MSP > ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/${ORGNAME_ALL_LOWER}.json

export FABRIC_CFG_PATH=${PWD}/../config/

infoln "Creating config transaction to add new org to network"

# Fetch the config for the channel, writing it to config.json
fetchChannelConfig 0 ${CHANNELNAME}

# Modify the configuration to append the new org
#In below statement if we directly use the env variable as ${ORGNAME_FIRST_UPPER} it doesnt work because everything is inside single quotes. So to solve that i have added another single quotes around my env variable
set -x
jq -s '.[0] * {"channel_group":{"groups":{"Application":{"groups": {"${ORGNAME_FIRST_UPPER}MSP":.[1]}}}}}' ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/config.json ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/${ORGNAME_ALL_LOWER}.json > ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/modified_config.json
{ set +x; } 2>/dev/null

# Compute a config update, based on the differences between config.json and modified_config.json, write it as a transaction to org3_update_in_envelope.pb
createConfigUpdate ${CHANNELNAME} ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/config.json ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/modified_config.json 

#infoln "Signing config transaction"
#signConfigtxAsPeerOrg 0 

# Update command automatically signs it so no need to call sign seperately. Also by default admins policy is referred for adding members to channel. You can see that in the application section in the configtx file. The policy by defaults requires majority of the admins to sign it but we will change it to only admin of the owner org need to sign
infoln "Signing and sending an update transaction as admin of owner org."
updateConfigtxAsPeerOrg 0 


successln "Config transaction to add the new org submitted successfully"

