#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

. ${PWD}/scm_script_files/channel_scripts/${OWNERNAME_ALL_LOWER}/envVar.sh
. ${PWD}/scm_script_files/channel_scripts/${OWNERNAME_ALL_LOWER}/utils.sh

# fetchChannelConfig <org> <channel_id> <output_json>
# Writes the current channel config for a given channel to a JSON file
# NOTE: this must be run in a CLI container since it requires configtxlator
fetchChannelConfig() {
  PEER=$1
  CHANNEL=$2
  setGlobals $PEER

  infoln "Fetching the most recent configuration block for the channel"
  set -x
  peer channel fetch config ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/config_block.pb -o localhost:${ORDERER_NODE_LISPORT} --ordererTLSHostnameOverride ${ORDERER_NODE_NAME}.${PROJECT_NAME}.com -c $CHANNEL --tls --cafile "$ORDERER_CA"
  { set +x; } 2>/dev/null

  infoln "Decoding config block to JSON and isolating config and storing it in the location seen below"
  set -x
  configtxlator proto_decode --input ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/config_block.pb --type common.Block --output ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/config_block.json
  
  jq ".data.data[0].payload.data.config" ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/config_block.json > ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/config.json
  { set +x; } 2>/dev/null
}

# createConfigUpdate <channel_id> <original_config.json> <modified_config.json> <output.pb>
# Takes an original and modified config, and produces the config update tx
# which transitions between the two
# NOTE: this must be run in a CLI container since it requires configtxlator
createConfigUpdate() {
  CHANNEL=$1
  ORIGINAL=$2
  MODIFIED=$3
  

  set -x
  configtxlator proto_encode --input ${ORIGINAL} --type common.Config --output ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/original_config.pb
  configtxlator proto_encode --input ${MODIFIED} --type common.Config --output ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/modified_config.pb
  configtxlator compute_update --channel_id ${CHANNEL} --original ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/original_config.pb --updated ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/modified_config.pb --output ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/config_update.pb
  configtxlator proto_decode --input ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/config_update.pb --type common.ConfigUpdate --output ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/config_update.json
  echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL'", "type":2}},"data":{"config_update":'$(cat ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/config_update.json)'}}}' | jq . >${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/config_update_in_envelope.json
  configtxlator proto_encode --input ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/config_update_in_envelope.json --type common.Envelope --output ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/${ORGNAME_ALL_LOWER}_update_in_envelope.pb
  { set +x; } 2>/dev/null
}

# signConfigtxAsPeerOrg <org> <configtx.pb>
# Set the peerOrg admin of an org and sign the config update
signConfigtxAsPeerOrg() {
  PEER=$1
  setGlobals $PEER
  set -x
  peer channel signconfigtx -f ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/${ORGNAME_ALL_LOWER}_update_in_envelope.pb
  { set +x; } 2>/dev/null
}

updateConfigtxAsPeerOrg() {
  PEER=$1
  setGlobals $PEER
  set -x
  peer channel update -f ${PWD}/${PROJECT_NAME}_script_files/channel_scripts/${ORGNAME_ALL_LOWER}/${CHANNELNAME}/channel_update_files/${ORGNAME_ALL_LOWER}_update_in_envelope.pb -c ${CHANNELNAME} -o localhost:${ORDERER_NODE_LISPORT} --ordererTLSHostnameOverride ${ORDERER_NODE_NAME}.${PROJECT_NAME}.com --tls --cafile "$ORDERER_CA"
{ set +x; } 2>/dev/null
}

