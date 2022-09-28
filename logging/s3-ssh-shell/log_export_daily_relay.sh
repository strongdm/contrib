#!/bin/bash

#  strongDM - Daily relay log export to AWS/GCP
#
#  Extracts all previous day relay logs from local gw/relay
#  Compresses, pushes up to AWS and/or GCP
#
#  Cron daily, at midnight. crontab 0 0 * * *
#  Hosts - all gateways and relays
#
#  Dan Keller - 12/9/2020


export AWS_ACCESS_KEY_ID=[key]
export AWS_SECRET_ACCESS_KEY=[key]

HOSTNAME=`hostname -s`
YESTERDAY_DATE=`date --date="yesterday" +'%Y-%m-%d'`

CLOUD_BUCKET="your-bucket"
SDM_LOG_DIR="/home/sdm/.sdm/logs"

AWS_STORAGE_PATH="s3://$CLOUD_BUCKET/log_exports/local-relay/$HOSTNAME/$YESTERDAY_DATE"
GCP_STORAGE_PATH="gs://$CLOUD_BUCKET/log_exports/local-relay/$HOSTNAME/$YESTERDAY_DATE"

# COMBINE ALL EVENTS FROM PREVIOUS DAY LOGS
PREV_DAY_LOG=`cat $SDM_LOG_DIR/relay.*.log | grep "$YESTERDAY_DATE"`

# PUSH LOGGING TO AWS/GCP
echo "Exporting relay logs for $YESTERDAY_DATE..."

echo "$PREV_DAY_LOG" | gzip | aws s3 cp - "$AWS_STORAGE_PATH/relay.$HOSTNAME.$YESTERDAY_DATE.gz"
echo "$PREV_DAY_LOG" | gzip | gsutil cp - "$GCP_STORAGE_PATH/relay.$HOSTNAME.$YESTERDAY_DATE.gz"

echo; echo "Process complete!"
exit 0
