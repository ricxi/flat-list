#!/bin/bash

# curl -X POST -H "Content-Type: application/json" -d @register.json http://localhost:8080/v1/user/register
# curl -X POST -H "Content-Type: application/json" -d @login.json http://localhost:8080/v1/user/login

REGISTER_ENDPOINT="http://localhost:9000/v1/user/register"
LOGIN_ENDPOINT="http://localhost:9000/v1/user/login"
REGISTER_JSON="register.json"
LOGIN_JSON="login.json"

post_request() {
   data="$(cat "$1")"
   endpoint="$2"

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
       echo "script usage: $(basename \$0) [-l] [-h] [-a somevalue]" >&2
       exit 1
       ;;       
    esac
done
shift $((OPTIND-1))
