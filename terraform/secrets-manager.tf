resource "aws_secretsmanager_secret" "rds_credentials" {
  count = var.existing_secret ? 1 : 0
  name  = "rds-credentials-v1"
}

resource "aws_secretsmanager_secret_version" "rds_credentials_version" {
  count = var.existing_secret ? 1 : 0
  secret_id     = aws_secretsmanager_secret.rds_credentials[0].id
  secret_string = jsonencode({
    username = var.rds_db_username
    password = var.rds_db_password
  })
}