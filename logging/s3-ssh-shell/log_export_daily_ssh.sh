#!/bin/bash

#  strongDM - Daily SSH session log export to AWS/GCP
#
#  Extracts all previous day ssh sessions from local gw/relay logging
#  Creates stdout files from resulting ssh session files
#  Compresses both ssh filetypes, pushes up to AWS and/or GCP
#
#  Cron daily, at midnight. crontab 0 0 * * *
#  Hosts - all gateways or relays used to access ssh server endpoints
#
#  Dan Keller - 12/9/2020

export AWS_ACCESS_KEY_ID=[key]
export AWS_SECRET_ACCESS_KEY=[key]

HOSTNAME=`hostname -s`
YESTERDAY_DATE=`date --date="yesterday" +'%Y-%m-%d'`

CLOUD_BUCKET="your-bucket"
SDM_LOG_DIR="/home/sdm/.sdm/logs"
TMP_DIR="/home/sdm/.sdm/tmp"

AWS_STORAGE_PATH="s3://$CLOUD_BUCKET/log_exports/local-relays/$HOSTNAME/$YESTERDAY_DATE/ssh"
GCP_STORAGE_PATH="gs://$CLOUD_BUCKET/log_exports/local-relays/$HOSTNAME/$YESTERDAY_DATE/ssh"

# CHECK DIRS
test -d $SDM_LOG_DIR || exit 1
test -d $TMP_DIR || exit 1

# FLUSH TMP DIR
rm -f $TMP_DIR/*

# COMBINE ALL EVENTS FROM PREVIOUS DAY LOGS
cat $SDM_LOG_DIR/relay.*.log | grep "$YESTERDAY_DATE" > $TMP_DIR/combined_$YESTERDAY_DATE.log

# EXTRACT SSH SESSIONS FROM COMBINED LOGGING
cd $TMP_DIR; sdm ssh split ./combined_$YESTERDAY_DATE.log

# CHECK FOR SSH FILES, OTHERWISE TRY AGAIN TOMORROW
ls *.ssh > /dev/null 2>&1 || exit 1

# EXTRACT SSH STDOUT FROM SSH SESSION FILES AND COMPRESS
for ssh_file in `ls *.ssh`; do sdm ssh dump -f $ssh_file > $ssh_file.stdout; done
gzip *.ssh*

# PUSH LOGGING TO AWS/GCP
echo "Exporting SSH session logs for $YESTERDAY_DATE..."
aws s3 cp $TMP_DIR/ "$AWS_STORAGE_PATH/sessions/" --recursive --exclude "*" --include "*.ssh.gz"
aws s3 cp $TMP_DIR/ "$AWS_STORAGE_PATH/stdout/"   --recursive --exclude "*" --include "*.stdout.gz"

gsutil -m cp $TMP_DIR/*.ssh.gz "$GCP_STORAGE_PATH/sessions/"
gsutil -m cp $TMP_DIR/*.stdout.gz "$GCP_STORAGE_PATH/stdout/"

echo; echo "Process complete!"
exit 0
