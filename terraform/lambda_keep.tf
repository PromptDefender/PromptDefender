resource "aws_iam_role" "lambda_role_keep" {
  name = "${terraform.workspace}-lambda_role_keep"

  assume_role_policy = jsonencode({
    Statement = [
      {
        Action    = "sts:AssumeRole",
        Effect    = "Allow",
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      },
    ],
    Version = "2012-10-17"
  })
}

resource "aws_cloudwatch_log_group" "lambda_log_group_keep" { #tfsec:ignore:aws-cloudwatch-log-group-customer-key
  name              = "/aws/lambda/${aws_lambda_function.aws_Lambda_keep.function_name}"
  retention_in_days = 14
}

resource "aws_iam_policy" "lambda_logging_policy" {
  name   = "${terraform.workspace}-lambda_logging_policy_keep"
  policy = jsonencode({
    Version   = "2012-10-17",
    Statement = [
      {
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        Effect   = "Allow",
        Resource = "${aws_cloudwatch_log_group.lambda_log_group_keep.arn}:*"
      },
    ],
  })
}

resource "aws_iam_role_policy_attachment" "xray_policy_attachment_keep" {
  role       = aws_iam_role.lambda_role_keep.name
  policy_arn = "arn:aws:iam::aws:policy/AWSXRayDaemonWriteAccess"
}

resource "aws_iam_role_policy_attachment" "lambda_logging_attach" {
  role       = aws_iam_role.lambda_role_keep.name
  policy_arn = aws_iam_policy.lambda_logging_policy.arn
}

resource "aws_iam_role_policy_attachment" "dynamodb_read_write_policy_attachment_keep" {
  role       = aws_iam_role.lambda_role_keep.name
  policy_arn = aws_iam_policy.dynamodb_read_write_policy_wall.arn
}

resource "aws_lambda_function" "aws_Lambda_keep" {
  function_name    = "${terraform.workspace}-PromptDefender-Keep"
  handler          = "bootstrap"
  role             = aws_iam_role.lambda_role_keep.arn
  filename         = data.archive_file.lambda_keep_zip.output_path
  runtime          = "provided.al2"
  source_code_hash = data.archive_file.lambda_keep_zip.output_base64sha256
  timeout          = 60

  tracing_config {
    mode = "Active"
  }


  environment {
    variables = {
      open_ai_api_key = var.openai_secret_key
      version         = var.commit_version
      CACHE_TABLE_NAME             = aws_dynamodb_table.cache_table.name
    }
  }
}

data "archive_file" "lambda_keep_zip" {
  type        = "zip"
  source_file = var.lambda_keep_path
  output_path = "keep_function.zip"
}

variable "lambda_keep_path" {
  type    = string
  default = "../cmd/lambda_keep/bootstrap"
}