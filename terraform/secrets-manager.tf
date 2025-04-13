resource "aws_secretsmanager_secret" "rds_credentials" {
  name  = "rds-credentials"
}

resource "aws_secretsmanager_secret_version" "rds_credentials_version" {
  secret_id     = aws_secretsmanager_secret.rds_credentials.id
  secret_string = jsonencode({
    username = var.rds_db_username
    password = var.rds_db_password
  })
}