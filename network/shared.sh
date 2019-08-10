#!/bin/sh

PATH=$GOPATH/src/github.com/hyperledger/fabric/build/bin:${PWD}/../bin:${PWD}:$PATH
FABRIC_CFG_PATH=${PWD}
CHANNEL_NAME=default
PHARMATIC_CHANNEL_NAME=pharmatic_default