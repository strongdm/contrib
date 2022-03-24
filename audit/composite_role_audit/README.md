# List Composite Roles

This folder contains a Python script that lists all composite roles, and each of their members, sub-roles, and related resources, if any.

## Requirements
* Python3
* A [strongDM API key pair](https://www.strongdm.com/docs/admin-ui-guide/settings/admin-tokens/api-keys) with the following permissions:
  * Roles: List
  * Datasources: List
  * Grants: Read
  * Accounts: Read

## Usage
* Run `pip install strongdm==1.0.35`
* Run the script e.g. `python comp_users.py`.

## Notes
* Composite roles are deprecated in Python SDK 2.0.0 and later. Please use version 1.0.35 or earlier. 
* Composite roles no longer apply if you have migrated to the latest strongDM Admin UI. If you are not sure if this applies to you, please contact your Customer Success Manger, or email support@strongdm.com.



