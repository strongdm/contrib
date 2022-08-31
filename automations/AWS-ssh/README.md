# Overview

These scripts use a combination of shell script (for AWS operations) and strongDM Python SDK (to register the machine as both Gateway and SSH Resource). It can be used standalone in AWS EC2, or as a user-data script in an EC2 Autoscaling Launch Template.

## Details

Here is how this solution works:

- the shell script uses Bash to get an AWS machine-level token, and the token to get machine metadata
- it sets those values as environment variables for the Python script to use
- that script uses our SDK to:
  - register the current machine as an SSH Public Key resource
  - write the new public key to the user's `authorized_keys` file
  - register the current machine as a Gateway
  - write the new token to a local file
- the shell script reads that file, uses it to install the Gateway, then deletes the file

## Requirements

As follows:

- AWS EC2 environment
- API key pair in the script. (Or refactor to read from an external source like Secrets Manager.)
- Key scope is: `Datasource:create`, and `Relays:create`.
- The AMI in your Launch Template should have Python installed (tested with Python 3 on Ubuntu), or you can edit the shell script to install it.