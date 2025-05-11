resource "aws_db_instance" "postgres" {

  identifier          = "crypto-db"
  engine              = "postgres"
  engine_version      = "16.3"
  instance_class      = "db.t3.micro"
  allocated_storage   = 20
  db_name             = var.rds_db_name

  username = local.db_creds.username
  password = local.db_creds.password

  skip_final_snapshot = true

  tags = {
    Name = var.rds_db_name
  }

  lifecycle {
   prevent_destroy = true
  }
}