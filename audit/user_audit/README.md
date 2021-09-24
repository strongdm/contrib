# List all Users with Roles and Resources

This folder contains a Python script that lists all users, with their Role(s) and all resources granted by those Roles.

## Requirements
* Python3
* A [strongDM API key pair](https://www.strongdm.com/docs/admin-ui-guide/settings/admin-tokens/api-keys) with the following permissions:
  * Roles: List
  * Datasources: List
  * Users: List

## Usage
* Run `pip install strongdm`
* Run the script e.g. `python audit_sdm_users`.
* The script will create a CSV file (path defined at top of script).
* It will attempt to open the file using the default application defined on your system. 

## Notes
* This script does not list permanent or temporary direct grants -- though it could be modified to do so!.
