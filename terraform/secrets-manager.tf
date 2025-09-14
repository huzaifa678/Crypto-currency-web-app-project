resource "aws_secretsmanager_secret" "production_credentials" {
  name = "production-credentials-v1"
}

data "external" "app_env" {
  program = ["bash", "${path.module}/parse_env.sh"]
}

resource "aws_secretsmanager_secret_version" "production_credentials_version" {
  secret_id = aws_secretsmanager_secret.production_credentials.id
  secret_string = jsonencode({
    USERNAME              = var.rds_db_username
    PASSWORD              = local.rds_db_password
    DB_SOURCE             = "postgresql://${var.rds_db_username}:${local.rds_db_password}@${aws_db_instance.postgres.endpoint}/${var.rds_db_name}?sslmode=disable"
    HTTP_SERVER_ADDR      = data.external.app_env.result.HTTP_SERVER_ADDR
    GRPC_SERVER_ADDR      = data.external.app_env.result.GRPC_SERVER_ADDR
    REDIS_ADDR            = data.external.app_env.result.REDIS_ADDR
    MIGRATION_URL         = data.external.app_env.result.MIGRATION_URL
    TOKEN_SYMMETRIC_KEY   = var.token_symmetric_key
    ACCESS_TOKEN_EXPIRE   = data.external.app_env.result.ACCESS_TOKEN_DURATION
    REFRESH_TOKEN_EXPIRE  = data.external.app_env.result.REFRESH_TOKEN_DURATION
    EMAIL_SENDER_NAME     = data.external.app_env.result.SENDER_NAME
    EMAIL_SENDER_ADDRESS  = data.external.app_env.result.SENDER_EMAIL
    EMAIL_SENDER_PASSWORD = data.external.app_env.result.SENDER_PASSWORD
    GOOGLE_CLIENT_ID      = var.google_client_id
    GOOGLE_CLIENT_SECRET  = var.google_client_secret
    REDIRECT_URL          = var.redirect_url
  })

  depends_on = [
    aws_db_instance.postgres
  ]
}