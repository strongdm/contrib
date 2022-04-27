# List all Roles and associated Resources

This folder contains a Python script that lists all Roles and the Resources linked to each. It works with the legacy strongDM UI, and under the new Access Overhaul system.

## Requirements
* Python3
* A [strongDM API key pair](https://www.strongdm.com/docs/admin-ui-guide/settings/admin-tokens/api-keys) with the following permissions:
  * Roles: List
  * Datasources: List
  * Grants: Read

## Usage
* Run `pip install strongdm`
* Run the script e.g. `python role_resources.py`.
