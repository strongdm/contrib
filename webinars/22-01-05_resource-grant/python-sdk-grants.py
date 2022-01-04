#!/usr/bin/env python

import strongdm, datetime
from datetime import timezone

# Global Variables

# API Key/Secret
api_key = "Your-API-Key-Here"
secret_key = "Your-Secret-Here"

# Role to grant access to
roleName = "YourRoleName"
# User to grant the role to
user = 'YourEmail@Address.com'
# Grant start and end date/time
startDate = "2022-01-05"
startTime = "00:01:00"
endDate = "2022-01-06"
endTime = "00:01:00"

def grant_access():
    # Gets SDM client based on the provided api key
    client = strongdm.Client(api_key, secret_key)
    # Gets the SDM Role specified in the roleName global variable
    roleResponse = client.roles.list("name:{role}".format(role=roleName))
    # Sets the start and end date/time for the grant based on the global variables
    start_date = startDate+"T"+startTime+"Z"
    end_date = endDate+"T"+endTime+"Z"
    s = datetime.datetime.strptime(start_date, '%Y-%m-%dT%H:%M:%SZ')
    e = datetime.datetime.strptime(end_date, '%Y-%m-%dT%H:%M:%SZ')
    start = s.replace(tzinfo=timezone.utc)
    end = e.replace(tzinfo=timezone.utc)
    #Get all SDM users and filter based on the user global variable 
    users = list(client.accounts.list('email:{}'.format(user)))
    #Select that user from the returned users list
    myUserID = users[0].id
     # Adds each resource from the role to the user
    for r in roleResponse:
        # Using the Role, User id, gets a list of associated resource grants
        rgResponse = client.role_grants.list(
            'role_id:{id}'.format(id=r.id))
        # Cycle through that list
        for g in rgResponse:
            # Create a temporary grant for each one
            myGrant = strongdm.AccountGrant(resource_id='{}'.format(g.resource_id), account_id='{}'.format(myUserID),
                        start_from=start, valid_until=end)
            # Assign the grant by "creating" it
            try:
                respGrant = client.account_grants.create(myGrant)
            except Exception as ex:
                print("\nSkipping user " + user + " because of error: " + str(ex))
            else:
                print("\nGrant succeeded for user " + user + " to a resource in role " + 
                r.name + " from {} to {}".format(start, end))
        print('---\n')

def main():
    grant_access()

main()