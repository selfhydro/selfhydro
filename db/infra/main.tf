provider "aws" {
  region = "us-east-2"
}

resource "aws_cloudwatch_event_rule" "once_a_day" {
    name = "once_a_day"
    description = "Fires every 24 hours"
    schedule_expression = "cron(45 23 * * ? *)"
}

resource "aws_cloudwatch_event_target" "check_foo_every_five_minutes" {
    rule = "${aws_cloudwatch_event_rule.every_five_minutes.name}"
    target_id = "check_foo"
    arn = "${aws_lambda_function.check_foo.arn}"
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_check_foo" {
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
  filename      = "create_dynamo_db_tables_payload.zip"
  function_name = "create_dynamo_db_tables"
  role          = "${aws_iam_role.iam_for_lambda.arn}"
  handler       = "dynamoDBTableCreater.CreateTable"
  runtime = "go1.12"

  environment {
    variables = {
      type = "CREATE_NEW_TABLE"
    }
  }
}
