import os
import time
from typing import Any
from pydantic import BaseModel
import strongdm
from prometheus_client import start_http_server, Gauge

"""
sdm_exporter.py

This script serves as an example exporter that can monitor the
health of resources ("Infrastructure") and nodes ("Gateways/Relays").

The script uses the following workflow:

- Make an API call to strongDM's API to retrive information about resources
and nodes. The frequency of the API call is configurable by updating the
"update_interval" variable in "main()"

- Collect data about any resource or node that is tagged with "alert" in strongDM.
This tag is configurable by updating the "alert_tag" variable in "main()"

- Export metrics to a prometheus endpoint as a gauge (0 for healthy, 1 for unhealthy)
"""

class SdmObject(BaseModel):   
    """
    Pydantic class with type enforcement
    Prometheus object of type "Gauge"
    Possible values: healthy: 0, unhealthy: 1
    """
    health_metric: Gauge

    class Config:
        validate_assignment = True
        arbitrary_types_allowed = True


def strip_invalid_chars(name):
    """
    Dashes are invalid for prometheus metric but allowed in strongDM
    """
    return name.replace('-','_')


def get_sdm_keys():
    """
    Loads API keys from environment variables.
    Raises: KeyError if environment variables are not found
    Returns: strings for API ID and secret
    """

    try:
        api_id = os.environ['SDM_API_ACCESS_KEY']
        api_secret = os.environ['SDM_API_SECRET_KEY']
    except KeyError as ke:
        print(f'FATAL: Missing env variable: {ke}. Exiting...')
        exit()

    return api_id, api_secret


def get_sdm_resources(client, alert_tag, sdm_objects):
    """
    Issue strongDM API calls every <update_interval> to collect information on resources
    Export any resource that's tagged with <alert_tag> in strongDM to a prometheus exporter
    Returns: dictionary of SdmObject instances, keyed by resource name
    """

    # prometheus labels that will be attached to the exported metric
    labels = ['id','healthy','name','tags','task']

    for resource in client.resources.list(''):
        
        # make sure there are no characters that are invalid for the prometheus exporter
        resource_name = strip_invalid_chars(resource.name)

        # only process resources that are tagged with <alert_tag> in strongDM
        if alert_tag not in resource.tags:
            continue

        # only register a new prometheus collector if it's not an existing object
        # create a new "Gauge" and define the labels that will be used for the metric
        # create new SdmObject with this information and add it to "sdm_objects" dictionary
        if resource_name not in sdm_objects:
            health_metric = Gauge(resource_name, 'health of resource', labelnames=labels)
            sdm_objects[resource_name] = SdmObject(health_metric=health_metric)

        # resource health is returned as True or False
        # set the label values based on information retrieved from the API call
        # set the metric to 0 (healthy) or 1 (unhealthy)
        sdm_objects[resource_name].health_metric.labels(
            id=resource.id,
            healthy=resource.healthy,
            name=resource_name,                       
            tags=resource.tags,
            task='health_check'
        ).set(0 if resource.healthy else 1)

    return sdm_objects


def get_sdm_nodes(client, alert_tag, sdm_objects):
    """
    Issue strongDM API calls every <update_interval> to collect information on nodes
    Export any resource that's tagged with <alert_tag> in strongDM to a prometheus exporter
    Returns: dictionary of SdmObject instances, keyed by node name
    """

    # prometheus labels that will be attached to the exported metric
    labels = ['id','name','state','tags','task']

    for node in client.nodes.list(''):
        
        # make sure there are no characters that are invalid for the prometheus exporter
        node_name = strip_invalid_chars(node.name)

        # only process nodes that are tagged with <alert_tag> in strongDM
        if alert_tag not in node.tags:
            continue

        # only register a new prometheus collector if it's not an existing object
        # create a new "Gauge" and define the labels that will be used for the metric
        # create new SdmObject with this information and add it to "sdm_objects" dictionary
        if node_name not in sdm_objects:
            health_metric = Gauge(node_name, 'health of resource', labelnames=labels)
            sdm_objects[node_name] = SdmObject(health_metric=health_metric)

        # node health is returned as "started", "stopped", or "new"
        # anything other than "stated" is considered unhealthy
        if node.state not in "started":

            # set the label values based on information retrieved from the API call
            # set the metric to 0 (healthy) or 1 (unhealthy)
            sdm_objects[node_name].health_metric.labels(
                id=node.id,
                name=node_name,
                state=node.state,
                tags=node.tags,
                task='health_check'
            ).set(0 if node.state == "started" else 1)

    return sdm_objects


def main():

    # frequency to issue API call to strongDM
    update_interval = 60

    # filter strongDM results to resources that are tagged with <alert_tag> in strongDM
    alert_tag = 'alert'

    # retrieve strongDM API keys
    api_id, api_secret = get_sdm_keys()

    # keeps track of SdmObject instances, keyed by resource name
    sdm_objects = {}

    # start prometheus endpoint
    start_http_server(8337) 

    # strongDM API client
    client = strongdm.Client(api_id, api_secret)

    # issue API calls to obtain health status indefinitely with a frequency of <update_interval>
    while True:    
    
        # collect resources and objects and write them out to a Prometheus exporter
        # returns dictionary of SdmObject instances, keyed by resource or node name 
        sdm_objects = get_sdm_resources(client, alert_tag, sdm_objects)
        sdm_objects = get_sdm_nodes(client, alert_tag, sdm_objects)
        
        # wait "update_interval" amount of time before issuing the next API call
        time.sleep(update_interval)


if __name__ == '__main__':
    main()
