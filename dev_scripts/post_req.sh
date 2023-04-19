#!/bin/bash
# Make requests to test user service
cd "$(dirname "$0")"

REGISTER_ENDPOINT="http://localhost:9000/v1/user/register"
LOGIN_ENDPOINT="http://localhost:9000/v1/user/login"
REGISTER_JSON="./json_req_data/register.json"
LOGIN_JSON="./json_req_data/login.json"

post_request() {
   local data="$(cat "$1")"
   local endpoint="$2"

   curl -X POST \
     -v \
      -H "Content-Type: application/json" \
      -d "$data" \
     "$endpoint" 
}

while getopts 'rl' flag; do
    case "$flag" in
     r)
        echo "register"
        post_request "$REGISTER_JSON" "$REGISTER_ENDPOINT" 
        ;;
     l)
        echo "login"
        post_request "$LOGIN_JSON" "$LOGIN_ENDPOINT"
        ;;
     ?)
       echo "script usage: $(basename \$0) [-l] [-r] [-a somevalue]" >&2
       exit 1
       ;;       
    esac
done
shift $((OPTIND-1))

# curl -X POST -H "Content-Type: application/json" -d @register.json http://localhost:8080/v1/user/register
# curl -X POST -H "Content-Type: application/json" -d @login.json http://localhost:8080/v1/user/login