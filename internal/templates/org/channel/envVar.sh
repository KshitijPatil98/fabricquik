#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

# This is a collection of bash functions used by different scripts

# imports
. ${PWD}/script_files/channel_scripts/${ORGNAME}/utils.sh

export ORDERER_CA=${PWD}/organizations/ordererOrganizations/orderer.${NETWORK}.com/tlsca/tlsca.orderer.${NETWORK}.com-cert.pem


# Set environment variables for the peers of ${ORGNAME} org
setGlobals() {
  local SET_PEER=""
  if [ -z "$OVERRIDE_PEER" ]; then
    SET_PEER=$1
  else
    SET_PEER="${OVERRIDE_PEER}"
  fi
  infoln "Setting env variables for peer${SET_PEER} of ${ORGNAME} org"
  if [ $SET_PEER -eq 0 ]; then
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="${ORGNAMEFC}MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/tlsca/tlsca.${ORGNAME}.${NETWORK}.com-cert.pem
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/users/Admin@${ORGNAME}.${NETWORK}.com/msp
    export CORE_PEER_ADDRESS=peer0.${ORGNAME}.${NETWORK}.com:${PEER0LIS}
  elif [ $SET_PEER -eq 1 ]; then
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="${ORGNAMEFC}MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/tlsca/tlsca.${ORGNAME}.${NETWORK}.com-cert.pem
    export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/users/Admin@${ORGNAME}.${NETWORK}.com/msp
    export CORE_PEER_ADDRESS=peer1.${ORGNAME}.${NETWORK}.com:${PEER1LIS}
  else
    errorln "The peer doesnt exist"
  fi

}


verifyResult() {
  if [ $1 -ne 0 ]; then
    fatalln "$2"
  fi
}
