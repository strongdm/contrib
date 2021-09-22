import logging
import strongdm

from config_template import config
from exceptions import NotFoundException

def create_sdm_service():
    client = strongdm.Client(config['SDM_API_ACCESS_KEY'], config['SDM_API_SECRET_KEY'])
    return SdmService(client)

# Class copied from https://github.com/strongdm/accessbot/blob/main/plugins/sdm/lib/sdm_service.py
class SdmService:
    def __init__(self, client):
        self.__client = client

    def grant_temporal_access(self, account_email, resource_name, start_from, valid_until):
        """
        Grant temporary access to a SDM resource for an account
        """
        try:
            logging.debug(
                "##SDM## SdmService.grant_temporary_access resource_id: %s account_id: %s start_from: %s valid_until: %s",
                resource_name, account_email, str(start_from), str(valid_until)
            )            
            sdm_grant = strongdm.AccountGrant(
                resource_id = self.get_resource_by_name(resource_name).id,
                account_id = self.get_account_by_email(account_email).id,
                start_from = start_from,
                valid_until = valid_until
            )
            self.__client.account_grants.create(sdm_grant)
        except Exception as ex:
            raise Exception("Grant failed: " + str(ex)) from ex

    def get_resource_by_name(self, name):
        """
        Return a SDM resouce by name
        """
        try:
            logging.debug("##SDM## SdmService.get_resource_by_name name: %s", name)
            sdm_resources = list(self.__client.resources.list('name:"{}"'.format(name)))
        except Exception as ex:
            raise Exception("List resources failed: " + str(ex)) from ex
        if len(sdm_resources) == 0:
            raise NotFoundException("Sorry, cannot find that resource!")
        return sdm_resources[0]

    def get_account_by_email(self, email):
        """
        Return a SDM account by email
        """
        try:
            logging.debug("##SDM## SdmService.get_account_by_email email: %s", email)
            sdm_accounts = list(self.__client.accounts.list('email:{}'.format(email)))
        except Exception as ex:
            raise Exception("List accounts failed: " + str(ex)) from ex
        if len(sdm_accounts) == 0:
            raise Exception("Sorry, cannot find your account!")
        return sdm_accounts[0]
