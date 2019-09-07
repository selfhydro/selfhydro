provider "aws" {
  region = "${var.region}"
}

data "aws_caller_identity" "current" {}

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

data "aws_iam_policy_document" "iam_for_lambda" {
  statement {
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]

    resources = [
      "arn:aws:logs:${var.region}:${data.aws_caller_identity.current.account_id}:*"
    ]
  }

  statement {
    actions = [
      "sts:assumeRole"
    ]

    principal = [
      type = "Service"
      identifiers = [*]
    ]
  }
  statement {
    actions = [
      "dynamodb:*"
    ]

    resources = [
      "arn:aws:dynamodb:${var.region}:${data.aws_caller_identity.current.account_id}:*"
    ]
  }
}

resource "aws_iam_role" "iam_for_lambda" {
  name    = "iam_for_lambda"
  assume_role_policy  = "${data.aws_iam_policy_document.iam_for_lambda.json}"
}

resource "aws_lambda_function" "create_dynamo_db_tables" {
  s3_bucket     = "selfhydro-releases"
  s3_key        = "selfhydro-state-db/selfhydro-state-db-release-${var.lamdba-version}.tar"
  function_name = "selfhydroStateTableCreater"
  role          = "${aws_iam_role.iam_for_lambda.arn}"
  handler       = "dynamoDBTableCreater.CreateTable"
  runtime       = "go1.x"

  environment {
    variables = {
      type = "CREATE_NEW_TABLE"
    }
  }
}
