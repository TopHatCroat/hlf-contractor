#!/bin/bash

API_HOST=localhost:8000

source ../_common.sh

POST /register '{ "email": "username1@mail.com", "password": "asdf" }'