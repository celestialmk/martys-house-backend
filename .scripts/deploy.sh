#!/bin/bash
set -e

echo "Deployment started ...showing status information"


# Pull the latest version of the app
git pull git@github.com:celestialmk/martys-house-backend.git main
echo "New changes copied to server !"



#Build statically linked executable
CGO_ENABLED=0 go build


# Reloading Application So New Changes could reflect on website
echo "Reload Pocketbase systemd"
systemctl restart manage_pocketbase


echo "Deployment Finished!"