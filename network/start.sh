#!/bin/bash

set -ev

source ./shared.sh

docker-compose -f docker-compose.yaml down --volumes --remove-orphans

docker-compose -f docker-compose.yaml up -d
docker ps -a

# wait for Hyperledger Fabric to start
# incase of errors when running later commands, issue export FABRIC_START_TIMEOUT=<larger number>
export FABRIC_START_TIMEOUT=10
#echo ${FABRIC_START_TIMEOUT}
sleep ${FABRIC_START_TIMEOUT}

# Create the channel
docker exec \
       -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/users/Admin@awesome.agency/msp" \
       anchor.awesome.agency \
       peer channel create -o orderer.foi.org:7050 -c default -f /var/hyperledger/fabric/artifacts/channel.tx \
            --tls true --cafile /etc/hyperledger/fabric/orderer_tls/tlsca.foi.org-cert.pem \
            --outputBlock /var/hyperledger/fabric/artifacts/default.block

# Join anchor.awesome.agency to the channel.
docker exec \
       -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/users/Admin@awesome.agency/msp" \
       anchor.awesome.agency \
       peer channel join -b /var/hyperledger/fabric/artifacts/default.block

# docker exec \
#        -e "CORE_PEER_LOCALMSPID=AwesomeAgencyMSP"\
#        -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/users/Admin@awesome.agency/msp" \
#        peer.awesome.agency \
#        peer channel join -b default.block \

# Join anchor.pharmatic.com to the channel.
docker exec \
      -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/users/Admin@pharmatic.com/msp" \
      peer.pharmatic.com \
      peer channel join -b /var/hyperledger/fabric/artifacts/default.block

# Join anchor.magik.org to the channel.
#docker exec \
#       -e "CORE_PEER_LOCALMSPID=MagikMSP" \
#       -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/msp" \
#       anchor.magik.org \
#       peer channel join -b /var/hyperledger/fabric/artifacts/default.block

# Install user chaincode on Awseome.agency peer and anchor
docker exec \
      -e "CORE_PEER_ADDRESS=anchor.awesome.agency:7051" \
      -e "CORE_PEER_LOCALMSPID=AwesomeAgencyMSP" \
      -e "CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/awesome.agency/peers/anchor.awesome.agency/tls/server.crt" \
      -e "CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/awesome.agency/peers/anchor.awesome.agency/tls/server.key" \
      -e "CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/awesome.agency/peers/anchor.awesome.agency/tls/ca.crt" \
      -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/awesome.agency/users/Admin@awesome.agency/msp" \
      api.awesome.agency \
      peer chaincode install -n users -v 0.0.1 -l golang -p github.com/TopHatCroat/hlf-contractor/chaincode/users \
      --tls --cafile=/etc/hyperledger/fabric/crypto-config/ordererOrganizations/foi.org/orderers/orderer.foi.org/msp/tlscacerts/tlsca.foi.org-cert.pem

# Install user chaincode on Awseome.agency peer and anchor
# docker exec \
#       -e "CORE_PEER_ADDRESS=peer.awesome.agency:7051" \
#       -e "CORE_PEER_LOCALMSPID=AwesomeAgencyMSP" \
#       -e "CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/awesome.agency/peers/peer.awesome.agency/tls/server.crt" \
#       -e "CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/awesome.agency/peers/peer.awesome.agency/tls/server.key" \
#       -e "CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/awesome.agency/peers/peer.awesome.agency/tls/ca.crt" \
#       -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/awesome.agency/users/Admin@awesome.agency/msp" \
#       api.awesome.agency \
#       peer chaincode install -n users -v 0.0.1 -l golang -p github.com/TopHatCroat/hlf-contractor/chaincode/users \
#       --tls --cafile=/etc/hyperledger/fabric/crypto-config/ordererOrganizations/foi.org/orderers/orderer.foi.org/msp/tlscacerts/tlsca.foi.org-cert.pem

# Install charger chaincode on pharmatic.com
docker exec \
      -e "CORE_PEER_ADDRESS=peer.pharmatic.com:7051" \
      -e "CORE_PEER_LOCALMSPID=PharmaticMSP" \
      -e "CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/pharmatic.com/peers/peer.pharmatic.com/tls/server.crt" \
      -e "CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/pharmatic.com/peers/peer.pharmatic.com/tls/server.key" \
      -e "CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/pharmatic.com/peers/peer.pharmatic.com/tls/ca.crt" \
      -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/pharmatic.com/users/Admin@pharmatic.com/msp" \
      api.awesome.agency \
      peer chaincode install -n charger -v 0.0.1 -l golang -p github.com/TopHatCroat/hlf-contractor/chaincode/charger \
      --tls --cafile=/etc/hyperledger/fabric/crypto-config/ordererOrganizations/foi.org/orderers/orderer.foi.org/msp/tlscacerts/tlsca.foi.org-cert.pem

# Install charger chaincode on pharmatic.com
# docker exec \
#       -e "CORE_PEER_ADDRESS=peer.magik.org:7051" \
#       -e "CORE_PEER_LOCALMSPID=MagikMSP" \
#       -e "CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/magik.org/peers/peer.magik.org/tls/server.crt" \
#       -e "CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/magik.org/peers/peer.magik.org/tls/server.key" \
#       -e "CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/crypto-config/peerOrganizations/magik.org/peers/peer.magik.org/tls/ca.crt" \
#       -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/magik.org/users/Admin@magik.org/msp" \
#       api.awesome.agency \
#       peer chaincode install -n charger -v 0.0.1 -l golang -p github.com/TopHatCroat/hlf-contractor/chaincode/charger \
#       --tls --cafile=/etc/hyperledger/fabric/crypto-config/ordererOrganizations/foi.org/orderers/orderer.foi.org/msp/tlscacerts/tlsca.foi.org-cert.pem

