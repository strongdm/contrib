# Overview

These scripts use a combination of shell script (for AWS operations) and strongDM Python SDK (to register the machine as an SSH Resource). It can be used standalone in AWS EC2, or as a user-data script in an EC2 Autoscaling Launch Template.

## Details

Here is how this solution works:

- the shell script uses Bash to get an AWS machine-level token, and the token to get machine metadata
- it sets those values as environment variables for the Python script to use
- that script uses our SDK to: 
  - register the current machine as an SSH Public Key resource
  - write the new public key to the user's `authorized_keys` file

## Requirements

As follows:

- AWS EC2 environment
- API key pair in the script. (Or refactor to read from an external source like Secrets Manager.)
- Key scope is: `Datasource:create`.
- The AMI in your Launch Template should have Python installed (tested with Python 3 on Ubuntu), or you can edit the shell script to install it.
