#!/bin/bash

API_HOST=localhost:8000

source ../_common.sh

POST /login '{ "email": "username1@mail.com", "password": "asdf" }'