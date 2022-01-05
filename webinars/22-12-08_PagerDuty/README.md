### Overview
If you're using PagerDuty, then you already have on-call schedules mapped out for critical roles. But when someone is on-call, they may need more database or server access than they'd otherwise use. This is where strongDM temporary grants come in: you can integrate your PagerDuty on-call schedule with strongDM to automatically grant strongDM users access to additional resources during their on-call shifts. This Python example shows a simple way of managing the process.

The script has two major portions: first, look up who is on call for a specific schedule over a certain time period; second, parse these assignments with the strongDM SDK to grant temporary access to a datasource or server. One wrinkle is that two API calls are necessary to PagerDuty: first, getting the list of who is on call will give a list of users and user IDs, but not email addresses. Second, specific user lookups get us the email addresses of who is on call.

#### Requirements:
- [SDM API keys](https://www.strongdm.com/docs/admin-ui-guide/settings/admin-tokens/api-keys) defined in runtime environment
- Python 3
- strongDM Python module

#### Steps:
To get this script working in your environment, you'll need the following:

- A strongDM API key and secret with Datasource list and grant and User assign and list rights
- A strongDM Resource name
- A PagerDuty API key with read-only rights
- The schedule ID of a PagerDuty schedule you wish to use as the basis of the temporary grants

#### Notes:
- In order for this automation to work, your users will need to be identified by the same email addresses in PagerDuty and in strongDM.

_This script was presented in a strongDM webinar on 8 December 2021._

