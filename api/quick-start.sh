#!/bin/bash

rm -rf ./on && dlv debug --headless --listen=:8008 --api-version=2 --accept-multiclient
