# Synchronize Okta Users with Multiple Groups
Shim go script for synchronizing Okta users/groups with strongDM.

Codebase guidelines from: https://github.com/strongdm/strongdm-sdk-go-examples/tree/master/contrib/okta-sync

## How it works
This script reads a separate JSON file, [matchers.yml](matchers.yml), which maps Okta groups to resources in SDM by type or name.

For each Group defined in the YML, an SDM Role will be created, and access to the defined resources will be granted to that Role, using the [filters spec](https://www.strongdm.com/docs/automation/getting-started/filters).

Any user that matches the Okta search filter and has groups associated listed in [matchers.yml](matchers.yml) will be created in SDM, and assigned to the corresponding SDM Role. 

strongDM only supports 1:1 user-role mappings, when there are multiple groups assigned to a okta user, a composite role is created with multiple sub-roles assigned to it.

The script won't remove any Roles or Users in SDM, unless you use the flags: `-delete-roles-not-in-okta` or `-delete-users-not-in-okta`.

## How to use this script
1. Set the following environment variables: SDM_API_ACCESS_KEY, SDM_API_SECRET_KEY, OKTA_CLIENT_TOKEN, and OKTA_CLIENT_ORGURL.
2. Edit the [matchers.yml](matchers.yml) file to a) define which groups to sync from Okta to strongDM as [Roles](https://www.strongdm.com/docs/admin-ui-guide/user-management/roles), and b) which strongDM resources users in those groups will receive access to.

  > For example, the sample file in this folder would create strongDM Roles with access to all `mysql` and `postgres` resources in your organization.

## Sample
```
$ go run . -delete-roles-not-in-okta -delete-users-not-in-okta
5 Okta users, 3 strongDM users in okta, 3 strongDM roles in okta
```

Considerations:
* For better reporting add `-plan` to your command.
* When using `-delete-users-not-in-okta​` remember to add your SDM admin emails to Okta, otherwise you could remove the account administrators.
