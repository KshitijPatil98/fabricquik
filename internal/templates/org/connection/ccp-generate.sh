#!/bin/bash

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function json_ccp {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        ${PWD}/${PROJECT_NAME}_script_files/ccp_scripts/${ORGNAME_ALL_LOWER}/ccp-template.json
}


ORG=${ORGNAME_FIRST_UPPER}
P0PORT=${PEER0LIS}
CAPORT=${ORGCAPORT}
PEERPEM=${PWD}/organizations/peerOrganizations/${ORGNAME_ALL_LOWER}.${PROJECT_NAME}.com/tlsca/tlsca.${ORGNAME_ALL_LOWER}.${PROJECT_NAME}.com-cert.pem
CAPEM=${PWD}/organizations/peerOrganizations/${ORGNAME_ALL_LOWER}.${PROJECT_NAME}.com/ca/ca.${ORGNAME_ALL_LOWER}.${PROJECT_NAME}.com-cert.pem

echo "$(json_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > ${PWD}/organizations/peerOrganizations/${ORGNAME_ALL_LOWER}.${PROJECT_NAME}.com/connection-${ORGNAME_ALL_LOWER}.json

