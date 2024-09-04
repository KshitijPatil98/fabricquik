#!/bin/bash

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function json_ccp {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${ORGFC}/$2/" \
        -e "s/\${NETWORKNAME}/$3/" \
        -e "s/\${P0PORT}/$4/" \
        -e "s/\${CAPORT}/$5/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        ${PWD}/script_files/ccp_scripts/${ORGNAME}/ccp-template.json > ${PWD}/connection_files/${ORGNAME}/connection-${ORGNAME}.json
}


#Here the reason i am using different names for env variables on left and right is because, when we generate this file it will also replace the variables of left which we dont want. The variables on left should stay as env variables because we the purpose of this file is to generate a template.
ORG=${ORGNAME}
ORGFC=${ORGNAMEFC}
NETWORKNAME=${NETWORK}
P0PORT=${PEER0LIS}
CAPORT=${ORGCAPORT}
PEERPEM=${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/tlsca/tlsca.${ORGNAME}.${NETWORK}.com-cert.pem
CAPEM=${PWD}/organizations/peerOrganizations/${ORGNAME}.${NETWORK}.com/ca/ca.${ORGNAME}.${NETWORK}.com-cert.pem

json_ccp $ORG $ORGFC $NETWORKNAME $P0PORT $CAPORT $PEERPEM $CAPEM

