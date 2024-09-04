export ${ORGNAMEC}_CA=${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/tlsca/tlsca.${ORGNAME}.${NETWORK}.com-cert.pem
export ${ORGNAMEC}_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/users/Admin@${ORGNAME}.${NETWORK}.com/tls/server.crt
export ${ORGNAMEC}_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/${ORGNAME}Organizations/${ORGNAME}.${NETWORK}.com/users/Admin@${ORGNAME}.${NETWORK}.com/tls/server.key

#Here remember this script runs from local terminal. So it wont know ${ORGNAME}0.${NETWORK}.com by its name. If you want you can add an entry in /etc/hosts file of my computer. Instead i will use localhost

osnadmin channel join --channelID ${CHANNELNAME} --config-block ${PWD}/channel_files/channel_artifacts/${CHANNELNAME}/${CHANNELNAME}.block -o orderer0.${NETWORK}.com:${ORDERER0ADM} --ca-file "$${ORGNAMEC}_CA" --client-cert "$${ORGNAMEC}_ADMIN_TLS_SIGN_CERT" --client-key "$${ORGNAMEC}_ADMIN_TLS_PRIVATE_KEY"

osnadmin channel join --channelID ${CHANNELNAME} --config-block ${PWD}/channel_files/channel_artifacts/${CHANNELNAME}/${CHANNELNAME}.block -o orderer1.${NETWORK}.com:${ORDERER1ADM} --ca-file "$${ORGNAMEC}_CA" --client-cert "$${ORGNAMEC}_ADMIN_TLS_SIGN_CERT" --client-key "$${ORGNAMEC}_ADMIN_TLS_PRIVATE_KEY"

osnadmin channel join --channelID ${CHANNELNAME} --config-block ${PWD}/channel_files/channel_artifacts/${CHANNELNAME}/${CHANNELNAME}.block -o orderer2.${NETWORK}.com:${ORDERER2ADM} --ca-file "$${ORGNAMEC}_CA" --client-cert "$${ORGNAMEC}_ADMIN_TLS_SIGN_CERT" --client-key "$${ORGNAMEC}_ADMIN_TLS_PRIVATE_KEY"
