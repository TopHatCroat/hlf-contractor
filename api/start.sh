#!/bin/bash

if ! [[ -x "$(command -v dlv)" ]]; then
	echo "Installing delve..."
	GO111MODULE=off go get -u github.com/go-delve/delve/cmd/dlv
fi

if ! [[ -x "$(command -v modd)" ]]; then
	echo "Installing modd..."
	GO111MODULE=off go get -u github.com/cortesi/modd/cmd/modd
fi

echo "Monitoring filesystem for changes..."
modd
