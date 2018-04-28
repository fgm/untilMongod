#!/usr/bin/env bash

export MONGO_URL="mongodb://127.0.0.1:27017"
export TIMEOUT=60

if [ $(untilMongod -url $MONGO_URL -timeout $TIMEOUT) ]; then
  echo "Failed connecting to server in {$TIMEOUT} seconds. Aborting"
  exit
fi

node app.js
