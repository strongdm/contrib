# Convert SDM roles for Kubernetes User Impersonation

This folder contains a Python script that helps with automation of roles for Kubernetes user impersonation.

## Requirements
* Python3
* A [strongDM API key pair](https://www.strongdm.com/docs/admin-ui-guide/settings/admin-tokens/api-keys) with the following permissions:
  * Roles: List, Create
  * Grants: Write
* An existing Role(s) that you wish to convert and/or add to a Composite Role for k8s user impersonation.


## Usage
* Run `pip install -r requirements.txt`
* Run the script with at least one role name: `k8s_auto.py -r CurrentRole`. This will:
  * Create a corresponding Composite Role with the same name, plus `_k8s`.
  * Add the specified Role to that new Composite Role.
* You may optionally include a "mapping role" via the `-m` flag. This is a role which has already been mapped via a YML file to some RBAC group in Kubernetes. For example, `k8s_auto.py -r CurrentRoleName -m MapRoleName`.
* When users from the role you passed into the script access a k8s cluster via SDM, the client will pass the user's name and all three Role names to the cluster for authorization and auditing.

## Considerations:
* If the Composite Role already exists, we assume this script has been run before. The script will warn and exit.
