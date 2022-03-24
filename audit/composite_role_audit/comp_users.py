#!/usr/bin/env python
import strongdm, os

api_access_key = "key"
api_secret_key = "secret"

client = strongdm.Client(api_access_key, api_secret_key)

def list_roles():
    comp_roles = list(client.roles.list("composite:true"))
    for c in comp_roles:
        print("Comp Role Name: **" + c.name + "** has these users:")
        # are any users attached to the comp role?
        comp_users = list(client.account_attachments.list('role_id:{}'.format(c.id)))
        if (len(comp_users) > 0):
            for u in comp_users:
                user = client.accounts.get(u.account_id)
                print("    " + user.account.email)
        else:
            print("    No direct members!")
        print("and these sub-roles:")
        role_attachments = list(client.role_attachments.list('composite_role_id:{}'.format(c.id)))
        if (len(role_attachments) > 0):
            for ra in role_attachments:
                role = client.roles.get(ra.attached_role_id)
                print("  Role: *" + role.role.name +  "*, whose members are:")
                # are any users attached to the sub-role?
                role_users = list(client.account_attachments.list('role_id:{}'.format(ra.attached_role_id)))
                if (len(role_users) > 0):
                    for u in role_users:
                        user = client.accounts.get(u.account_id)
                        print("    " + user.account.email)
                else:
                    print("    No members!")
                print("  and linked resources are:")
                grants = list(client.role_grants.list('role_id:{}'.format(ra.attached_role_id)))
                if (len(grants) > 0):
                    for g in grants:
                        res = client.resources.get(g.resource_id)
                        print("    Resource: ", res.resource.name)
                else:
                    print("    No resources!")
                    
        else:
            print("  No sub-roles!")
        print("=================")

def main():
    list_roles()

main()