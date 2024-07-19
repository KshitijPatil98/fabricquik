#!/bin/bash


# installChaincode PEER ORG
function installChaincode() {
  PEER=$1
  setGlobalsCLI $PEER
  set -x
  peer lifecycle chaincode install ${CC_SRC_PATH} >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  verifyResult $res "Chaincode installation on peer${PEER} of ${OWNERNAME_ALL_LOWER} has failed"
  successln "Chaincode is installed on peer${PEER} of ${OWNERNAME_ALL_LOWER}"
}

# queryInstalled PEER ORG
function queryInstalled() {
  PEER=$1
  setGlobalsCLI $PEER
  set -x
  peer lifecycle chaincode queryinstalled --output json | jq -r 'try (.installed_chaincodes[].package_id)' | grep ^${PACKAGE_ID}$ >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  verifyResult $res "Query installed on peer${PEER} of ${ORGNAME_ALL_LOWER} has failed"
  successln "Query installed successful on peer${PEER} of ${ORGNAME_ALL_LOWER} on channel "
}

# approveForMyOrg VERSION PEER ORG
function approveForMyOrg() {
  PEER=$1
  setGlobalsCLI $PEER
  set -x
  peer lifecycle chaincode approveformyorg -o ${ORDERER_NODE_NAME}.${PROJECT_NAME}.com:${ORDERER_NODE_LISPORT} --ordererTLSHostnameOverride ${ORDERER_NODE_NAME}.${PROJECT_NAME}.com --tls --cafile "$ORDERER_CA" --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${CC_VERSION} --package-id ${PACKAGE_ID} --sequence ${CC_SEQUENCE} ${INIT_REQUIRED} ${CC_END_POLICY} ${CC_COLL_CONFIG} >&log.txt
  res=$?
  { set +x; } 2>/dev/null
  cat log.txt
  verifyResult $res "Chaincode definition approved on peer${PEER} of ${ORGNAME_ALL_LOWER} on channel '$CHANNEL_NAME' failed"
  successln "Chaincode definition approved on peer${PEER} on channel '$CHANNEL_NAME'"
}

# checkCommitReadiness VERSION PEER ORG
function checkCommitReadiness() {
  PEER=$1
  shift 1
  setGlobalsCLI $PEER
  infoln "Checking the commit readiness of the chaincode definition on peer${PEER} of ${ORGNAME_ALL_LOWER} on channel '$CHANNEL_NAME'..."
  local rc=1
  local COUNTER=1
  # continue to poll
  # we either get a successful response, or reach MAX RETRY
  while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ]; do
    sleep $DELAY
    infoln "Attempting to check the commit readiness of the chaincode definition on peer${PEER} of ${ORGNAME_ALL_LOWER}, Retry after $DELAY seconds."
    set -x
    peer lifecycle chaincode checkcommitreadiness --channelID $CHANNEL_NAME --name ${CC_NAME} --version ${CC_VERSION} --sequence ${CC_SEQUENCE} ${INIT_REQUIRED} ${CC_END_POLICY} ${CC_COLL_CONFIG} --output json >&log.txt
    res=$?
    { set +x; } 2>/dev/null
    let rc=0
    for var in "$@"; do
      grep "$var" log.txt &>/dev/null || let rc=1
    done
    COUNTER=$(expr $COUNTER + 1)
  done
  cat log.txt
  if test $rc -eq 0; then
    infoln "Checking the commit readiness of the chaincode definition successful on peer${PEER} of ${ORGNAME_ALL_LOWER} on channel '$CHANNEL_NAME'"
  else
    fatalln "After $MAX_RETRY attempts, Check commit readiness result on peer${PEER} is INVALID!"
  fi
}


# queryCommitted ORG
function queryCommitted() {
  PEER=$1
  setGlobalsCLI $PEER
  EXPECTED_RESULT="Version: ${CC_VERSION}, Sequence: ${CC_SEQUENCE}, Endorsement Plugin: escc, Validation Plugin: vscc"
  infoln "Querying chaincode definition on peer${PEER} of ${ORGNAME_ALL_LOWER} on channel '$CHANNEL_NAME'..."
  local rc=1
  local COUNTER=1
  # continue to poll
  # we either get a successful response, or reach MAX RETRY
  while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ]; do
    sleep $DELAY
    infoln "Attempting to Query committed status on peer${ORG} of ${ORGNAME_ALL_LOWER}, Retry after $DELAY seconds."
    set -x
    peer lifecycle chaincode querycommitted --channelID $CHANNEL_NAME --name ${CC_NAME} >&log.txt
    res=$?
    { set +x; } 2>/dev/null
    test $res -eq 0 && VALUE=$(cat log.txt | grep -o '^Version: '$CC_VERSION', Sequence: [0-9]*, Endorsement Plugin: escc, Validation Plugin: vscc')
    test "$VALUE" = "$EXPECTED_RESULT" && let rc=0
    COUNTER=$(expr $COUNTER + 1)
  done
  cat log.txt
  if test $rc -eq 0; then
    successln "Query chaincode definition successful on peer${PEER} of ${ORGNAME_ALL_LOWER} on channel '$CHANNEL_NAME'"
  else
    fatalln "After $MAX_RETRY attempts, Query chaincode definition result on peer${PEER} of ${ORGNAME_ALL_LOWER} is INVALID!"
  fi
}


