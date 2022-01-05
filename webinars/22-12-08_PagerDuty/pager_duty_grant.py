#!/usr/bin/env python

import requests,json,datetime,subprocess,strongdm,re
from datetime import timezone

# PagerDuty API key
API_KEY = 'PD_API_KEY'

# strongDM API keys: requires [datasource: list,grant] and [user: list,assign]
access_key = "SDM_ACCESS_KEY"
secret_key = "SDM_SECRET_KEY"

# Name of strongDM Datasource to which you are granting access, from Admin UI
DATASOURCE = 'RESOURCE_NAME'

# Set your time zone for PagerDuty
TIME_ZONE = 'UTC'
# Get this ID from the PagerDuty admin UI, or via their 'schedules' API endpoint
SCHEDULE_IDS = ['ID']
# for the PD API requests. Modify UNTIL with the proper time offset
UNTIL = (datetime.timedelta(days=1) + datetime.datetime.utcnow()).isoformat() + 'Z'

def get_oncalls():
  url = 'https://api.pagerduty.com/oncalls'
  headers = {
    'Accept': 'application/vnd.pagerduty+json;version=2',
    'Authorization': 'Token token={token}'.format(token=API_KEY)
  }
  payload = {
    'time_zone': TIME_ZONE,
    'schedule_ids[]': SCHEDULE_IDS,
    'until': UNTIL,
  }
  
  r = requests.get(url, headers=headers, params=payload)
  struct = r.json()
  output = []

  for record in struct["oncalls"]:
  # get user's email address
    r = requests.get(record["user"]["self"], headers=headers)
    output.append({"email" : r.json()["user"]["email"],
            "from" : record["start"],
            "to" : record["end"]})
  return output

def grant_access(access_list):

  client = strongdm.Client(access_key, secret_key)

  # Get Datasource(s)
  resources = list(client.resources.list('name:{}'.format(DATASOURCE)) )
  resourceID = resources[0].id

  # Cycle through the output from PagerDuty
  for item in access_list:
    # Use the PD email address to get the user from SDM
    print('Current PD user is: ' + item["email"])
    users = list(client.accounts.list('email:{}'.format(item["email"])))
    if len(users) > 0:
      print('SDM user found!')
      myUserID = users[0].id
      # Convert the date strings from PD into a datetime object
      s = datetime.datetime.strptime(item["from"], '%Y-%m-%dT%H:%M:%SZ')
      e = datetime.datetime.strptime(item["to"], '%Y-%m-%dT%H:%M:%SZ')
      # Make both objects 'aware' (with TZ) as required by the strongDM SDK
      start = s.replace(tzinfo=timezone.utc)
      end = e.replace(tzinfo=timezone.utc)
      # Create the grant object
      myGrant = strongdm.AccountGrant(resource_id='{}'.format(resourceID),account_id='{}'.format(myUserID), 
        start_from=start, valid_until=end)
      # Perform the grant
      try:
        respGrant = client.account_grants.create(myGrant)
      except Exception as ex:
        print("\nSkipping user " + item["email"] + " on account of error: " + str(ex))
      else:
        print("\nGrant succeeded for user " + item["email"] + " to resource " + DATASOURCE + " from {} to {}".format(start,end))
    print('---\n')

def main():
  access = get_oncalls()
  grant_access(access)

main()