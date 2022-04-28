#! python

import json
import os
import strongdm, time, os, argparse, sys, logging, csv, subprocess, platform

# This will hold user/role/resource data.
audit_dict = {}
# Output file
filepath = 'audit.csv'

access_key=os.getenv("SDM_API_ACCESS_KEY")
secret_key=os.getenv("SDM_API_SECRET_KEY")

client = strongdm.Client(access_key, secret_key)

def get_role_by_id(id):
    try:
        sdm_roles = list(client.roles.list('id:"{}"'.format(id)))
    except Exception as ex:
        raise Exception("List roles failed: " + str(ex)) from ex
    if len(sdm_roles) == 0:
        raise Exception("Sorry, cannot find that role!")
    return sdm_roles[0]

def get_all_resources_by_role(role_id, filter = ''):
    """
    Return all resources by role name
    """
    try:
        sdm_role = get_role_by_id(role_id)
        resources_filters = get_resources_filters_by_role(sdm_role)
        if filter:
            resources_filters = [f"{rf},{filter}" for rf in resources_filters]
        return get_unique_resources(resources_filters)
    except Exception as ex:
        raise Exception("List resources by role failed: " + str(ex)) from ex

def get_resources_filters_by_role(sdm_role):
    if not hasattr(sdm_role, 'access_rules') or sdm_role.access_rules is None:
        sdm_role_grants = list(client.role_grants.list(f"role_id:{sdm_role.id}"))
        return [f"id:{rg.resource_id}" for rg in sdm_role_grants]
    # then this org is using Access Overhaul
    access_rules = json.loads(sdm_role.access_rules) if isinstance(sdm_role.access_rules, str) else sdm_role.access_rules
    resources_filters = []
    for ar in access_rules:
        filter = []
        if ar.get('ids'):
            filter.append(",".join([f"id:{id}" for id in ar['ids']]))
        if ar.get('type'):
            filter.append(f"type:{ar['type']}")
        if ar.get('tags'):
            tags = []
            for key, value in ar['tags'].items():
                tags.append('tag:"{}"="{}"'.format(key, value))
            filter.append(",".join(tags))
        resources_filters.append(",".join(filter))
    return resources_filters

def get_unique_resources(resources_filter):
    resources_map = {}
    for filter in resources_filter:
        resources = remove_none_values(client.resources.list(filter))
        resources_map |= {r.id: r for r in resources if resources_map.get(r.id) is None}
    return resources_map.values()

def remove_none_values(elements):
    return [e for e in elements if e is not None]

def get_all_users():
    users = list(client.accounts.list('type:user,suspended:false'))
    return users

def get_user_roles(users):
    for user in users:
        role_list = []
        resource_list = []
        attachments = list(client.account_attachments.list('account_id:{}'.format(user.id)))
        for a in attachments:
            role = client.roles.get(a.role_id)
            resources = get_all_resources_by_role(a.role_id)
            role_list.append(role.role.name)
            role_list.append("/")

            for resource in resources:
                # Append the name to resource_list
                resource_list.append(resource.name)
                resource_list.append("/")
            
    # Store the userID and email, for printing later
        audit_dict[user.id] = [user.email, role_list, resource_list]

def print_csv():
  with open(filepath, 'w', newline='') as csvfile:
      spamwriter = csv.writer(csvfile, delimiter=',',
                              quotechar='|', quoting=csv.QUOTE_MINIMAL)
      spamwriter.writerow(["AccountID", "Email", "Roles", "Resources"])
      for key, val in audit_dict.items():
        spamwriter.writerow([key,val[0], ' '.join(val[1]).strip(" / "), ' '.join(val[2]).strip(" / ") ] )

def open_csv():
  if platform.system() == 'Darwin':        # Mac
    subprocess.call(('open', filepath))
  elif platform.system() == 'Windows':    # Windows
      os.startfile(filepath)
  else:                                   # linux variants
      subprocess.call(('xdg-open', filepath))

def main():
    users = get_all_users()
    get_user_roles(users)
    print_csv()
    open_csv()

if __name__ == "__main__":
    main()


