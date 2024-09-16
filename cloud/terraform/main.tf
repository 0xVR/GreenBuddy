provider "aws" {
  region = "us-east-2"
}

resource "aws_dynamodb_table" "plant_moisture_data" {
  name         = "PlantMoistureData"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "PlantID"
  range_key    = "Timestamp"

  attribute {
    name = "PlantID"
    type = "S"
  }

  attribute {
    name = "Timestamp"
    type = "N"
  }
}

resource "aws_sns_topic" "plant_notifications" {
  name = "PlantNotifications"
}

resource "aws_lambda_function" "plant_monitor" {
  function_name = "PlantMonitorFunction"
  role          = aws_iam_role.lambda_exec.arn
  handler       = "main"
  runtime       = "go1.x"
  filename      = "path_to_zipped_lambda_package"
  environment {
    variables = {
      SNS_TOPIC_ARN = aws_sns_topic.plant_notifications.arn
    }
  }
}

resource "aws_iam_role" "lambda_exec" {
  name = "lambda_exec_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action    = "sts:AssumeRole"
      Effect    = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })

  inline_policy {
    name   = "lambda-dynamo-sns-policy"
    policy = jsonencode({
      Version = "2012-10-17"
      Statement = [
        {
          Action   = ["dynamodb:PutItem"]
          Effect   = "Allow"
          Resource = aws_dynamodb_table.plant_moisture_data.arn
        },
        {
          Action   = ["sns:Publish"]
          Effect   = "Allow"
          Resource = aws_sns_topic.plant_notifications.arn
        }
      ]
    })
  }
}

resource "aws_apigatewayv2_api" "plant_monitor_api" {
  name          = "PlantMonitorAPI"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_integration" "lambda_integration" {
  api_id           = aws_apigatewayv2_api.plant_monitor_api.id
  integration_type = "AWS_PROXY"
  integration_uri  = aws_lambda_function.plant_monitor.invoke_arn
}

resource "aws_apigatewayv2_route" "lambda_route" {
  api_id    = aws_apigatewayv2_api.plant_monitor_api.id
  route_key = "POST /trigger"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.plant_monitor_api.id
  name        = "$default"
  auto_deploy = true
}