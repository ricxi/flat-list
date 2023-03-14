#!/bin/bash

# JSON_PAYLOAD=''

END_POINT="http://localhost:5002/v1/token/activation"
# AUTH_TOKEN=''
# AUTH_HEADER='Authorization:Bearer "'"$AUTH_TOKEN"'"'

# this will be the user id in this example
url_param="fakeuserid"

curl_args=(
    -X POST "${END_POINT}/${url_param}"
    # --data "$payload"
    -H "Content-Type:application/json"
    # -H "$AUTH_HEADER"
)

curl "${curl_args[@]}"
