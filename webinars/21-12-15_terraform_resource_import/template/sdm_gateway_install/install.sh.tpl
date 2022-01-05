#!/bin/bash
sudo apt install unzip -y
curl -J -O -L https://app.strongdm.com/releases/cli/linux && unzip sdmcli* && rm -f sdmcli*
sudo ./sdm install --relay --token="${SDM_GATEWAY_TOKEN}"