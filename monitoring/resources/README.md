### Overview
This Python script uses the strongDM SDK to find unhealthy resources (those that appear yellow in the UI) and "tags" them with a timestamp. The tagging process forces a health check on each resource, just as if you had clicked the `Check Now` button in the Admin UI.

This process doesn't affect any other tags you might use.

#### Requirements:
- [SDM API keys](https://www.strongdm.com/docs/admin-ui-guide/settings/admin-tokens/api-keys) defined in runtime environment
- Python 3
- strongDM Python module

#### Notes:
- This script can take some time to run if you have many unhealthy resources.
- We strongly recommend that you run this script single-threaded.