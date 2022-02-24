# Purpose

This script serves as an example exporter that can monitor the health of resources ("Infrastructure") and nodes ("Gateways/Relays"). The script uses the following workflow:

1. Make an API call to strongDM's API to retrive information about resources and nodes (configurable by updating `update_interval` variable in `main()`)

2. Collect data about any resource or node that is tagged with the `alert_tag` in strongDM (`alert_tag` variable is configurable in `main()`)

3. Export metrics to a prometheus endpoint as a "Gauge" (0 for healthy, 1 for unhealthy)

*IMPORTANT NOTE*: Currently, the resources and nodes only perform automatic health checks every 12 hours, or when the check is manually initiated through the UI. There is a feature request in place to lower that automatic health check interval, and/or make it possible to initiate a manual check through the API.


# Setup

- Create a new strongDM API key

    https://www.strongdm.com/docs/admin-ui-guide/access/api-keys

- Configure the  environment variables:`

    `export SDM_API_ID="<id>"`

    `export SDM_API_SECRET="<secret>"`

- Create a new virtual environment

    `python3 -m venv venv`

- Activate the new environment

    `source venv/bin/activate`

- Install requirements with `pip`
    
    `pip install -r requirements.txt`


# Sample `/metrics` data

After starting the exporter, the new metrics will be available on `http://<ip or hostname>:8337/metrics`. Here is an example of what the exported metrics will look:

```
# HELP example_server1 health of resource
# TYPE example_server1 gauge
example_server1{id="rs-116<redacted>",name="example_server1",tags="{'lab-infra': '', 'alert': ''}",task="health_check"} 0.0
# HELP example_server2 health of resource
# TYPE example_server2 gauge
example_server2{id="rs-5c6<redacted>",name="example_server2",tags="{'lab-infra': '', 'alert': ''}",task="health_check"} 0.0
# HELP example_server3 health of resource
# TYPE example_server3 gauge
example_server3{id="n-181<redacted>",name="example_server3",tags="{'alert': ''}",task="health_check"} 0.0
```

# Scraping the metrics with Prometheus

Adding the following job to the `prometheus.yml` file will scrape the metrics from our new endpoint (replace `localhost` with the appropriate IP or hostname):

```
  - job_name: "sdm_health_exporter"
    static_configs:
      - targets: ["localhost:8337"]
```


# Alerting on failures

If you use Alert Manager with Prometheus, you can use a PromQL expression similar to the one below to alert on unhealthy resources/nodes:

`{job="sdm_health_exporter", task="health_check"} == 1`
