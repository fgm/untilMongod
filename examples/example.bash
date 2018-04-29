#!/usr/bin/env bash

export MONGO_URL="mongodb://127.0.0.1:27017"
export TIMEOUT=15

untilMongod -v -url "${MONGO_URL}" -timeout "${TIMEOUT}"
if [ $? ]; then
  echo "Failed connecting to server in ${TIMEOUT} seconds. Aborting"
  exit
fi

node app.js
