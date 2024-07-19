set -x

#Here we are not making directory because it will already be created during tls certs generation
#mkdir -p organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/

:'The below is equivalent of tlsca_admin which we had in tlscerts. Just that here we dont really store it inside a folder called same orgca_admin but its the same function. The orgcaadmin certificates will be stored in msp folder in ${ORGNAME}.${NETWORK}.com. In case of tlsca it used to get stored inside tlscaadmin folder. The orgcaadmin certificates are primararily used during registering identities. The register function looks for msp folder inside the directory specified in FABRIC_CA_CLIENT_HOME
'
export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/


#When this enroll happens a msp folder gets generated in the fabric ca client home directory and these are basically certificates of the orgca admin.    
fabric-ca-client enroll -u https://orgcaadmin:orgcaadminpw@orgca-${ORGNAME}:${ORGCAPORT} --caname orgca-${ORGNAME} --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem"

echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/orgca-${ORGNAME}-${ORGCAPORT}-orgca-${ORGNAME}.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/orgca-${ORGNAME}-${ORGCAPORT}-orgca-${ORGNAME}.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/orgca-${ORGNAME}-${ORGCAPORT}-orgca-${ORGNAME}.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/orgca-${ORGNAME}-${ORGCAPORT}-orgca-${ORGNAME}.pem
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/msp/config.yaml"
    
    
# Copy ${ORGNAME} org's tls CA cert to ${ORGNAME} org's /msp/tlscacerts directory (for use in the channel MSP definition)
mkdir -p "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/msp/tlscacerts"
cp "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/msp/tlscacerts/ca.crt"

# Copy ${ORGNAME} org's tls CA cert to ${ORGNAME} org's /tlsca directory (for use by clients)
mkdir -p "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/tlsca"
cp "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/tlsca/tlsca.${ORGNAME}.${NETWORK}.com-cert.pem"

# Copy ${ORGNAME} org's CA cert to ${ORGNAME} org's /ca directory (for use by clients)
mkdir -p "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/ca"
cp "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/ca/ca.${ORGNAME}.${NETWORK}.com-cert.pem"



#Here we are not specifying msp of bootstrap admin using msp dir flag because be default it is set to point to a msp folder in in fabric-ca-client home directory.  

fabric-ca-client register --caname orgca-${ORGNAME} --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem"

fabric-ca-client register --caname orgca-${ORGNAME} --id.name peer1 --id.secret peer1pw --id.type peer --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem"
 
fabric-ca-client register --caname orgca-${ORGNAME} --id.name user1 --id.secret user1pw --id.type client --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem"

fabric-ca-client register --caname orgca-${ORGNAME} --id.name ${ORGNAME}orgadmin --id.secret ${ORGNAME}orgadminpw --id.type admin --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem"




fabric-ca-client enroll -u https://peer0:peer0pw@orgca-${ORGNAME}:${ORGCAPORT} --caname orgca-${ORGNAME} -M "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/peers/peer0.${ORGNAME}.${NETWORK}.com/msp" --csr.hosts peer0.${ORGNAME}.${NETWORK}.com  --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" 
cp "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/peers/peer0.${ORGNAME}.${NETWORK}.com/msp/config.yaml"

fabric-ca-client enroll -u https://peer1:peer1pw@orgca-${ORGNAME}:${ORGCAPORT} --caname orgca-${ORGNAME} -M "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/peers/peer1.${ORGNAME}.${NETWORK}.com/msp" --csr.hosts peer1.${ORGNAME}.${NETWORK}.com  --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" 
cp "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/peers/peer1.${ORGNAME}.${NETWORK}.com/msp/config.yaml"


fabric-ca-client enroll -u https://user1:user1pw@orgca-${ORGNAME}:${ORGCAPORT} --caname orgca-${ORGNAME} -M "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/users/User1@${ORGNAME}.${NETWORK}.com/msp" --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" 
cp "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/users/User1@${ORGNAME}.${NETWORK}.com/msp/config.yaml"

fabric-ca-client enroll -u https://${ORGNAME}orgadmin:${ORGNAME}orgadminpw@orgca-${ORGNAME}:${ORGCAPORT} --caname orgca-${ORGNAME} -M "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/users/Admin@${ORGNAME}.${NETWORK}.com/msp" --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" 
cp "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/users/Admin@${ORGNAME}.${NETWORK}.com/msp/config.yaml"

{ set +x; } 2>/dev/null

