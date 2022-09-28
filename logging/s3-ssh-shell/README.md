# SDM log export to AWS or GCP

This folder contains shell scripts to export various SDM logs (e.g. gateway logs, Activities) to cloud storage in AWS (S3) or GCP. They are intended to be run on gateway hosts _before_ log export to the Cloud.

## Requirements
* Linux
* SDM logs
* AWS or GCP credentials
* A [strongDM admin token](https://www.strongdm.com/docs/admin-ui-guide/settings/admin-tokens/admin-tokens) with all `Audit` permissions.

## Configuration
Please review the header in each script for details on what it does, and how to run it.
