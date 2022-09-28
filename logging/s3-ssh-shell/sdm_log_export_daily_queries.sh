#!/bin/bash

#  strongDM - Daily datasource query log export to AWS/GCP
#
#  Extracts all previous day queries from the strongDM servers
#  Compresses resulting file, pushes up to AWS and/or GCP
#
#  Cron daily, at midnight. crontab 0 0 * * *
#  Required hosts: Any one single host within your environment.
#
#  Dan Keller - 12/8/2020


export SDM_ADMIN_TOKEN=[token]
export AWS_ACCESS_KEY_ID=[key]
export AWS_SECRET_ACCESS_KEY=[key]

TODAY_DATE=`date +'%Y-%m-%d'`
YESTERDAY_DATE=`date --date="yesterday" +'%Y-%m-%d'`

CLOUD_BUCKET="your-bucket"
AWS_STORAGE_PATH="s3://$CLOUD_BUCKET/log_exports/query"
GCP_STORAGE_PATH="gs://$CLOUD_BUCKET/log_exports/query"

# PULL QUERY LOGS FROM SDM SERVERS
PREV_DAY_LOG=`sdm audit queries --from "$YESTERDAY_DATE" --to "$TODAY_DATE"`

# PUSH LOGGING to AWS/GCP
echo "Exporting datasource query logs for $YESTERDAY_DATE..."
echo "$PREV_DAY_LOG" | gzip | aws s3 cp - "$AWS_STORAGE_PATH/query_log.$YESTERDAY_DATE.gz"
echo "$PREV_DAY_LOG" | gzip | gsutil cp - "$GCP_STORAGE_PATH/query_log.$YESTERDAY_DATE.gz"

echo; echo "Process complete!"
exit 0
