#!/usr/bin/env python3

import strongdm, time, os, argparse, sys, logging

# Set logging level
logging.basicConfig(level = logging.INFO)

# Parse command line arguments
# The role/-r argument is required, the map/-m is optional
parser = argparse.ArgumentParser()
parser.add_argument("-r", "--role", help="SDM role name to convert", required=True)
parser.add_argument("-m", "--map", help="SDM mapping role to add")
args = parser.parse_args()

# Get SDM API keys from environment
access_key = os.environ['SDM_API_ACCESS_KEY']
secret_key = os.environ['SDM_API_SECRET_KEY']

client = strongdm.Client(access_key, secret_key)

# Function to create the Composite Role
def createCompRole(roleName):
  compRole = strongdm.Role(
    name=roleName + "_k8s",
    composite=True
  )
  comp_response = client.roles.create(compRole, timeout=30)
  return comp_response

# Function to assign normal roles to Composite
def addRoleToComp(role, compRoleId):
    compAttachment = strongdm.RoleAttachment(
    composite_role_id=compRoleId,
    attached_role_id=role.id,
    )
    client.role_attachments.create(compAttachment, timeout=30)

def main():
  if args.role:
    try:
      # Create the new Composite role
      resp = createCompRole(args.role)
    except strongdm.AlreadyExistsError as ex:
      logging.error('The Composite Role already exists! Please check the Admin UI to confirm the specified role is assigned to it.')
      exit(-1)
    except Exception as ex:
      logging.error('Failed to create Composite Role: '+ str(ex))
      exit(-1)
    logging.info('Composite Role %s_k8s created successfully!' % args.role)

    # Get the role passed in via CLI
    list = client.roles.list('name:"%s"' % args.role)
    for u in list:
      try:
        # Add that role to the new Composite Role
        addRoleToComp(u, resp.role.id)
      except Exception as ex:
        logging.error('Failed to add role %s to Composite Role: '+ str(ex) % args.role)
        exit(-1)
      logging.info('Role %s assigned successfully!' % args.role)

      # Perform the same operation for the 'map' role, if specified
      if args.map:
        list = client.roles.list('name:"%s"' % args.map)
        for v in list:
          try:
            addRoleToComp(v, resp.role.id)
          except Exception as ex:
            logging.error('Failed to add role to Composite Role: '+ str(ex))
            exit(-1)
          logging.info('Map role %s assigned successfully!' % args.map)

if __name__ == "__main__":
    main()





