#!/bin/bash
set -e
pushd ../../internal/wordle
sqlc generate
popd
go build
sudo systemctl stop punktulis-dev
sudo mv bot /usr/local/bin/punktulis-dev
sudo systemctl start punktulis-dev
