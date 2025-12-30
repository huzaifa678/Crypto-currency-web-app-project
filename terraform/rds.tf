resource "random_password" "rds_db_password" {
  length           = 16
  special          = false
  override_special = "!#$&*"
  upper            = true  
  lower            = true  
  numeric          = true
}

resource "aws_db_subnet_group" "rds_subnet_group" {
  name       = "rds-subnet-group"
  subnet_ids = module.vpc.private_subnets

  tags = {
    Name = "rds-subnet-group"
  }
}

resource "aws_db_instance" "postgres" {
  identifier          = "crypto-db"
  engine              = "postgres"
  engine_version      = "16.6"
  instance_class      = "db.t3.micro"
  allocated_storage   = 20
  db_name             = var.rds_db_name

  username = var.rds_db_username
  password = local.rds_db_password

  db_subnet_group_name   = aws_db_subnet_group.rds_subnet_group.name
  vpc_security_group_ids = [aws_security_group.rds_sg.id]

  skip_final_snapshot = true

  tags = {
    Name = var.rds_db_name
  }
}