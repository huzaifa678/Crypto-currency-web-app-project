resource "aws_secretsmanager_secret" "rds_credentials" {
  name = "rds-credentials"
}

resource "aws_secretsmanager_secret_version" "rds_credentials_version" {
  secret_id     = aws_secretsmanager_secret.rds_credentials.id
  secret_string = jsonencode({
    username = var.rds_db_username
    password = var.rds_db_password
  })
}

resource "aws_db_instance" "postgres" {
    identifier = "crypto-db"
    engine     = "postgres"
    engine_version    = "12.15"
    instance_class    = "db.t3.micro"
    allocated_storage = 20

    db_name = var.rds_db_name
    username = var.rds_db_username
    password = var.rds_db_password

    skip_final_snapshot = true

    tags = {
        Name = var.rds_db_name
    }
}