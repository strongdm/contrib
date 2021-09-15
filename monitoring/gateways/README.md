### Overview
This shell script is intended to help troubleshoot the use of system resources on Unix machines running a strongDM Gateway or Relay. It may be useful in particular when a Gateway is failing, and you are unable to obtain logs from the given machine.

NB: for long-term monitoring, we recommend a more robust third-party tool like [Node Exporter](https://github.com/prometheus/node_exporter).

#### Requirements:
- An S3 bucket in AWS
- AWS credentials that allow writing to that bucket.

#### Notes:
- This script was tested on Ubuntu 18.04 and 20.04.
