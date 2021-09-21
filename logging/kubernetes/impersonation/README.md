### Important

This solution is designed to work with multiple access roles. If your strongDM organization has not yet been transitioned to this system, please refer to this [archived folder](https://github.com/strongdm/contrib/tree/main/archive/logging/kubernetes/impersonation) to use the older code for Composite Roles.

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
* Run the script with at two role names. This will copy the users in the "r" role to the "m" role one role name: `k8s_auto.py -r CurrentRole`. For example, `k8s_auto.py -r CurrentRoleName -m MapRoleName`.
* When users from the role you passed into the script access a k8s cluster via SDM, the client will pass the user's name and both Role names to the cluster for authorization and auditing.

