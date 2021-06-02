### Overview
This Terraform script uses the AWS provider to create an EC2 instance, and the strongDM provider to create an [SSH-certificate server resource](https://www.strongdm.com/docs/admin-ui-guide/settings/ssh/ssh-certificate-auth) that points to it.

#### Requirements:
- [AWS credentials configured locally](https://docs.aws.amazon.com/sdk-for-java/v1/developer-guide/setup-credentials.html)
- [SDM API keys](https://www.strongdm.com/docs/admin-ui-guide/settings/admin-tokens/api-keys) defined in runtime environment
- Terraform 14 or later