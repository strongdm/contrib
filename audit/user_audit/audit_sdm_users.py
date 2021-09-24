#!/usr/bin/env python3

import strongdm, time, os, argparse, sys, logging, csv, subprocess, platform

# This will hold user/role/resource data.
audit_dict = {}
# Output file
filepath = 'audit.csv'

access_key = os.environ['SDM_API_ACCESS_KEY']
secret_key = os.environ['SDM_API_SECRET_KEY']

client = strongdm.Client(access_key, secret_key)

def get_data():
  user_response = client.accounts.list('type:user')
  # Loop through the user list
  for u in user_response:
    # Create/reset role and resource lists at user level
    role_list = []
    resource_list = []
    # Use the account.id to get ACCOUNT ATTACHMENTS
    account_attach = client.account_attachments.list('account_id:{}'.format(u.id))
    # Loop 2
    for r in account_attach:
        roles = list(client.roles.list('id:{}'.format(r.role_id)))
        # Loop through the roles to which this user is assigned
        for o in roles:
          # Append the name to role_list
          role_list.append(o.name)
          role_list.append("/")
          # Using the role_id, get role_attachments
          role_grants = client.role_grants.list('role_id:{}'.format(o.id))
          for a in role_grants:
            resources = list(client.resources.list('id:{}'.format(a.resource_id)))
            for e in resources:
              # Append the name to resource_list
              resource_list.append(e.name)
              resource_list.append("/")

    # Store the userID and email, for printing later
    audit_dict[u.id] = [u.email, role_list, resource_list]


def print_csv():
  with open(filepath, 'w', newline='') as csvfile:
      spamwriter = csv.writer(csvfile, delimiter=',',
                              quotechar='|', quoting=csv.QUOTE_MINIMAL)
      spamwriter.writerow(["AccountID", "Email", "Roles", "Resources"])
      for key, val in audit_dict.items():
        spamwriter.writerow([key,val[0], ' '.join(val[1]).strip(" / "), ' '.join(val[2]).strip(" / ") ] )

def open_csv():
  if platform.system() == 'Darwin':
    subprocess.call(('open', filepath))
  elif platform.system() == 'Windows':    # Windows
      os.startfile(filepath)
  else:                                   # linux variants
      subprocess.call(('xdg-open', filepath))




def main():
  get_data()
  print_csv()
  open_csv()

if __name__ == "__main__":
    main()