resource "aws_iam_role" "sdm_ecs_task_execution_role" {
  name = "sdm-ecs-role-name-${random_id.id.hex}"

  assume_role_policy = jsonencode(
    {
      "Version" : "2012-10-17",
      "Statement" : [
        {
          "Action" : "sts:AssumeRole",
          "Principal" : {
            "Service" : "ecs-tasks.amazonaws.com"
          },
          "Effect" : "Allow",
          "Sid" : ""
        }
      ]
  })
}

resource "aws_iam_role" "sdm_ecs_task_role" {
  name = "sdm-ecs-role-name-task-${random_id.id.hex}"

  assume_role_policy = jsonencode(
    {
      "Version" : "2012-10-17",
      "Statement" : [
        {
          "Action" : "sts:AssumeRole",
          "Principal" : {
            "Service" : "ecs-tasks.amazonaws.com"
          },
          "Effect" : "Allow",
          "Sid" : ""
        }
      ]
  })
}

resource "aws_iam_policy" "ecs-policy" {
  name        = "sdm-ecs-policy-${random_id.id.hex}"
  path        = "/"
  description = "sdm-ecs-policy"
  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Action" : [
          "ecr:GetAuthorizationToken",
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
          "logs:CreateLogStream",
          "logs:PutLogEvents",
          "secretsmanager:*",
          "cloudformation:CreateChangeSet",
          "cloudformation:DescribeChangeSet",
          "cloudformation:DescribeStackResource",
          "cloudformation:DescribeStacks",
          "cloudformation:ExecuteChangeSet",
          "ec2:DescribeSecurityGroups",
          "ec2:DescribeSubnets",
          "ec2:DescribeVpcs",
          "kms:*",
          "lambda:ListFunctions",
          "rds:DescribeDBClusters",
          "rds:DescribeDBInstances",
          "redshift:DescribeClusters",
          "tag:GetResources"
        ],
        "Effect" : "Allow",
        "Resource" : "*"
      },
      {
        "Action" : [
          "lambda:AddPermission",
          "lambda:CreateFunction",
          "lambda:GetFunction",
          "lambda:InvokeFunction",
          "lambda:UpdateFunctionConfiguration"
        ],
        "Effect" : "Allow",
        "Resource" : "arn:aws:lambda:*:*:function:SecretsManager*"
      },
      {
        "Action" : [
          "serverlessrepo:CreateCloudFormationChangeSet",
          "serverlessrepo:GetApplication"
        ],
        "Effect" : "Allow",
        "Resource" : "arn:aws:serverlessrepo:*:*:applications/SecretsManager*"
      },
      {
        "Action" : [
          "s3:GetObject"
        ],
        "Effect" : "Allow",
        "Resource" : [
          "arn:aws:s3:::awsserverlessrepo-changesets*",
          "arn:aws:s3:::secrets-manager-rotation-apps-*/*"
        ]
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "sdm-ecs-task-execution-role-policy-attachment" {
  role       = aws_iam_role.sdm_ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role_policy_attachment" "task_secrets" {
  role       = aws_iam_role.sdm_ecs_task_role.name
  policy_arn = aws_iam_policy.ecs-policy.arn
}