set -x

export FABRIC_CFG_PATH=${PWD}/channel_files/configtx_files/${CHANNELNAME}


configtxgen -profile TwoOrgsApplicationGenesis -outputBlock ${PWD}/channel_files/channel_artifacts/${CHANNELNAME}/${CHANNELNAME}.block -channelID ${CHANNELNAME}

{ set +x; } 2>/dev/null
