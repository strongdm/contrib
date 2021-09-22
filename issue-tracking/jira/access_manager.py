import datetime
from exceptions import NotFoundException, PermissionDeniedException
import logging
import re
import requests

from config_template import config
from sdm_service import create_sdm_service

ACCESS_REGEX = r"^access to (.+)$"

class AccessManager:
    def __init__(self):
        self.__jira_api_url = f"{config['JIRA_BASE_URL']}/rest/api/2"
        self.__http_auth = (config['JIRA_USER'], config['JIRA_TOKEN'])
        self.__sdm_service = create_sdm_service()

    def process_issue(self, issue_id):
        resp = requests.get(
            f"{self.__jira_api_url}/issue/{issue_id}", 
            auth = self.__http_auth
        )
        fields = resp.json()['fields']
        description = fields['description']
        if not re.compile(ACCESS_REGEX).match(description):
            return
        creator_account_id = fields['creator']['accountId']
        creator_email = self.__get_account_email(creator_account_id)
        assignee_email = self.__get_assignee_email(fields['assignee'])
        if assignee_email not in config['SDM_ADMINS']:
            raise PermissionDeniedException(f"{assignee_email} cannot approve access requests, not an SDM_ADMIN")
        resource_name = re.sub(ACCESS_REGEX, "\\1", description)
        self.__grant_temporal_access(creator_email, resource_name)        

    def __get_account_email(self, account_id):
        resp = requests.get(
            f"{self.__jira_api_url}/user?accountId={account_id}", 
            auth = self.__http_auth
        )
        data = resp.json()
        if 'emailAddress' not in data:
            raise NotFoundException("Creator email not available, please check your profile and visibility settings")
        return data['emailAddress']

    def __get_assignee_email(self, assignee):
        if 'emailAddress' not in assignee:
            raise NotFoundException("Assignee email not available, please check your profile and visibility settings")
        return assignee['emailAddress']

    def __grant_temporal_access(self, account_email, resource_name):
        grant_start_from = datetime.datetime.now(datetime.timezone.utc)
        grant_valid_until = grant_start_from + datetime.timedelta(minutes = config['GRANT_TIMEOUT'])
        self.__sdm_service.grant_temporal_access(account_email, resource_name, grant_start_from, grant_valid_until)

