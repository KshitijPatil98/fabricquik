set -x


export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/tlsca_certs/${ORGNAME}/users/tlsca_admin


fabric-ca-client enroll -d -u https://tlscaadmin:tlscaadminpw@tlsca-${ORGNAME}:${TLSCAPORT} --caname tlsca-${ORGNAME} --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" --enrollment.profile tls 


fabric-ca-client register -d --caname tlsca-${ORGNAME} --id.name rcaadmin --id.secret rcaadminpw  --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" --mspdir "${PWD}/organizations/tlsca_certs/${ORGNAME}/users/tlsca_admin/msp"

fabric-ca-client enroll -d -u https://rcaadmin:rcaadminpw@tlsca-${ORGNAME}:${TLSCAPORT} --caname tlsca-${ORGNAME} --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" --enrollment.profile tls --csr.hosts orgca-${ORGNAME} --mspdir "${PWD}/organizations/tlsca_certs/${ORGNAME}/users/rca_admin/msp"


#------------------------------------------REGISTERING AND ENROLLING Orderer0 NODE OF ${ORGNAME} ORGANIZATION-----------

#you can also run this command without mspdir flag as by defualt it looks for msp folder inside fabric ca client directory
fabric-ca-client register -d --caname tlsca-${ORGNAME} --id.name orderer0 --id.secret orderer0pw  --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" --mspdir "${PWD}/organizations/tlsca_certs/${ORGNAME}/users/tlsca_admin/msp"

fabric-ca-client enroll -u https://orderer0:orderer0pw@tlsca-${ORGNAME}:${TLSCAPORT} --caname tlsca-${ORGNAME} -M "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer0.${NETWORK}.com/tls" --enrollment.profile tls --csr.hosts orderer0.${NETWORK}.com --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem"

# Copy the tls CA cert, server cert, server keystore to well known file names in the orderer's tls directory that are referenced by peer startup config
cp "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer0.${NETWORK}.com/tls/tlscacerts/"* "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer0.${NETWORK}.com/tls/ca.crt"
cp "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer0.${NETWORK}.com/tls/signcerts/"* "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer0.${NETWORK}.com/tls/server.crt"
cp "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer0.${NETWORK}.com/tls/keystore/"* "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer0.${NETWORK}.com/tls/server.key"  


#---------------------------------------REGISTERING AND ENROLLING ORDERER1 NODE OF ${ORGNAME} ORGANIZATION----------

#you can also run this command without mspdir flag as by defualt it looks for msp folder inside fabric ca client directory
fabric-ca-client register -d --caname tlsca-${ORGNAME} --id.name orderer1 --id.secret orderer1pw  --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" --mspdir "${PWD}/organizations/tlsca_certs/${ORGNAME}/users/tlsca_admin/msp"

fabric-ca-client enroll -u https://orderer1:orderer1pw@tlsca-${ORGNAME}:${TLSCAPORT} --caname tlsca-${ORGNAME} -M "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer1.${NETWORK}.com/tls" --enrollment.profile tls --csr.hosts orderer1.${NETWORK}.com  --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem"

# Copy the tls CA cert, server cert, server keystore to well known file names in the orderer's tls directory that are referenced by peer startup config
cp "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer1.${NETWORK}.com/tls/tlscacerts/"* "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer1.${NETWORK}.com/tls/ca.crt"
cp "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer1.${NETWORK}.com/tls/signcerts/"* "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer1.${NETWORK}.com/tls/server.crt"
cp "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer1.${NETWORK}.com/tls/keystore/"* "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer1.${NETWORK}.com/tls/server.key"  


#--------------------------------------REGISTERING AND ENROLLING ORDERER2 NODE OF ${ORGNAME} ORGANIZATION---------

#you can also run this command without mspdir flag as by defualt it looks for msp folder inside fabric ca client directory
fabric-ca-client register -d --caname tlsca-${ORGNAME} --id.name orderer2 --id.secret orderer2pw  --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" --mspdir "${PWD}/organizations/tlsca_certs/${ORGNAME}/users/tlsca_admin/msp"

