# AWS ECS Terraform

This code will create a StrongDM Gateway running in ECS Fargate, along with related network and IAM objects.

Please note that deployments created by this code may incur costs in your AWS account.

## Pre-Requirements

1. AWS CLI Installed
2. SDM CLI Installed

3. SDM API Key and Secret

    * [How to make API Key for SDM](https://www.strongdm.com/docs/admin-ui-guide/settings/admin-tokens/api-keys)

4. VPC ID, Public Subnet ID, and Private Subnet ID to be used. Can be found in AWS Console. NB the 2 subnets need to be in the same Availability Zone

## Variables

You can provide variable values interactively, or you can create a `tfvars` file. For example, create a file in the same folder named `terraform.tfvars` like this:

```HCL
sdm_access_key = "value"
sdm_secret_key = "value"
region = "value"
vpc_id = "value"
public_subnet_id = "value"
private_subnet_id = "value"
```

The code will use the default AWS credentials used for the AWS CLI you have installed on your machine.

## Infrastructure Generation

This will automatically generate the IAM objects, Network Load Balancer, Security Groups, etc for an ECS gateway and will attach a random string at the end of most of the generated resources.

