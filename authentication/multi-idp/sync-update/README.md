# Synchronize Users from different Identity Providers 

🚧 _This script is written for legacy strongDM organizations, and will not work properly with the new Access Overhaul features. If you are unsure whether your organization is enabled for Access Overhaul, please contact your strongDM Account Manager, or write to support@strongdm.com_

Shim go script for synchronizing users/groups from different Identity Providers (IdP) with strongDM. Current version supports Okta and Google Directory.

Original version: https://github.com/strongdm/strongdm-sdk-go-examples/tree/master/contrib/okta-sync

## How it works
This script reads a separate JSON file, [matchers.yml](matchers.yml), which maps IdP groups to resources in SDM by type or name.

For each Group defined in the YML, an SDM Role will be created, and access to the defined resources will be granted to that Role, using the [filters spec](https://www.strongdm.com/docs/automation/getting-started/filters).

Any active user present in the IdP (Okta, Google Directory) and with associated group(s) listed in [matchers.yml](matchers.yml) will be created in SDM, and assigned to the corresponding SDM Role. 

strongDM only supports 1:1 user-role mappings, when there are multiple groups assigned to a okta user, a composite role is created with multiple sub-roles assigned to it.

The script won't remove any Roles or Users in SDM, unless you use the flags: `-delete-unmatching-roles` or `-delete-unmatching-users`. SDM admin users are ignored during deletion.

## How to use this script
1. Set the environment variables: SDM_API_ACCESS_KEY and SDM_API_SECRET_KEY. 
  * For Okta set OKTA_CLIENT_TOKEN and OKTA_CLIENT_ORGURL.
  * For Google set credentials.json 
2. Edit the [matchers.yml](matchers.yml) file to a) define which groups to sync from the IdP to strongDM as [Roles](https://www.strongdm.com/docs/admin-ui-guide/user-management/roles), and b) which strongDM resources users in those groups will receive access to.

  > For example, the sample file in this folder would create strongDM Roles with access to all `mysql` and `postgres` resources in your organization.

## Sample
Help: 
```
$ go run . -help
  -delete-unmatching-roles
    	delete roles present in SDM but not in matchers.yml
  -delete-unmatching-users
    	delete users present in SDM but not in the selected IdP or assigned to any role in matchers.yml
  -google
    	use Google as IdP
  -json
    	dump a JSON report for debugging
  -log
    	include logging information
  -okta
    	use Okta as IdP
  -plan
    	do not apply changes just show initial queries
````

Okta:
```
$ go run . -okta -delete-unmatching-roles -delete-unmatching-users
5 IdP users, 3 strongDM users in IdP, 3 strongDM roles in Idp
```

Google:
```
$ go run . -google -delete-unmatching-roles -delete-unmatching-users
5 IdP users, 3 strongDM users in IdP, 3 strongDM roles in Idp
```

Considerations:
* For better reporting add `-plan` to your command.

## Google 
1. Enable OAuth Consent: https://console.cloud.google.com/apis/credentials/consent (Internal is OK)
2. Create credentials for a Desktop App: https://console.cloud.google.com/apis/credentials
3. Enable Admin SDK API: https://console.cloud.google.com/apis/api/admin.googleapis.com/overview
4. Administrate Users and OrgUnits: https://admin.google.com/u/2/ac/users
  * A user can only be assigned to one OrgUnit at a time

Considerations:
* Reference: https://developers.google.com/admin-sdk/directory/v1/quickstart/go
* Google uses paths for OrgUnits. For indicating `/` use the flag Root in `matchers.yml`.
