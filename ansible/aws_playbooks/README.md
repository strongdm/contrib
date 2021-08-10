# Self register AWS Ansible Playbooks

## AWS SDM Gateway

Within the playbook there is a vars section you'll need to update within the AWS task and down in the scripts task. Some of the information you'll need to pull from AWS. You can find all EC2 vars examples [here](https://docs.ansible.com/ansible/latest/collections/amazon/aws/ec2_module.html)

- [Self Registering SDM AWS Gateway Playbook](aws_self_register_playbooks/aws-self-register-gateway.yml)

Inside the script you'll need to add your SDM Admin Token.

- [Ansible SDM Gateway Self Register Script](aws_self_register_playbooks/scripts/sdm-gatewayadd.sh)

## AWS SSH Server

Within the playbook there is a vars section you'll need to update within the AWS task and down in the scripts task. Some of the information you'll need to pull from AWS. You can find all EC2 vars examples [here](https://docs.ansible.com/ansible/latest/collections/amazon/aws/ec2_module.html)

- [Self Registering SDM SSH Resource Playbook](aws_self_register_playbooks/aws-self-register-ssh.yml)

Inside the script you'll need to add your SDM Admin Token.

- [Ansible SDM SSH Self Register Script](aws_self_register_playbooks/scripts/sdm-sshadd.sh)

# Single Ansible Playbooks

## SDM Gateway Install

This playbook will run on any host within the inventory file. I've built a full playbook without the need of a script. To target a specific group changes the `hosts:` parameter. It will auto register any AWS machine with a public address.

- [Self Registering SDM Gateway Playbook](playbooks/sdm_gateway_install.yml)

_Example: `ansible-playbook sdm_gateway_install.yml -i sdm-gateways --extra-vars 'SDM_ADMIN_TOKEN={{string for sdm token}}'`_

## SDM Relay Install

This playbook will run on any host within the inventory file. To target a specific group changes the `hosts:` parameter.

- [Self Registering SDM Relay Playbook](playbooks/sdm_relay_install.yml)

_Example: `ansible-playbook sdm_relay_install.yml -i sdm-relays --extra-vars 'SDM_ADMIN_TOKEN={{string for sdm token}}'`_

## SDM SSH Public Cert Install

This playbook will run on any host within the inventory file. To target a specific group changes the `hosts:` parameter. You'll need to pass `SDM_PUB_CA` using `--extra-vars` to append in the public CA. 

_Example: `ansible-playbook sdm_pub_cert_ssh_install.yml -i ssh-servers --extra-vars 'SDM_ADMIN_TOKEN={{string for sdm token}} SDM_PUB_CA={{string for sdm ca}}'`_

- [Self Registering SDM Public SSH Cert Playbook](playbooks/sdm_pub_cert_ssh_install.yml)

## SDM SSH Install

This playbook will run on any host within the inventory file. To target a specific group changes the `hosts:` parameter.

- [Self Registering SDM SSH Playbook](playbooks/sdm_ssh_install.yml)

_Example: `ansible-playbook sdm_ssh_install.yml -i ssh-servers --extra-vars 'SDM_ADMIN_TOKEN={{string for sdm token}}'`_