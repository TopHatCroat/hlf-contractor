#!/bin/bash

set -ev

source ./shared.sh
VERSION=1.0.2
CHAINCODE_NAME=charger

# Install new code
docker exec \
      -e "CORE_PEER_ADDRESS=anchor.awesome.agency:7051" \
      -e "CORE_PEER_LOCALMSPID=AwesomeAgencyMSP" \
      -e "CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/awesome.agency/peers/anchor.awesome.agency/tls/server.crt" \
      -e "CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/awesome.agency/peers/anchor.awesome.agency/tls/server.key" \
      -e "CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/awesome.agency/peers/anchor.awesome.agency/tls/ca.crt" \
      -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/awesome.agency/users/Admin@awesome.agency/msp" \
      api.awesome.agency \
      peer chaincode install -n $CHAINCODE_NAME -v $VERSION -l golang -p github.com/TopHatCroat/hlf-contractor/chaincode/$CHAINCODE_NAME \
      --tls --cafile=/etc/hyperledger/fabric/crypto-config/ordererOrganizations/foi.org/orderers/orderer.foi.org/msp/tlscacerts/tlsca.foi.org-cert.pem

docker exec \
      -e "CORE_PEER_ADDRESS=peer.pharmatic.com:7051" \
      -e "CORE_PEER_LOCALMSPID=PharmaticMSP" \
      -e "CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/pharmatic.com/peers/peer.pharmatic.com/tls/server.crt" \
      -e "CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/pharmatic.com/peers/peer.pharmatic.com/tls/server.key" \
      -e "CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/pharmatic.com/peers/peer.pharmatic.com/tls/ca.crt" \
      -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/pharmatic.com/users/Admin@pharmatic.com/msp" \
      api.awesome.agency \
      peer chaincode install -n $CHAINCODE_NAME -v $VERSION -l golang -p github.com/TopHatCroat/hlf-contractor/chaincode/$CHAINCODE_NAME \
      --tls --cafile=/etc/hyperledger/fabric/crypto-config/ordererOrganizations/foi.org/orderers/orderer.foi.org/msp/tlscacerts/tlsca.foi.org-cert.pem

# Instantiate new code
docker exec \
      -e "CORE_PEER_ADDRESS=anchor.awesome.agency:7051" \
      -e "CORE_PEER_LOCALMSPID=AwesomeAgencyMSP" \
      -e "CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/awesome.agency/peers/anchor.awesome.agency/tls/server.crt" \
      -e "CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/awesome.agency/peers/anchor.awesome.agency/tls/server.key" \
      -e "CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/awesome.agency/peers/anchor.awesome.agency/tls/ca.crt" \
      -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/awesome.agency/users/Admin@awesome.agency/msp" \
      api.awesome.agency \
      peer chaincode instantiate -o orderer.foi.org:7050 \
        -C $CHANNEL_NAME -n $CHAINCODE_NAME -v $VERSION -l golang -c '{"Args":["2"]}' -P "AND ('AwesomeAgencyMSP.peer', 'PharmaticMSP.peer')" \
        --tls --cafile=/etc/hyperledger/fabric/crypto-config/ordererOrganizations/foi.org/orderers/orderer.foi.org/msp/tlscacerts/tlsca.foi.org-cert.pem

docker exec \
      -e "CORE_PEER_ADDRESS=peer.pharmatic.com:7051" \
      -e "CORE_PEER_LOCALMSPID=PharmaticMSP" \
      -e "CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/pharmatic.com/peers/peer.pharmatic.com/tls/server.crt" \
      -e "CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/pharmatic.com/peers/peer.pharmatic.com/tls/server.key" \
      -e "CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/pharmatic.com/peers/peer.pharmatic.com/tls/ca.crt" \
      -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/pharmatic.com/users/Admin@pharmatic.com/msp" \
      api.awesome.agency \
      peer chaincode instantiate -o orderer.foi.org:7050 \
        -C $CHANNEL_NAME -n $CHAINCODE_NAME -v $VERSION -l golang -c '{"Args":["2"]}' -P "AND ('AwesomeAgencyMSP.peer', 'PharmaticMSP.peer')" \
        --tls --cafile=/etc/hyperledger/fabric/crypto-config/ordererOrganizations/foi.org/orderers/orderer.foi.org/msp/tlscacerts/tlsca.foi.org-cert.pem
