module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  version = "5.19.0"

  name                 = "crypto-based-web-system-vpc-network"
  cidr                 = var.vpc_cidr
  azs                  = data.aws_availability_zones.available.names 
  public_subnets       = local.ordered_public_subnets
  private_subnets      = local.ordered_private_subnets
  enable_nat_gateway   = true
  single_nat_gateway   = true
  enable_dns_hostnames = true
  enable_dns_support   = true

  public_subnet_tags = {
    "kubernetes.io/role/elb" = "1"
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
  }

  private_subnet_tags = {
    "kubernetes.io/role/internal-elb" = "1"
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
  }

  tags = {
      Name = "crypto-based-web-system-vpc-network"
  }
}