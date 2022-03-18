#!/usr/bin/env python3

import strongdm, time, os

access_key="key"
secret_key="secret"

# Get AWS environment vars
instance_ip = os.getenv("INSTANCE_IP")
instance_id = os.getenv("INSTANCE_ID")
ssh_user = os.getenv("TARGET_USER")

# Create SDM client
client = strongdm.Client(access_key, secret_key)

# Create the SSH resource
server = strongdm.SSH(hostname=instance_ip,name=instance_id,username=ssh_user,port=22)
ssh = client.resources.create(server)
# Append the new public key to authorized_keys
with open("/home/{}/.ssh/authorized_keys".format(ssh_user),'a',) as f:
   f.write(ssh.resource.public_key)
f.close()

# Create gateway
gateway = strongdm.Gateway(
    name="name-us-west-2------{}".format(instance_id),
    listen_address="{}:5000".format(instance_ip),
)
resp = client.nodes.create(gateway)

# Write gateway token to file for installation
# Shell script will delete after it's consumed
with open("token.txt",'a',) as f:
   f.write(resp.token)
f.close()

