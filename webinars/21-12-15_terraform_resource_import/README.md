# Importing resources into strongDM using Terraform.

This is a small demo script showing how to add resources into strongDM if you are spinning them up in Terraform. This module will create a whole new VPC with public and private subnet and create a PSQL database and ssh resource. This is purely referencial and not for production environments.

1. Download this repo

2. Navigate into it this folder in your terminal

3. Run `terraform init`

4. You'll need to edit/create the `terraform.tfvars` file with your required variables. To generate an API key log in to the [strongDM admin UI](https://app.strongdm.com/) and generate a new API Key. You can see this docs page for instructions on how to [generate an API Key](https://www.strongdm.com/docs/admin-ui-guide/access/api-keys). File contentes example below:

    ```HCL
    sdm_access_key = "YOURSDMAPIACCESSKEY"
    sdm_secret_key = "YOURSDMAPISECRETKEY"
    region         = "us-west-1"
    ```

5. Run `terraform plan` to verify all the changes and any errors.

6. Run `terraform apply` to deploy the terraform. Give about 15-20 minutes to execute as the Relay and PSQL instance wait until the NAT Gateway has come online to validate network connectivity and installation script to run.

7. Check your [strongDM Admin Panel](https://app.strongdm.com/) and verify resources are online.

8. If you wish to test resources and connectivity assign a role that has permissions to the PSQL resource and ssh resources.