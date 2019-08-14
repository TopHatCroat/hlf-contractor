#!/bin/bash

print_usage() {
	echo "  Usage: "
	echo "  API_HOST GET|POST|PUT|DELETE endpoint [json_data]"
}

API_HOST=${API_HOST:-''}
VERBOSE=${VERBOSE:-'s'}
PRETTY=${PRETTY:-'python'}

if [[ -z "$API_HOST" ]]; then echo "You must set env var API_HOST"; exit 1; fi

GET(){
    if [[ -z "$1" ]]; then print_usage; exit 1; fi

    SHOW "`
          curl -${VERBOSE} -X GET -H "Content-Type: application/json"  \
          ${API_HOST}$1
         `"
}

DELETE(){
    if [[ -z "$1" ]]; then print_usage; exit 1; fi

    SHOW "`
         curl -${VERBOSE} -X DELETE -H "Content-Type: application/json" \
         ${API_HOST}$1
        `"
}

POST(){
	echo $APi
    if [[ -z "$1" ]]; then print_usage; exit 1; fi
    if [[ -z "$2" ]]; then print_usage; exit 1; fi

    SHOW "`
        curl -${VERBOSE} -X POST -H "Content-Type: application/json"  -d "$2" \
             ${API_HOST}$1
        `"
}


PUT(){
    if [[ -z "$1" ]]; then print_usage; exit 1; fi
    if [[ -z "$2" ]]; then print_usage; exit 1; fi

    SHOW "`
          curl -${VERBOSE} -X PUT -H "Content-Type: application/json"  -d "$2"\
          ${API_HOST}$1
         `"
}

SHOW() {
  PRETTY=${PRETTY:-''}

  if [[ -z "$PRETTY" ]]; then
    echo "$1"
  elif [[ "$PRETTY" == "python" ]]; then
    printf '%s' "$1" | python -c 'import json,sys; print(json.dumps(sys.stdin.read()))'
  else
    echo $1 | jq
  fi
}