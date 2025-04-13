resource "aws_secretsmanager_secret" "rds_credentials" {
  count = length(data.aws_secretsmanager_secret.existing_secret) > 0 ? 0 : 1
  name = "rds-credentials"
}

resource "aws_secretsmanager_secret_version" "rds_credentials_version" {
  secret_id     = aws_secretsmanager_secret.rds_credentials.id
  secret_string = jsonencode({
    username = var.rds_db_username
    password = var.rds_db_password
  })
}