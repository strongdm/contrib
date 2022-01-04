### Overview
This Python script (python-sdk-grants.py) uses the strongDM SDK to assign resources with in a role to a user for set time period. Follow the steps below to add in the details for your environment to get started.

#### Requirements:
- [SDM API keys](https://www.strongdm.com/docs/admin-ui-guide/settings/admin-tokens/api-keys) defined in runtime environment
- Python 3
- strongDM Python module

#### Steps:
You will need to update the following lines with entries that exist in your environment.
- line 9: Your SDM API key
- line 10: Your SDM Secret
- line 13: The role you wish to assign to a user
- line 15: The user you wish to assign the resources too, defined by their email address as it exists with in SDM
- Line 17 & 18: The start date and time for the access to begin
- Line 19 & 20: The end date and time for the access to be removed  

#### Notes:
- [SDM API Documentation](https://www.strongdm.com/docs/api)
