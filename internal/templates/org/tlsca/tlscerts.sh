#!/bin/bash


set -x


#Here /organizations/xyz and  organizations/xyz is a 
#Make sure you already have a folder named users created and it should contain a folder named tlsca_admin

export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/tlsca_certs/${ORGNAME}/users/tlsca_admin


:'This script will run in my pc and my pc doesnt know what tlsca-${ORGNAME} is. So make an entry in /etc/hosts file and you are all set.
'
fabric-ca-client enroll -d -u https://tlscaadmin:tlscaadminpw@tlsca-${ORGNAME}:${TLSCAPORT} --caname tlsca-${ORGNAME} --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" --enrollment.profile tls 


fabric-ca-client register -d --caname tlsca-${ORGNAME} --id.name rcaadmin --id.secret rcaadminpw  --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" --mspdir "${PWD}/organizations/tlsca_certs/${ORGNAME}/users/tlsca_admin/msp"

fabric-ca-client enroll -d -u https://rcaadmin:rcaadminpw@tlsca-${ORGNAME}:${TLSCAPORT} --caname tlsca-${ORGNAME} --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" --enrollment.profile tls --csr.hosts orgca-${ORGNAME} --mspdir "${PWD}/organizations/tlsca_certs/${ORGNAME}/users/rca_admin/msp"


fabric-ca-client register -d --caname tlsca-${ORGNAME} --id.name peer0 --id.secret peer0pw  --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" --mspdir "${PWD}/organizations/tlsca_certs/${ORGNAME}/users/tlsca_admin/msp"

fabric-ca-client enroll -u https://peer0:peer0pw@tlsca-${ORGNAME}:${TLSCAPORT} --caname tlsca-${ORGNAME} -M "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/peers/peer0.${ORGNAME}.${NETWORK}.com/tls" --enrollment.profile tls --csr.hosts peer0.${ORGNAME}.${NETWORK}.com --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem"

cp "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/peers/peer0.${ORGNAME}.${NETWORK}.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/peers/peer0.${ORGNAME}.${NETWORK}.com/tls/ca.crt"
cp "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/peers/peer0.${ORGNAME}.${NETWORK}.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/peers/peer0.${ORGNAME}.${NETWORK}.com/tls/server.crt"
cp "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/peers/peer0.${ORGNAME}.${NETWORK}.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/peers/peer0.${ORGNAME}.${NETWORK}.com/tls/server.key"  



fabric-ca-client register -d --caname tlsca-${ORGNAME} --id.name peer1 --id.secret peer1pw  --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" --mspdir "${PWD}/organizations/tlsca_certs/${ORGNAME}/users/tlsca_admin/msp"

fabric-ca-client enroll -u https://peer1:peer1pw@tlsca-${ORGNAME}:${TLSCAPORT} --caname tlsca-${ORGNAME} -M "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/peers/peer1.${ORGNAME}.${NETWORK}.com/tls" --enrollment.profile tls --csr.hosts peer1.${ORGNAME}.${NETWORK}.com  --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem"

cp "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/peers/peer1.${ORGNAME}.${NETWORK}.com/tls/tlscacerts/"* "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/peers/peer1.${ORGNAME}.${NETWORK}.com/tls/ca.crt"
cp "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/peers/peer1.${ORGNAME}.${NETWORK}.com/tls/signcerts/"* "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/peers/peer1.${ORGNAME}.${NETWORK}.com/tls/server.crt"
cp "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/peers/peer1.${ORGNAME}.${NETWORK}.com/tls/keystore/"* "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/peers/peer1.${ORGNAME}.${NETWORK}.com/tls/server.key"  




{ set +x; } 2>/dev/null

