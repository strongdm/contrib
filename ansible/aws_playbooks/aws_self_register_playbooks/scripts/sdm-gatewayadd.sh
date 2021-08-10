#!/bin/bash
apt update && apt upgrade -y
apt install zip curl wget -y
curl -J -O -L https://app.strongdm.com/releases/cli/linux
unzip *.zip
export SDM_ADMIN_TOKEN={{ SDM_ADMIN_TOKEN }}
export INSTANCE_HOSTNAME=$(curl http://169.254.169.254/latest/meta-data/public-hostname)
export SDM_RELAY_TOKEN=`./sdm relay create-gateway $INSTANCE_HOSTNAME:5000 0.0.0.0:5000`
unset SDM_ADMIN_TOKEN
./sdm install --relay --token=$SDM_RELAY_TOKEN