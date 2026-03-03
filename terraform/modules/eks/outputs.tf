output "eks_cluster_name" {
  value = aws_eks_cluster.eks_cluster.name
}

output "eks_cluster_endpoint" {
  value = aws_eks_cluster.eks_cluster.endpoint
}

output "eks_cluster_ca" {
  value = aws_eks_cluster.eks_cluster.certificate_authority[0].data
}

output "eks_service_ipv4_cidr" {
  value = aws_eks_cluster.eks_cluster.kubernetes_network_config[0].service_ipv4_cidr
}

output "launch_template_user_data" {
  value = aws_launch_template.eks_nodes.user_data
}

output "cert_manager_irsa_role_arn" {
  value = aws_iam_role.cert_manager_irsa_role.arn
}

output "eks_node_group" {
  value = aws_eks_node_group.eks_node_group
}

output "external_dns_irsa_role_arn" {
  value = aws_iam_role.external_dns_irsa_role.arn
}

output "aws_lb_controller_irsa_role_arn" {
  value = aws_iam_role.aws_lb_controller_irsa_role.arn
}

output "rds_sg_id" {
  value = aws_security_group.rds_sg.id
}

output "redis_sg_id" {
  value = aws_security_group.redis_sg.id
}