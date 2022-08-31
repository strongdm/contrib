#!/usr/bin/env python3

import strongdm, time, os

access_key="iL6etBGqRbyW2FeRFtqIUz7OmS62fdyQtpvhSR2W68EYo6HvyyvR09VySEsTeJUi"
secret_key="FyJhqwLo7W234aN+LR/bj8um4LZ1w8l217+UAMIWfLV3RiV3ddRGEQ=="

# Get AWS environment vars
instance_ip = os.getenv("INSTANCE_IP")
instance_id = os.getenv("INSTANCE_ID")
ssh_user = os.getenv("TARGET_USER")

# Create SDM client
client = strongdm.Client(access_key, secret_key)

# Create the SSH resource
try:
   server = strongdm.SSH(hostname=instance_ip,name=instance_id,username=ssh_user,port=22)
   ssh = client.resources.create(server)
except Exception as ex:
   print(ex)

# Append the new public key to authorized_keys
try:
   with open("/home/{}/.ssh/authorized_keys".format(ssh_user),'a',) as f:
      f.write(ssh.resource.public_key)
   f.close()
except Exception as ex:
   print(ex)


