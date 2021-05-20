## Overview
This script is an update to our standard example located here:
		 https://github.com/strongdm/strongdm-sdk-go-examples/tree/master/contrib/okta-sync

It is written against Okta golang API version 1.x. 

It partially implements user/Group sync from Okta > SDM, and switches access grants from the user to the Role (Group) level.

---

## How to use this script

1. Set the following environment variables: SDM_API_ACCESS_KEY, SDM_API_SECRET_KEY, OKTA_CLIENT_TOKEN, and OKTA_CLIENT_ORGURL.
2. Edit the matchers.yml file to a) define which groups to sync from Okta to strongDM as [Roles](https://www.strongdm.com/docs/admin-ui-guide/user-management/roles), and b) which strongDM resources users in those groups will receive access to.

  > For example, the sample file in this folder would create strongDM Roles named Support and Engineering. Users would receive access to all `mysql` and `postgres` resources in your organization, respectively.



3. Edit the variable `oktaQueryString` in the .GO file to specify which users to sync to strongDM. (The Okta SDK does not provide a method to retrieve _group members_, unfortunately.)

---

## How it works

This script reads a separate JSON file, matchers.yml, which maps Okta groups to resources in SDM by type or name.
For each Group defined in the YML, an SDM Role will be created, and access to the defined resources will be granted to that Role.
Any user that matches the Okta search filter will be created in SDM (see the "oktaQueryString" definition just below).
If they belong to an Okta Group that is defined in the YML, they will be assigned to the corresponding SDM Role.

The script won't remove any Roles or Users in SDM. However, it will remove any grants for Groups/Roles that are not defined in the YML. It will also add/remove grants for Groups/Roles if you change the mapping in the YML.

An important consideration is that Okta supports multiple group assignment, but strongDM does not. This means that a user with multiple Group memberships will be assigned to the first Group/Role provided by Okta.

---

## How to set up Okta
We recommend that you consider creating SDM-specific Groups in Okta, e.g. sdm-qa, sdm-dev, and assign users in Okta accordingly.
Then define only these groups in the YML, with appropriate resource mapping.

You may wish to modify the oktaQueryString to match only users who belong to the SDM-specific Groups.