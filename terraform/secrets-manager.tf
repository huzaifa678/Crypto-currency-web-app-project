resource "aws_secretsmanager_secret" "rds_credentials" {
  name = "rds-credentials-v2"
}

resource "aws_secretsmanager_secret_version" "rds_credentials_version" {
  secret_id = aws_secretsmanager_secret.rds_credentials.id
  secret_string = jsonencode({
    username = var.rds_db_username,
    password = var.rds_db_password
  })

  lifecycle {
   prevent_destroy = true
  }
}