# Export Logs with Lifecycle Hooks

In this tutorial, you'll learn how to configure lifecycle hooks to export logs from a self-terminating instance to an AWS S3 bucket. Lifecycle hooks tell Amazon EC2 Auto Scaling what action to take when it launches or terminates instances. In this case, we will be configuring lifecycle hooks to allow the instance to export logs before being powered off or destroyed.

If you don't use S3 buckets, you can easily use the principles explained in this tutorial to export logs to other storage destinations.

---

``This technique is unsupported and is only provided as a potentially useful idea.``

---

## Prerequisites

- You must have an AWS account and be familiar with the AWS Console.
- All resources must have AWS CLI installed and keys configured.
- All have to be able to ping AWS resources to access and back them up.

## Steps

### Create a Simple Notification System (SNS) topic

1. In the Services menu, click **Simple Notification System**.
2. Click **Topics**.
3. Click **Create Topic**.
4. Select **Standard** type.
5. Enter a name.
6. Click **Create**.
7. In the topics menu, save the Amazon Resource Name (ARN), the unique identifier of the resource, for later use. Example: `arn:aws:sns:{{region}}:XXXXXXXXXXXX:{{name}}`

### Create a policy

1. In the Services menu, click **IAM**.
2. Click **Policies**.
3. Click **Create Policy**.
4. Click the **JSON Editor**.
5. Copy the content of [create_policy.json](./create_policy.json) file and paste it into the JSON Editor.
6. Click **Review**.
7. Enter a name for the policy.
8. Click **Create Policy**.

### Create an EC2 role

1. In the IAM menu, click **Roles**.
2. Click **Create Role**.
3. Click on **EC2 Service**.
4. Attach the policy created in [Create a Policy](#create-a-policy).
5. Click **Tags**.
6. Click **Review**.
7. Enter a name for the role.
8. Click **Create Role**.

### Create a Lambda role

1. In the IAM menu, click **Roles**.
2. lick **Create Role**.
3. Click on **Lambda Service**.
4. Attach the policy created in [Create a Policy](#create-a-policy).
5. Click **Tags**.
6. Click **Review**.
7. Enter a name for the role.
8. Click **Create Role**.

### Create a lifecycle hook

1. In the Services menu, click **EC2**.
2. Click **Auto-Scaling Groups**.
3. Click on the Auto-Scaling Group you want to edit.
4. Click **Instance Management**.
5. Click **Create Lifecycle Hook** and set the following:
    1. **Name**: Enter a name for the lifecycle hook.
    2. **Lifecycle transition**: Set "Instance terminate."
    3. **Heartbeat timeout**: Set "600."
    4. **Default result**: Set to "Continue."
6. Click **Create**.

### Create an S3 bucket

1. In the Services menu, click **S3**.
2. Click **Create Bucket**.
3. Enter **Bucket Name**.
4. Select the **Region** for the bucket.
5. Click **Create Bucket**.

### Create a Systems Manager (SSM) document

1. In the Services menu, click **System Manager**.
2. Click **Documents**.
3. Click **Owned By Me**.
4. Click **Create Command**.
5. Enter a **Name**.
6. Copy the content of [create_ssm_document.json](./create_ssm_document.json) file and paste it.
7. Click **Create document**.

### Create a Lambda Function

1. In the Services menu, click **Lambda**.
2. Click **Create function**.
3. Select **author from scratch** and then set the following:
    1. **Name**: Enter a name.
    2. **Runtime**: Set "Node.js 10.x."
    3. **Role**: Set "use existing role" and select the Lambda role created in [Create a Lambda Role](#create-a-lambda-role).
4. Click **Create function**.
5. Copy the content of [create_lambda_function.js](./create_lambda_function/index.js) file (and replace everything in the window).
6. Go down to **Environment Variables**.
7. Add the variables and the values for your variables using the content of [s3_vars.env](./create_lambda_function/s3_vars.env) as an example.
8. Click **Deploy**.

### Create a CloudWatch trigger

1. In Services menu, click **CloudWatch**.
2. Click **Rules**.
3. Click **Create Rule** and set the following:
    1. **Service Name**: Set "Auto-Scaling."
    2. **Event Type**: Set "Instance Launch and Terminate."
    3. **Instance Event**: Edit this to be "specific instance event(s)" and set "EC2 Instance-terminate Lifecycle Action."
    4. **Group**: Edit the group to be '"specific group name(s)" and set it to the auto-scaling group you want to monitor.
4. Click **Add target**.
5. Select the Lambda function you created in [Create a Lambda Function](#create-a-lambda-function).
6. Click **Configure details**.
7. Enter a **Name**.
8. Click **Create rule**.

### Test that logs are exported to S3

1. Go back to **EC2**.
2. Change size of your auto-scaling group in order for it to terminate an instance. **Note: Manual termination doesn't trigger the lifecycle event**.
3. Wait for the instance to terminate and then check S3 for logs.
4. For troubleshooting, you should be able to check the logs in CloudWatch for the Lambda Function and System Manager.
