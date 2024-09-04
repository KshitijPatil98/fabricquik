set -x

ORGNAME=$1
TYPE=$2
CHANNEL=$3
NETWORK=$4

. ${PWD}/script_files/channel/${ORGNAME}/utils.sh

infoln "STOPPING ${ORGNAME} organization"
#-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

infoln "Lets first stop the docker container for orgca of ${ORGNAME} organization"

docker-compose -f ${PWD}/compose_files/orgca/${ORGNAME}/orgca.yaml down --volumes 


infoln "Let us now delete all the generated artifacts for orgca and org artifacts of all the components of ${ORGNAME} organization"

sudo rm -rf ${PWD}/organizations/orgca_certs/${ORGNAME}/msp  ${PWD}/organizations/orgca_certs/${ORGNAME}/ca-cert.pem ${PWD}/organizations/orgca_certs/${ORGNAME}/fabric-ca-server.db ${PWD}/organizations/orgca_certs/${ORGNAME}/IssuerPublicKey ${PWD}/organizations/orgca_certs/${ORGNAME}/IssuerRevocationPublicKey  

sudo rm -rf ${PWD}/organizations/orgca_certs/${ORGNAME}/tls/cert.pem ${PWD}/organizations/orgca_certs/${ORGNAME}/tls/key.pem
  
sudo rm -rf ${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/ca  ${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/msp ${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/tlsca ${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/users ${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/fabric-ca-client-config.yaml 
 
sudo rm -rf ${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/peers/peer0.${ORGNAME}.${NETWORK}.com/msp
sudo rm -rf ${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/peers/peer1.${ORGNAME}.${NETWORK}.com/msp


#------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

infoln "Lets us now stop the docker container for tlsca of ${ORGNAME} organization"

docker-compose -f ${PWD}/compose_files/tlsca/${ORGNAME}/tlsca.yaml down --volumes 

infoln "Let us now delete all the generated artifacts for tlsca and tls artifacts of all the components of the ${ORGNAME} organization"

#Removing the tlsca related artifacts created for ${ORGNAME} and then removing tls certificates of ${ORGNAME}
sudo rm -rf ${PWD}/organizations/tlsca_certs/${ORGNAME}/msp ${PWD}/organizations/tlsca_certs/${ORGNAME}/tls-cert.pem ${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem ${PWD}/organizations/tlsca_certs/${ORGNAME}/IssuerPublicKey ${PWD}/organizations/tlsca_certs/${ORGNAME}/IssuerRevocationPublicKey ${PWD}/organizations/tlsca_certs/${ORGNAME}/fabric-ca-server.db ${PWD}/organizations/tlsca_certs/${ORGNAME}/users/tlsca_admin/msp ${PWD}/organizations/tlsca_certs/${ORGNAME}/users/tlsca_admin/fabric-ca-client-config.yaml
 
 sudo rm -rf ${PWD}/organizations/tlsca_certs/${ORGNAME}/users/rca_admin


if [ "$TYPE" == "org" ] || [ "$TYPE" == "owner" ]; then
 infoln "Finally let us delete the org folder of ${ORGNAME} organization."

 sudo rm -rf ${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com
fi
#---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

 infoln "Let us now stop peer and couchdb containers of ${ORGNAME} organization"

 docker-compose -f ${PWD}/compose_files/peer_couch_orderer/${ORGNAME}/peer.yaml -f ${PWD}/compose_files/peer_couch_orderer/${ORGNAME}/couch.yaml down --volumes 

 infoln "Let us now stop cli container of ${ORGNAME} organization"
 docker-compose -f ${PWD}/compose_files/peer_couch_orderer/${ORGNAME}/cli.yaml  down --volumes 

#----------------------------------------------------------------------------------------------------------------------------------------------------------
 
 #When the peers are stopped all the channel blocks become invalid because channel block contains certs of peers. When we stop , discard and   remove old container and then start new containers, new certs are generated and  the old certificates are invalid and hence old blocks are invalid as it contain those certificates
 infoln "Lets remove the channel block created as well" 
 sudo rm -rf ${PWD}/channel_files/channel_artifacts/${CHANNEL}/*


if [ "$TYPE" == "orderer" ]; then
 infoln "Finally let us delete the orderer folder of orderer organization."

 sudo rm -rf ${PWD}/organizations/ordererOrganizations/orderer.${NETWORK}.com
fi
#---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

 infoln "Let us now stop orderer containers of orderer organization"

 docker-compose -f ${PWD}/compose_files/peer_couch_orderer/orderer/orderer.yaml down --volumes 

 #--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

successln "All the artifacts are cleaned. Also all the containers are stopped. ${ORGNAME}  IS SUCCESSFULLY STOPPED" 

{ set +x; } 2>/dev/null

