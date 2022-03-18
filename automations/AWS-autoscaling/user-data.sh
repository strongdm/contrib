#!/bin/bash -xe
	apt update 
	apt install -y unzip jq awscli python3-pip
	curl -J -O -L https://app.strongdm.com/releases/cli/linux
	unzip sdm*
	cp sdm /usr/local/bin
  pip install strongdm
	
  TOKEN=$(curl -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600")
  INSTANCE_IDENTITY="$(curl --silent -H "X-aws-ec2-metadata-token: $TOKEN" -v http://169.254.169.254/latest/dynamic/instance-identity/document)"

  export INSTANCE_ID="$(echo $INSTANCE_IDENTITY | jq -r .instanceId)"
  export INSTANCE_IP="$(echo $INSTANCE_IDENTITY | jq -r .privateIp)"
  export TARGET_USER=ubuntu

# Call python file to create SSH, and write token to authorized_keys
python3 /home/ubuntu/gw.py

  sdm logout
  rm -rf /root/.sdm/*
  rm -rf ~ubuntu/.sdm/*

  unset SDM_VERBOSE

sudo ./sdm install --relay --token=`cat /home/$TARGET_USER/token.txt`
rm -rf /home/$TARGET_USER/token.txt
