import os

config = {
    'JIRA_USER': os.getenv('JIRA_USER'),
    'JIRA_TOKEN': os.getenv('JIRA_TOKEN'),
    'JIRA_BASE_URL': os.getenv('JIRA_BASE_URL'),
    'SDM_API_ACCESS_KEY': os.getenv('SDM_API_ACCESS_KEY'),
    'SDM_API_SECRET_KEY': os.getenv('SDM_API_SECRET_KEY'),
    'SDM_ADMINS': os.getenv("SDM_ADMINS").split(" "),
    'GRANT_TIMEOUT': int(os.getenv('GRANT_TIMEOUT', '60')),
}