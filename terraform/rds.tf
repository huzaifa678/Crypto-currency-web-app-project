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

resource "aws_iam_role" "rds_monitoring" {
  name = "crypto-rds-monitoring-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Action    = "sts:AssumeRole"
      Principal = { Service = "monitoring.rds.amazonaws.com" }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "rds_monitoring" {
  role       = aws_iam_role.rds_monitoring.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonRDSEnhancedMonitoringRole"
}

resource "aws_db_parameter_group" "postgres" {
  name   = "crypto-db-postgres16"
  family = "postgres16"

  parameter {
    name  = "log_statement"
    value = "all"
  }

  parameter {
    name  = "log_min_duration_statement"
    value = "1"
  }

  # Enforce TLS for all client connections (encryption in transit).
  parameter {
    name  = "rds.force_ssl"
    value = "1"
  }

  tags = {
    Name = "crypto-db-postgres16"
  }
}

resource "aws_db_instance" "postgres" {
  identifier        = "crypto-db"
  engine            = "postgres"
  engine_version    = "16.6"
  instance_class    = "db.t3.micro"
  allocated_storage = 20
  db_name           = var.rds_db_name

  username = var.rds_db_username
  password = local.rds_db_password

  db_subnet_group_name   = aws_db_subnet_group.rds_subnet_group.name
  vpc_security_group_ids = [module.eks.rds_sg_id]
  parameter_group_name   = aws_db_parameter_group.postgres.name

  storage_encrypted = true
  kms_key_id        = aws_kms_key.rds.arn

  multi_az                  = true
  deletion_protection       = true
  copy_tags_to_snapshot     = true
  backup_retention_period   = 7
  skip_final_snapshot       = false
  final_snapshot_identifier = "crypto-db-final"

  auto_minor_version_upgrade          = true
  iam_database_authentication_enabled = true

  enabled_cloudwatch_logs_exports = ["postgresql", "upgrade"]
  monitoring_interval             = 60
  monitoring_role_arn             = aws_iam_role.rds_monitoring.arn
  performance_insights_enabled    = true
  performance_insights_kms_key_id = aws_kms_key.rds.arn

  tags = {
    Name = var.rds_db_name
  }
}