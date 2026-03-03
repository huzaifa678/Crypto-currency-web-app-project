data "aws_availability_zones" "available" {}

locals {
  db_creds = jsondecode(aws_secretsmanager_secret_version.production_credentials_version.secret_string)
}

locals {
  rds_db_password = coalesce(var.rds_db_password, random_password.rds_db_password.result)
}

locals {
  ordered_public_subnets  = sort(var.public_subnets)   
  ordered_private_subnets = sort(var.private_subnets)
}

module "eks" {
  source = "./modules/eks"

  cluster_name                          = var.cluster_name
  kubernetes_version                     = var.kubernetes_version
  private_subnets                        = module.vpc.private_subnets
  cluster_endpoint_public_access_cidrs   = var.cluster_endpoint_public_access_cidrs
  vpc_id                                 = module.vpc.vpc_id
  environment                            = var.environment
  region                                 = var.region
}

module "k8s" {
  source = "./modules/k8s-and-helm"

  cluster_name              = module.eks.eks_cluster_name
  cluster_endpoint          = module.eks.eks_cluster_endpoint
  cluster_ca                = module.eks.eks_cluster_ca
  cert_manager_irsa_role_arn = module.eks.cert_manager_irsa_role_arn
  external_dns_irsa_role_arn = module.eks.external_dns_irsa_role_arn
  aws_lb_controller_irsa_role_arn = module.eks.aws_lb_controller_irsa_role_arn
  vpc_id = module.vpc.vpc_id
  environment = var.environment
  region = var.region
  eks_node_group = module.eks.eks_node_group
  vpc = module.vpc
}