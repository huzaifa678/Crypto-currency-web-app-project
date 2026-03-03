output "eks_cluster_name" {
  value = module.eks.eks_cluster_name
}

output "eks_cluster_endpoint" {
  value = module.eks.eks_cluster_endpoint
}

output "eks_cluster_ca" {
  value = module.eks.eks_cluster_ca
}

output "eks_service_ipv4_cidr" {
  value = module.eks.eks_service_ipv4_cidr
}

output "launch_template_user_data" {
  value = module.eks.launch_template_user_data
}

output "cert_manager_irsa_role_arn" {
  value = module.eks.cert_manager_irsa_role_arn
}

output "external_dns_irsa_role_arn" {
  value = module.eks.external_dns_irsa_role_arn
}

output "aws_lb_controller_irsa_role_arn" {
  value = module.eks.aws_lb_controller_irsa_role_arn
}

output "rds_sg_id" {
  value = module.eks.rds_sg_id
}

output "redis_sg_id" {
  value = module.eks.redis_sg_id
}