#!/bin/bash

#  strongDM - Daily datasource local query log (JSON) export to AWS/GCP
#
#  Extracts all previous day queries from the gateway/relay
#  Compresses resulting file, pushes up to AWS and/or GCP
#
#  Cron daily, at midnight. crontab 0 0 * * *
#  Hosts: all gateways or relays used to access database datasources
#
#  Dan Keller - 12/14/2020


export AWS_ACCESS_KEY_ID=[key]
export AWS_SECRET_ACCESS_KEY=[key]

HOSTNAME=`hostname -s`
YESTERDAY_DATE=`date --date="yesterday" +'%Y-%m-%d'`

CLOUD_BUCKET="your-bucket"
SDM_LOG_DIR="/home/sdm/.sdm/logs"

AWS_STORAGE_PATH="s3://$CLOUD_BUCKET/log_exports/$HOSTNAME/query"
GCP_STORAGE_PATH="gs://$CLOUD_BUCKET/log_exports/$HOSTNAME/query"

# COMBINE ALL EVENTS FROM PREVIOUS DAY LOGS
#PREV_DAY_LOG=`cat $SDM_LOG_DIR/relay.*.log | grep "$YESTERDAY_DATE" | grep "uuid\":\"0"` #deprecated?
PREV_DAY_LOG=`cat $SDM_LOG_DIR/relay.*.log | grep "$YESTERDAY_DATE" | grep ",start,0"`

# PUSH LOGGING to AWS/GCP
echo "Exporting local datasource query logs for $YESTERDAY_DATE..."
echo "$PREV_DAY_LOG" | gzip | aws s3 cp - "$AWS_STORAGE_PATH/query_log.$YESTERDAY_DATE.gz"
echo "$PREV_DAY_LOG" | gzip | gsutil cp - "$GCP_STORAGE_PATH/query_log.$YESTERDAY_DATE.gz"

echo; echo "Process complete!"
exit 0
