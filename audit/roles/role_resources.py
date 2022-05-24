import json
import os
import csv
import strongdm

access_key=os.getenv("SDM_API_ACCESS_KEY")
secret_key=os.getenv("SDM_API_SECRET_KEY")

client = strongdm.Client(access_key, secret_key)

def get_all_roles():
    """
    Return all roles
    """
    try:
        return list(client.roles.list(''))
    except Exception as ex:
        raise Exception("List roles failed: " + str(ex)) from ex

def get_role_by_name(name):
    """
    Return a SDM role by name
    """
    try:
        sdm_roles = list(client.roles.list('name:"{}"'.format(name)))
    except Exception as ex:
        raise Exception("List roles failed: " + str(ex)) from ex
    if len(sdm_roles) == 0:
        raise Exception("Sorry, cannot find that role!")
    return sdm_roles[0]

def get_all_resources_by_role(role_name, filter = ''):
    """
    Return all resources by role name
    """
    try:
        sdm_role = get_role_by_name(role_name)
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

def print_border():
  print("~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=")

def write_to_csv(roles):
    header = ['Role Name', 'Resources->']
    with open('role_resources.csv', 'w', encoding='UTF8') as f:
        # Write header to csv
        writer = csv.writer(f)
        writer.writerow(header)
        # Write a row for each Role, with columns for each resource
        # (opposite may serve better?)
        for role in roles:
            # Get resources that are accessible from a given role name
            resources = get_all_resources_by_role(role.name)
            row = [role.name] + [r.name for r in resources]
            writer.writerow(row)
        f.close()

def main():
  roles = get_all_roles()
  print("Creating CSV")
  write_to_csv(roles)
  print("Finished Writing CSV")
  exit()

if __name__ == "__main__":
    main()


