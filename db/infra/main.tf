provider "aws" {
  region = "ap-southeast-2"
}

resource "aws_cloudwatch_event_rule" "once_a_day" {
    name = "once_a_day"
    description = "Fires every 24 hours"
    schedule_expression = "cron(45 23 * * ? *)"
}

resource "aws_cloudwatch_event_target" "create_dynamo_db_tables_once_a_day" {
    rule = "${aws_cloudwatch_event_rule.once_a_day.name}"
    target_id = "create_dynamo_db_tables"
    arn = "${aws_lambda_function.create_dynamo_db_tables.arn}"
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_create_dynamo_db_tables" {
    statement_id = "AllowExecutionFromCloudWatch"
    action = "lambda:InvokeFunction"
    function_name = "${aws_lambda_function.create_dynamo_db_tables.function_name}"
    principal = "events.amazonaws.com"
    source_arn = "${aws_cloudwatch_event_rule.once_a_day.arn}"
}

resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_lambda"

  assume_role_policy = <<EOF
{
	"Version": "2012-10-17",
	"Statement": [{
			"Effect": "Allow",
			"Action": [
				"dynamodb:BatchGetItem",
				"dynamodb:GetItem",
				"dynamodb:Query",
				"dynamodb:Scan",
				"dynamodb:BatchWriteItem",
				"dynamodb:PutItem",
				"dynamodb:UpdateItem"
			],
			"Resource": "*"
		}
	]
}
EOF
}

resource "aws_lambda_function" "create_dynamo_db_tables" {
  s3_bucket     = "selfhydro-releases"
  s3_key        = "selfhydro-state-db/selfhydro-state-db-release-${var.version}.tar"
  filename      = "selfhydro-state-db-release-${var.version}.zip"
  function_name = "selfhydroStateTableCreater"
  role          = "${aws_iam_role.iam_for_lambda.arn}"
  handler       = "dynamoDBTableCreater.CreateTable"
  runtime = "go1.12"

  environment {
    variables = {
      type = "CREATE_NEW_TABLE"
    }
  }
}
