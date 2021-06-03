import strongdm, time, os

access_key=os.getenv("SDM_API_ACCESS_KEY")
secret_key=os.getenv("SDM_API_SECRET_KEY")

# Create SDM client
client = strongdm.Client(access_key, secret_key)

# Get time details for key update
import time
seconds = int( time.time() )
humanTime = time.strftime('%Y-%m-%d %H:%M:%S', time.localtime(seconds))

# Get all resources of type DATASOURCE with healthy=FALSE
resources = list(client.resources.list('category:datasource, healthy:false'))
print("\nThe following Datasources are NOT healthy, and will be updated with a fresh tag:\n")
for i in resources:
    print("Name:", i.name," | Healthy?", i.healthy)
    # create the tag locally, with current timestamp
    i.tags = {"lastHealthcheck": "{}".format(seconds)}
    # this will add the tag if not present, or update if so
    response = client.resources.update(i)

# Wait 5 seconds per unhealthy resource
mTime = len(resources) * 5
print('~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=')
print('Running healthchecks ...')
time.sleep(mTime)
print('Work complete.')
print('~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=~=')

# Get a new list of resources, to see what changes the update above made
resources = list(client.resources.list('category:datasource, healthy:false'))
print("The following Datasources are still unhealthy:\n")
for i in resources:
    time.strftime('%Y-%m-%d %H:%M:%S', time.localtime(1347517370))
    print("Name:", i.name," | Healthy?", i.healthy, " | Last check:", humanTime )

print("\nDone.")

