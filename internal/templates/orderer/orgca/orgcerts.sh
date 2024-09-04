set -x

#Here we are not making directory because it will already be created during tls certs generation
#mkdir -p organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/

:'The below is equivalent of tlsca_admin. Just that here we dont really store it inside a folder called same orgca_admin but its the same function. The orgcaadmin certificates will be stored in
msp folder in ${ORGNAME}.${NETWORK}.com. In case of tlsca it used to get stored inside tlsca folder. 
'
export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/


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
    OrganizationalUnitIdentifier: orderer' > "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/msp/config.yaml"
    
    
# Copy ${ORGNAME} org's tls CA cert to ${ORGNAME} org's /msp/tlscacerts directory (for use in the channel MSP definition)
mkdir -p "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/msp/tlscacerts"
cp "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/msp/tlscacerts/ca.crt"

# Copy ${ORGNAME} org's tls CA cert to ${ORGNAME} org's /tlsca directory (for use by clients)
mkdir -p "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/tlsca"
cp "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/tlsca/tlsca.${ORGNAME}.${NETWORK}.com-cert.pem"

# Copy ${ORGNAME} org's CA cert to ${ORGNAME} org's /ca directory (for use by clients)
mkdir -p "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/ca"
cp "${PWD}/organizations/orgca_certs/${ORGNAME}/ca-cert.pem" "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/ca/ca.${ORGNAME}.${NETWORK}.com-cert.pem"



#Here we are not specifying msp of bootstrap admin using msp dir flag because be default it is set to point to a msp folder in in fabric-ca-client home directory.  

fabric-ca-client register --caname orgca-${ORGNAME} --id.name orderer0 --id.secret orderer0pw --id.type ${ORGNAME} --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem"

fabric-ca-client register --caname orgca-${ORGNAME} --id.name orderer1 --id.secret orderer1pw --id.type ${ORGNAME} --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem"

fabric-ca-client register --caname orgca-${ORGNAME} --id.name orderer2 --id.secret orderer2pw --id.type ${ORGNAME} --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem"
 
fabric-ca-client register --caname orgca-${ORGNAME} --id.name ${ORGNAME}orgadmin --id.secret ${ORGNAME}orgadminpw --id.type admin --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem"




fabric-ca-client enroll -u https://orderer0:orderer0pw@orgca-${ORGNAME}:${ORGCAPORT} --caname orgca-${ORGNAME} -M "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer0.${NETWORK}.com/msp" --csr.hosts orderer0.${NETWORK}.com  --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" 
cp "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/msp/config.yaml" "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer0.${NETWORK}.com/msp/config.yaml"

fabric-ca-client enroll -u https://orderer1:orderer1pw@orgca-${ORGNAME}:${ORGCAPORT} --caname orgca-${ORGNAME} -M "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer1.${NETWORK}.com/msp" --csr.hosts orderer1.${NETWORK}.com  --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" 
cp "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/msp/config.yaml" "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer1.${NETWORK}.com/msp/config.yaml"

fabric-ca-client enroll -u https://orderer2:orderer2pw@orgca-${ORGNAME}:${ORGCAPORT} --caname orgca-${ORGNAME} -M "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer2.${NETWORK}.com/msp" --csr.hosts orderer2.${NETWORK}.com  --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" 
cp "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/msp/config.yaml" "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer2.${NETWORK}.com/msp/config.yaml"


fabric-ca-client enroll -u https://${ORGNAME}orgadmin:${ORGNAME}orgadminpw@orgca-${ORGNAME}:${ORGCAPORT} --caname orgca-${ORGNAME} -M "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/users/Admin@${ORGNAME}.${NETWORK}.com/msp" --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" 
cp "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/msp/config.yaml" "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/users/Admin@${ORGNAME}.${NETWORK}.com/msp/config.yaml"

{ set +x; } 2>/dev/null

