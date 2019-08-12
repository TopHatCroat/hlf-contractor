#!/bin/bash

# If using TLS and running from localhost then before running this script run:
#   `sudo sh -c 'echo "127.0.0.1 ca.awesome.agency
#                      127.0.0.1 ca.pharmatic.com" >> /etc/hosts'`

set -ev

source ./shared.sh

function register_and_enroll() {
    role=$1
    username=$2
    fabric-ca-client register --id.name $username --id.secret=password \
                              --home $FABRIC_CA_CLIENT_HOME \
                              --id.attrs role=${role}:ecert \
                              --id.type user \
                              --url https://ca.awesome.agency:7054

    echo "User $username registered"

    fabric-ca-client enroll --home $FABRIC_CA_CLIENT_HOME \
                            --url https://$username:password@ca.awesome.agency:7054 \
                            --mspdir ../$USER_MSP_DIR/$username

    mkdir -p $USER_MSP_DIR/$username/admincerts
    cp $USER_MSP_DIR/$username/signcerts/cert.pem $USER_MSP_DIR/$username/admincerts/cert.pem

    echo "User $username enrolled"
}

FABRIC_CA_CLIENT_HOME=./ca-client-home
mkdir -p $FABRIC_CA_CLIENT_HOME
rm -rf $FABRIC_CA_CLIENT_HOME/msp
USER_MSP_DIR=./crypto-config/users
# Create the MSP directory if missing, but make sure it is empty
mkdir -p $USER_MSP_DIR
rm -rf $USER_MSP_DIR/*
fabric-ca-client enroll --home $FABRIC_CA_CLIENT_HOME --url https://admin:adminpw@ca.awesome.agency:7054

echo "Global admin enrolled"
echo "Creating global users..."
# Pairs of ROLE:USERNAME
for it in "admin:admin" "user:user1" "user:user2"; do
    role=$(echo $it | cut -f1 -d:)
    name=$(echo $it | cut -f2 -d:)
    username=$name@awesome.com

    register_and_enroll $role $username
done

echo
echo "Users created:"
fabric-ca-client identity list --home $FABRIC_CA_CLIENT_HOME \
                               --url https://admin:adminpw@ca.awesome.agency:7054


# Instatiate the installed users chaincode
docker exec \
      -e "CORE_PEER_ADDRESS=anchor.awesome.agency:7051" \
      -e "CORE_PEER_LOCALMSPID=AwesomeAgencyMSP" \
      -e "CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/awesome.agency/peers/anchor.awesome.agency/tls/server.crt" \
      -e "CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/awesome.agency/peers/anchor.awesome.agency/tls/server.key" \
      -e "CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/awesome.agency/peers/anchor.awesome.agency/tls/ca.crt" \
      -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/awesome.agency/users/Admin@awesome.agency/msp" \
      api.awesome.agency \
      peer chaincode instantiate -o orderer.foi.org:7050 \
        -C $CHANNEL_NAME -n users -l golang -v 0.0.1 -c '{"Args":["init"]}' -P "OR ('AwesomeAgencyMSP.peer')" \
        --tls --cafile=/etc/hyperledger/fabric/crypto-config/ordererOrganizations/foi.org/orderers/orderer.foi.org/msp/tlscacerts/tlsca.foi.org-cert.pem

# Instantiate charger code
docker exec \
      -e "CORE_PEER_ADDRESS=peer.pharmatic.com:7051" \
      -e "CORE_PEER_LOCALMSPID=PharmaticMSP" \
      -e "CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/pharmatic.com/peers/peer.pharmatic.com/tls/server.crt" \
      -e "CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/pharmatic.com/peers/peer.pharmatic.com/tls/server.key" \
      -e "CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/pharmatic.com/peers/peer.pharmatic.com/tls/ca.crt" \
      -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/pharmatic.com/users/Admin@pharmatic.com/msp" \
      api.awesome.agency \
      peer chaincode instantiate -o orderer.foi.org:7050 \
        -C $CHANNEL_NAME -n charger -l golang -v 0.0.1 -c '{"Args":["2"]}' -P "OR ('PharmaticMSP.peer')" \
        --tls --cafile=/etc/hyperledger/fabric/crypto-config/ordererOrganizations/foi.org/orderers/orderer.foi.org/msp/tlscacerts/tlsca.foi.org-cert.pem

# docker exec \
#       -e "CORE_PEER_ADDRESS=peer.magik.org:7051" \
#       -e "CORE_PEER_LOCALMSPID=MagikMSP" \
#       -e "CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/magik.org/peers/peer.magik.org/tls/server.crt" \
#       -e "CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/magik.org/peers/peer.magik.org/tls/server.key" \
#       -e "CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/magik.org/peers/peer.magik.org/tls/ca.crt" \
#       -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/magik.org/users/Admin@magik.org/msp" \
#       api.awesome.agency \
#       peer chaincode instantiate -o orderer.foi.org:7050 \
#         -C $CHANNEL_NAME -n charger -l golang -v 0.0.1 -c '{"Args":["3"]}' -P "OR ('MagikMSP.peer')" \
#         --tls --cafile=/etc/hyperledger/fabric/crypto-config/ordererOrganizations/foi.org/orderers/orderer.foi.org/msp/tlscacerts/tlsca.foi.org-cert.pem

# Wait for chaincode instantiation to propagate
sleep $FABRIC_WAIT_TIME

docker exec \
      -e "CORE_PEER_ADDRESS=peer.pharmatic.com:7051" \
      -e "CORE_PEER_LOCALMSPID=PharmaticMSP" \
      -e "CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/pharmatic.com/peers/peer.pharmatic.com/tls/server.crt" \
      -e "CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/pharmatic.com/peers/peer.pharmatic.com/tls/server.key" \
      -e "CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/pharmatic.com/peers/peer.pharmatic.com/tls/ca.crt" \
      -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/pharmatic.com/users/Admin@pharmatic.com/msp" \
      api.awesome.agency \
      peer chaincode query -o orderer.foi.org:7050 -n charger \
        -C $CHANNEL_NAME -c '{"Args":["QueryAll"]}' \
        --tls --cafile=/etc/hyperledger/fabric/crypto-config/ordererOrganizations/foi.org/orderers/orderer.foi.org/msp/tlscacerts/tlsca.foi.org-cert.pem

