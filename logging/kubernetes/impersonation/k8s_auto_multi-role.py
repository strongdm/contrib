#!/usr/bin/env python3

import strongdm, time, os, argparse, sys, logging

access_key = os.environ['SDM_API_ACCESS_KEY']
secret_key = os.environ['SDM_API_SECRET_KEY']

# Change INFO to ERROR if you don't care about success messages
logging.basicConfig(level = logging.INFO)

# In this multi-role version, both "r" and "m" keys/values are required
# Add tags to make params required in argparse?

parser = argparse.ArgumentParser()
parser.add_argument("-r", "--role", help="SDM role name to convert",required=True)
parser.add_argument("-m", "--map", help="SDM mapping role to add",required=True)
args = parser.parse_args()

client = strongdm.Client(access_key, secret_key)

def transfer_users(role, id):
  # get the role in the arg
    role_response = list(client.roles.list('name:\"{role_name}\"'.format(role_name=role)))
    if len(role_response) == 0:
      print("Could not find Role: " + args.role)
      exit(0)
    # loop through the list
    for r in role_response:
      # get the account attachments
      attachments = list(client.account_attachments.list('role_id:{id}'.format(id=r.id)))
      # loop through that list
      for a in attachments:
        # create a new account attachment using: user, and "m" role from arg
        k8s_attachment = strongdm.AccountAttachment(account_id=a.account_id, role_id=id)
        try:
          respGrant = client.account_attachments.create(k8s_attachment)
          logging.info('Role assignment to' + a.account_id + ' succeeded.')
        except Exception as ex:
          logging.info('Role assignment to ' + a.account_id + ' failed.')
          logging.error("Role assignment role assignment because of error: " + str(ex))

def get_k8s_role(role):
  k8s_role_id = ""
  role_response = list(client.roles.list('name:\"{role_name}\"'.format(role_name=role)))
  for r in role_response:
    k8s_role_id = r.id
  return k8s_role_id

def main():
  my_id = get_k8s_role(args.map)
  try:
    transfer_users(args.role, my_id)
    logging.info("Script execution is complete.")
  except Exception as ex:
    logging.error("Script failed to complete, because of error: " + str(ex))

if __name__ == "__main__":
    main()