#!/bin/bash

echo "Running full API integration tests"

./register_user.sh
./login_user.sh
./get_me.sh
./get_all_users.sh.sh
