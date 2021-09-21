const AWS = require('aws-sdk');
const ssm = new AWS.SSM();
const SSM_DOCUMENT_NAME = process.env.SSM_DOCUMENT_NAME;
const S3_BUCKET = process.env.S3_BUCKET;
const SNS_TARGET = process.env.SNS_TARGET;
const BACKUP_DIRECTORY = process.env.BACKUP_DIRECTORY;
const sendCommand = (instanceId, autoScalingGroup, lifecycleHook) => {
  var params = {
    DocumentName: SSM_DOCUMENT_NAME,
    InstanceIds: [instanceId],
    Parameters: {
      'ASGNAME': [autoScalingGroup],
      'LIFECYCLEHOOKNAME': [lifecycleHook],
      'BACKUPDIRECTORY': [BACKUP_DIRECTORY],
      'S3BUCKET': [S3_BUCKET],
      'SNSTARGET': [SNS_TARGET],
    },
    TimeoutSeconds: 300
  };
  return ssm.sendCommand(params).promise();
}
exports.handler = async (event) => {
  console.log('Received event ', JSON.stringify(event));
  try {
    const records = event.Records;
    if (!records || !records.length) {
      return;
    }
    for (const record of records) {
      if (record.EventSource !== 'aws:sns') {
        console.log('Record is not processed because record.EventSource is not aws:sns');
        continue;
      }
      const message = JSON.parse(record.Sns.Message);
      if (message.LifecycleTransition !== 'autoscaling:EC2_INSTANCE_TERMINATING') {
        console.log('Record is not processed because message.LifecycleTransition is not autoscaling:EC2_INSTANCE_TERMINATING');
        continue;
      }
      console.log("processing autoscaling event");
      const autoScalingGroup = message.AutoScalingGroupName;
      const instanceId = message.EC2InstanceId;
      const lifecycleHook = message.LifecycleHookName;
      await sendCommand(instanceId, autoScalingGroup, lifecycleHook);
      console.log('sent command');
    }
  } catch (error) {
    throw error;
  }
}