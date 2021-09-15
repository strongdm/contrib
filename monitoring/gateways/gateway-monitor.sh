#!/bin/bash
# day, hour, minute timestamp
TIMESTAMP=`date +'%Y%m%d%H%M'`

S3NAME=strongdm-gw-health-$TIMESTAMP.gz
S3PATH=s3://s3-bucket-name # no trailing slash

export AWS_ACCESS_KEY_ID=id
export AWS_SECRET_ACCESS_KEY=secret

(free -m | awk 'NR==2{printf "Memory Usage: %s/%sMB (%.2f%%)\n", $3,$2,$3*100/$2 }' ; \
df -h | awk '$NF=="/"{printf "Disk Usage: %d/%dGB (%s)\n", $3,$2,$5}' ; \
top -bn1 | grep load | awk '{printf "CPU Load: %.2f\n", $(NF-2)}' ; \
lsof | wc -l | awk '{printf "Open files: " $1"\n"}') | \
gzip | aws s3 cp - $S3PATH/$S3NAME