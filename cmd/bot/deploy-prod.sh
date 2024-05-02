#!/bin/bash
set -e
pushd ../../internal/wordle
sqlc generate
popd
go build
sudo systemctl stop punktulis
sudo mv bot /usr/local/bin/punktulis
sudo systemctl start punktulis
