data "aws_iam_role" "lambda" {
  name = "${var.project_name}-lambda-role"
}

resource "aws_lambda_function" "api" {
  function_name = "${var.project_name}-api"

  role = data.aws_iam_role.lambda.arn

  runtime = "provided.al2023"
  handler = "bootstrap"

  filename         = var.lambda_zip
  source_code_hash = filebase64sha256(var.lambda_zip)

  timeout     = 30
  memory_size = 256
}