#!/bin/bash
set -e
sudo systemctl stop punktulis
sudo cp /usr/local/bin/punktulis-dev /usr/local/bin/punktulis
sudo systemctl start punktulis
