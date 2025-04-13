resource "aws_db_instance" "postgres" {

  identifier          = "crypto-db"
  engine              = "postgres"
  engine_version      = "16.2"
  instance_class      = "db.t3.micro"
  allocated_storage   = 20
  db_name             = var.rds_db_name

  username = jsondecode(data.aws_secretsmanager_secret_version.rds_credentials.secret_string)["username"]
  password = jsondecode(data.aws_secretsmanager_secret_version.rds_credentials.secret_string)["password"]

  skip_final_snapshot = true

  tags = {
    Name = var.rds_db_name
  }
}
