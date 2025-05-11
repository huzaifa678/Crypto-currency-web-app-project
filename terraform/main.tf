data "aws_availability_zones" "available" {}

locals {
  db_creds = jsondecode(
    aws_secretsmanager_secret_version.rds_credentials_version.secret_string
  )
}

