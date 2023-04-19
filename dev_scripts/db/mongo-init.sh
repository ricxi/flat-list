#!/bin/bash

set -e

mongosh <<EOF

db = db.getSiblingDB('flatlist')
db.createCollection('users')
db.users.createIndex({ email: 1}, { unique: true})

EOF