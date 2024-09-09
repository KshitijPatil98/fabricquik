#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

#VVVIP This file runs inside cli container so path is set accordingly. Check the cli container to see the volume mapping. This diff is only for two files setAnchorPeer.sh and configUpdate.sh 


#ALSO ANCHOR PEER IS SET FOR AN ORG. SO ONCE THE ORG HAS JOINED THE CHANNEL THERE IS NO NEED TO GET ANY SIGNATURES From any other org


. ${PWD}/script_files/channel/${ORGNAME}/envVar.sh
. ${PWD}/script_files/channel/${ORGNAME}/configUpdate.sh
. ${PWD}/script_files/channel/${ORGNAME}/utils.sh


# NOTE: this must be run in a CLI container since it requires jq and configtxlator 
createAnchorPeerUpdate() {
  infoln "Fetching channel config for channel $CHANNEL_NAME"
  fetchChannelConfig $PEER $CHANNEL_NAME ${CORE_PEER_LOCALMSPID}config.json

  infoln "Generating anchor peer update transaction for peer${PEER} of ${ORGNAME} on channel $CHANNEL_NAME"

  if [ $PEER -eq 0 ]; then
    HOST="peer0.${ORGNAME}.${NETWORK}.com"
    PORT=${PEER0LIS}
  elif [ $PEER -eq 1 ]; then
    HOST="peer1.${ORGNAME}.${NETWORK}.com"
    PORT=${PEER1LIS}
  else
    errorln "Peer${PEER} unknown"
  fi

  set -x
  # Modify the configuration to append the anchor peer 
  jq '.channel_group.groups.Application.groups.'${CORE_PEER_LOCALMSPID}'.values += {"AnchorPeers":{"mod_policy": "Admins","value":{"anchor_peers": [{"host": "'$HOST'","port": '$PORT'}]},"version": "0"}}' ${CORE_PEER_LOCALMSPID}config.json > ${CORE_PEER_LOCALMSPID}modified_config.json
  { set +x; } 2>/dev/null

  # Compute a config update, based on the differences between 
  # {orgmsp}config.json and {orgmsp}modified_config.json, write
  # it as a transaction to {orgmsp}anchors.tx
  createConfigUpdate ${CHANNEL_NAME} ${CORE_PEER_LOCALMSPID}config.json ${CORE_PEER_LOCALMSPID}modified_config.json ${CORE_PEER_LOCALMSPID}anchors.tx
}

updateAnchorPeer() {
  peer channel update -o ${ORDERERNAME}.${NETWORK}.com:${ORDERERPORT} --ordererTLSHostnameOverride ${ORDERERNAME}.${NETWORK}.com -c $CHANNEL_NAME -f ${CORE_PEER_LOCALMSPID}anchors.tx --tls --cafile "$ORDERER_CA" >&log.txt
  res=$?
  cat log.txt
  verifyResult $res "Anchor peer update failed"
  successln "Anchor peer set for '$CORE_PEER_LOCALMSPID' on channel '$CHANNEL_NAME'"
}

PEER=$1
CHANNEL_NAME=$2

setGlobals $PEER

createAnchorPeerUpdate 

updateAnchorPeer 
