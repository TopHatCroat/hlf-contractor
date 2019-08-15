#!/bin/bash

print_usage() {
	echo "  Usage: "
	echo "  GET|POST|PUT|DELETE endpoint [json_data]"
	echo
	echo "  Environmental variables:"
	echo "    API_HOST - base host URL"
	echo "    API_AUTHORIZATION - Authorization token"
	echo "    VERBOSE -  Verbose option for curl"
	echo "    PRETTY -  Output format, options: jq, python"
}

API_HOST=${API_HOST:-''}
API_AUTHORIZATION=${API_AUTHORIZATION:-''}
VERBOSE=${VERBOSE:-'s'}
PRETTY=${PRETTY:-'jq'}

if [[ -z "$API_HOST" ]]; then echo "You must set env var API_HOST"; exit 1; fi

GET(){
    if [[ -z "$1" ]]; then print_usage; exit 1; fi

    SHOW "`
    	curl -${VERBOSE} -X GET \
        -H "Content-Type: application/json" \
        -H "Authorization: API_AUTHORIZATION" \
        ${API_HOST}$1
        `"
}

DELETE(){
    if [[ -z "$1" ]]; then print_usage; exit 1; fi

    SHOW "`
    	curl -${VERBOSE} -X DELETE \
        -H "Content-Type: application/json" \
	   	-H "Authorization: API_AUTHORIZATION" \
        ${API_HOST}$1
		`"
}

POST(){
	echo $APi
    if [[ -z "$1" ]]; then print_usage; exit 1; fi
    if [[ -z "$2" ]]; then print_usage; exit 1; fi

    SHOW "`
    	curl -${VERBOSE} -X POST -H "Content-Type: application/json"  -d "$2" \
	   	-H "Authorization: API_AUTHORIZATION" \
        ${API_HOST}$1
        `"
}


PUT(){
    if [[ -z "$1" ]]; then print_usage; exit 1; fi
    if [[ -z "$2" ]]; then print_usage; exit 1; fi

    SHOW "`
    	curl -${VERBOSE} -X PUT \
    	-H "Content-Type: application/json"  -d "$2"\
 	   	-H "Authorization: API_AUTHORIZATION" \
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