fabric-ca-client enroll -u https://orderer2:orderer2pw@tlsca-${ORGNAME}:${TLSCAPORT} --caname tlsca-${ORGNAME} -M "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer2.${NETWORK}.com/tls" --enrollment.profile tls --csr.hosts orderer2.${NETWORK}.com  --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem"

# Copy the tls CA cert, server cert, server keystore to well known file names in the orderer's tls directory that are referenced by peer startup config
cp "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer2.${NETWORK}.com/tls/tlscacerts/"* "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer2.${NETWORK}.com/tls/ca.crt"
cp "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer2.${NETWORK}.com/tls/signcerts/"* "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer2.${NETWORK}.com/tls/server.crt"
cp "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer2.${NETWORK}.com/tls/keystore/"* "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/${ORGNAME}s/orderer2.${NETWORK}.com/tls/server.key"  


#--------------------------------------REGISTERING AND ENROLLING ADMIN USER OF ${ORGNAME} ORGANIZATION------------------
:'
Now, technically tls certificates are generated only for the server endpoints. It means certificates are only generated for the identities which act as server(example fabric ca server),
peers, orderers etc.
Now we always have 2 admins. One admin is like the bootstrap admin which is basically admin of the ca server. Other is an admin user of the organization which is of type admin. This is like a 
priviledged user of the organization which is responsible for doing some administrative tasks like joining orderer nodes to the channel etc. 
Now, in the latest update you dont need to create a system channel and application channel seperately. You just create a channel block using the configtx file and the tool and then use the
osnadmin cli to join the ordering nodes to the channel. Now, for running this osnadmin  command you need to be an admim with certain privileges. Now you might wonder that, the admin which we are
talking of is already getting enrolled in the orgcerts_orderer.sh file why do we need tls certificates for this orderer admin user. The reason is when we hit the osnadmin command it is targetted to the 
admin port on the same orderer node(same IP) just a different port. It is configured during bootstraping the orderer node. Now the intercation between the osnadmin and the orderer admin server 
endpoint needs to be mutual tls (both client side and server side need to provide tls certificates during the interaction). Now the orderer admin server certificates are nothing but the tls certificates of the orderer node itself. And the admin client certificates will be generated below
'

#you can also run this command without mspdir flag as by defualt it looks for msp folder inside fabric ca client directory
fabric-ca-client register -d --caname tlsca-${ORGNAME} --id.name ${ORGNAME}orgadmin --id.secret ${ORGNAME}orgadminpw  --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem" --mspdir "${PWD}/organizations/tlsca_certs/${ORGNAME}/users/tlsca_admin/msp"

fabric-ca-client enroll -u https://${ORGNAME}orgadmin:${ORGNAME}orgadminpw@tlsca-${ORGNAME}:${TLSCAPORT} --caname tlsca-${ORGNAME} -M "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/users/Admin@${ORGNAME}.${NETWORK}.com/tls" --enrollment.profile tls --csr.hosts ${ORGNAME}orgadmin --tls.certfiles "${PWD}/organizations/tlsca_certs/${ORGNAME}/ca-cert.pem"

# Copy the tls CA cert, server cert, server keystore to well known file names. We are just renaming here.
cp "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/users/Admin@${ORGNAME}.${NETWORK}.com/tls/tlscacerts/"* "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/users/Admin@${ORGNAME}.${NETWORK}.com/tls/ca.crt"
cp "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/users/Admin@${ORGNAME}.${NETWORK}.com/tls/signcerts/"* "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/users/Admin@${ORGNAME}.${NETWORK}.com/tls/server.crt"
cp "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/users/Admin@${ORGNAME}.${NETWORK}.com/tls/keystore/"* "${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/users/Admin@${ORGNAME}.${NETWORK}.com/tls/server.key"  


{ set +x; } 2>/dev/null

