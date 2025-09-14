data "aws_availability_zones" "available" {}

locals {
  db_creds = jsondecode(aws_secretsmanager_secret_version.production_credentials_version.secret_string)
}

locals {
  rds_db_password = coalesce(var.rds_db_password, random_password.rds_db_password.result)
}