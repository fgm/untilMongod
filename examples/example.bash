#!/usr/bin/env bash

export MONGO_URL="mongodb://127.0.0.1:27017"
export TIMEOUT=15

# This assumes MongoDB Community has been installed by Homebrew, e.g. on macOS.
brew services stop mongodb-community
(
  sleep 1
  brew services start mongodb-community
) &

go run .. -v -url "${MONGO_URL}" -timeout "${TIMEOUT}"
if [ $? != 0 ]; then
  echo "Failed connecting to server in ${TIMEOUT} seconds. Aborting"
  exit
fi

node app.js
