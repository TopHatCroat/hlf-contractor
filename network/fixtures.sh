#!/bin/bash

# If using TLS and running from localhost then before running this script run:
#   `sudo echo "127.0.0.1 ca.awesome.agency" >> /etc/hosts`

set -e

FABRIC_CA_CLIENT_HOME=./ca-client-home
mkdir -p $FABRIC_CA_CLIENT_HOME
rm -rf $FABRIC_CA_CLIENT_HOME/msp

fabric-ca-client enroll --home $FABRIC_CA_CLIENT_HOME --url https://admin:adminpw@ca.awesome.agency:7054

echo "Admin enrolled"

for user in "user1" "user2"; do
    username=$user@mail.com

    fabric-ca-client register --id.name $username --id.secret=password \
                    --home $FABRIC_CA_CLIENT_HOME \
                    --id.attrs '"hf.Registrar.Roles=client","rights=consume:ecert"' \
                    --url https://ca.awesome.agency:7054 \

    echo "User $user registered"

    fabric-ca-client enroll --home $FABRIC_CA_CLIENT_HOME \
                            --url https://$username:password@ca.awesome.agency:7054 \
                            --mspdir $username

    echo "User $user enrolled"
done

echo
echo "Users created:"
fabric-ca-client identity list --home $FABRIC_CA_CLIENT_HOME \
                               --url https://admin:adminpw@ca.awesome.agency:7054
