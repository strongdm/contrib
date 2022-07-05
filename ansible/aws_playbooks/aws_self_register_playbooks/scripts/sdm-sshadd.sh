#!/bin/bash
 export SDM_ADMIN_TOKEN= {{ SDM_ADMIN_TOKEN }}
 apt update
 apt install -y unzip
 curl -o sdm.zip -L https://app.strongdm.com/releases/cli/linux
 unzip sdm.zip
 ./sdm admin ssh add \
   -p `curl http://169.254.169.254/latest/meta-data/instance-id` \
   $USERNAME@`curl http://169.254.169.254/latest/meta-data/public-hostname` \
   | tee -a "/home/$USERNAME/.ssh/authorized_keys"
 rm sdm.zip
