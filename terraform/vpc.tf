module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  version = "5.19.0"

  name                 = "crypto-based-web-system-vpc-network"
  cidr                 = var.vpc_cidr
  azs                  = data.aws_availability_zones.available.names 
  private_subnets      = var.private_subnets
  public_subnets       = var.public_subnets
  enable_nat_gateway   = true
  single_nat_gateway   = true
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
      Name = "crypto-based-web-system-vpc-network"
  }
}