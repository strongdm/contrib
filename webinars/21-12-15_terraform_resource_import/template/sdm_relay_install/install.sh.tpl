#!/bin/bash
echo "${SSH_PUB_KEY}" | sudo tee -a /etc/ssh/sdm_ca.pub
echo "TrustedUserCAKeys /etc/ssh/sdm_ca.pub" | sudo tee -a /etc/ssh/sshd_config
sudo systemctl restart ssh
sudo apt update && sudo apt upgrade -y
sudo apt install unzip -y
curl -J -O -L https://app.strongdm.com/releases/cli/linux && unzip sdmcli* && rm -f sdmcli*
sudo ./sdm install --relay --token="${SDM_RELAY_TOKEN}"