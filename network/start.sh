#!/bin/bash

set -ev

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
       -e "CORE_PEER_LOCALMSPID=AwesomeAgencyMSP" \
       -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/users/Admin@awesome.agency/msp" \
       anchor.awesome.agency \
       peer channel create -o orderer.foi.org:7050 -c default -f /var/hyperledger/fabric/artifacts/channel.tx --tls true --cafile /etc/hyperledger/fabric/orderer_tls/tlsca.foi.org-cert.pem

# Join anchor.awesome.agency to the channel.
docker exec \
       -e "CORE_PEER_LOCALMSPID=AwesomeAgencyMSP"\
       -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/users/Admin@awesome.agency/msp" \
       anchor.awesome.agency \
       peer channel join -b default.block

# Join anchor.magik.dev to the channel.
#docker exec \
#       -e "CORE_PEER_LOCALMSPID=MagikMSP" \
#       -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/msp" \
#       anchor.magik.dev \
#       peer channel join -b default.block

# Join anchor.pharmatic.com to the channel.
#docker exec \
#       -e "CORE_PEER_LOCALMSPID=PharmaticMSP" \
#       -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/users/Admin@pharmatic.com/msp" \
#       anchor.pharmatic.com \
#       peer channel join -b default.block
