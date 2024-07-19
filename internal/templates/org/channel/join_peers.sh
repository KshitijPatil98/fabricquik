#!/bin/bash

#Here basically we are setting env variables which will be used in rest of the code. So here the channel name will come from CHANNELNAME variable which we will pass during sed. But once that is done below everywhere we will use the value of the env variable set. Hence you see CHANNEL_NAME used in rest of the code.

: ${CHANNEL_NAME:="${CHANNELNAME}"}
: ${DELAY:="3"}
: ${MAX_RETRY:="5"}
: ${VERBOSE:="false"}
: ${CONTAINER_CLI:="docker"}

BLOCKFILE="${PWD}/channel_files/channel_artifacts/${CHANNEL_NAME}/${CHANNEL_NAME}.block"
#export FABRIC_CFG_PATH=${PWD}/../config/

# imports  
. ${PWD}/script_files/channel_scripts/${ORGNAME}/envVar.sh
. ${PWD}/script_files/channel_scripts/${ORGNAME}/utils.sh



# joinChannel ORG
joinChannel() {
  PEER=$1
  setGlobals $PEER
	local rc=1
	local COUNTER=1
	## Sometimes Join takes time, hence retry
	while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ] ; do
    sleep $DELAY
    set -x
    peer channel join -b $BLOCKFILE >&log.txt
    res=$?
    { set +x; } 2>/dev/null
		let rc=$res
		COUNTER=$(expr $COUNTER + 1)
	done
	cat log.txt
	verifyResult $res "After $MAX_RETRY attempts, peer${ORG} of ${ORGNAME} has failed to join channel '$CHANNEL_NAME' "
}

setAnchorPeer() {
  PEER=$1
  docker exec cli_${ORGNAME} ./script_files/channel_scripts/${ORGNAME}/${CHANNEL_NAME}/setAnchorPeer.sh $PEER $CHANNEL_NAME
}


## Join all the peers to the channel
infoln "Joining peer0 of ${ORGNAME} to the channel..."
joinChannel 0
infoln "Joining peer1 of ${ORGNAME} to the channel..."
joinChannel 1


#Here we will use orderer0 of orderer org to send a anchor peer update. For supplier we will use orderer1
## Set the anchor peers for each org in the channel
infoln "Setting peer0 as anchor peer for ${ORGNAME}..."
setAnchorPeer 0
infoln "Anchor peer set for ${ORGNAME}..."

#successln "Channel '$CHANNEL_NAME' joined by peer0 and peer1"
