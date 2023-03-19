#!/bin/bash
# check that http server is working to create and register a token for a user

# JSON_PAYLOAD=''

END_POINT="http://localhost:5002/v1/token/activation"
# AUTH_TOKEN=''
# AUTH_HEADER='Authorization:Bearer "'"$AUTH_TOKEN"'"'

# this will be the user id in this example
url_param="fakeuserid"

curl_args=(
    -X POST "${END_POINT}/${url_param}"
    # --data "$payload"
    -s
    -H "Content-Type:application/json"
    # -H "$AUTH_HEADER"
)

response=$( curl "${curl_args[@]}" )
echo $response
