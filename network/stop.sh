#!/bin/bash

set -ev

# Shut down the Docker containers that might be currently running.
docker-compose -f docker-compose.yaml stop
