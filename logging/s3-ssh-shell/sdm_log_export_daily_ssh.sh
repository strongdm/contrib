#!/bin/bash

#  strongDM - Daily SSH session commands export (from strongDM servers) to AWS/GCP
#
#  Extracts all previous day ssh events from strongDM logs
#  Creates a stdout file for each ssh session-id
#  Compresses, pushes up to AWS and/or GCP
#
#  Cron daily, at midnight. crontab 0 0 * * *
#  Required hosts: Any one single host within your environment
#
#  Dan Keller - 12/9/2020


export SDM_ADMIN_TOKEN=[token]
export AWS_ACCESS_KEY_ID=[id]
export AWS_SECRET_ACCESS_KEY=[key]

HOSTNAME=`hostname -s`
TODAY_DATE=`date +'%Y-%m-%d'`
YESTERDAY_DATE=`date --date="yesterday" +'%Y-%m-%d'`

CLOUD_BUCKET="your-bucket"
SDM_LOG_DIR="/home/sdm/.sdm/logs"
TMP_DIR="/home/sdm/.sdm/tmp"

AWS_STORAGE_PATH="s3://$CLOUD_BUCKET/log_exports/ssh/$YESTERDAY_DATE"
GCP_STORAGE_PATH="gs://$CLOUD_BUCKET/log_exports/ssh/$YESTERDAY_DATE"

# CHECK DIRS
test -d $SDM_LOG_DIR || exit 1
test -d $TMP_DIR     || exit 1

# FLUSH TMP DIR
rm -f $TMP_DIR/*

# PULL SSH EVENTS FROM SDM SERVERS, FILTER ONLY SESSION-ID
PREV_DAY_SSH_SESSIONS=`sdm audit ssh --from "$YESTERDAY_DATE" --to "$TODAY_DATE" | cut -d, -f7 | uniq | egrep "^s"`

# CREATE SSH DUMPS FOR EACH SESSION-ID
cd $TMP_DIR; for session in $PREV_DAY_SSH_SESSIONS; do echo $session; sdm ssh dump $session > ./$session.stdout; done

# TEST AND COMPRESS SESSION FILES, OTHERWISE TRY AGAIN TOMORROW
ls *.stdout > /dev/null 2>&1 || exit 1
gzip *.stdout

# PUSH LOGGING UP TO AWS/GCP
echo "Exporting SSH session commands for $YESTERDAY_DATE..."

aws s3 cp $TMP_DIR/ "$AWS_STORAGE_PATH/stdout/" --recursive --exclude "*" --include "*.stdout.gz"
gsutil -m cp $TMP_DIR/*.stdout.gz "$GCP_STORAGE_PATH/stdout/"

echo "Process complete!"
exit 0
