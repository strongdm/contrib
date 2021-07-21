#!/usr/bin/env python

import strongdm, os

access_key=os.getenv("SDM_API_ACCESS_KEY")
secret_key=os.getenv("SDM_API_SECRET_KEY")

client = strongdm.Client(access_key, secret_key)

def get_role_details():
    roleResponse = list(client.roles.list(""))
    print_border()
    for r in roleResponse:
      print("Role name: \"" + r.name + "\" includes the following resources:")
      rgResponse = client.role_grants.list('role_id:{id}'.format(id=r.id) )

      for g in rgResponse:
        res = list(client.resources.list("id:{}".format(g.resource_id)))
        print("\t" + res[0].name)
      print_border()

def print_border():
  print("~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=")

def main():
  get_role_details()

if __name__ == "__main__":
    main()






