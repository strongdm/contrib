#!/bin/bash
echo "${SSH_PUB_KEY}" | sudo tee -a /etc/ssh/sdm_ca.pub
echo "TrustedUserCAKeys /etc/ssh/sdm_ca.pub" | sudo tee -a /etc/ssh/sshd_config
sudo systemctl restart ssh