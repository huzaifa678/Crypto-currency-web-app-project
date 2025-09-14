output "rds_db" {
  description = "RDS DB"
  value       = aws_db_instance.postgres.db_name
}

output "rds_username" {
  description = "Username"
  sensitive = true
  value = local.db_creds
}

output "vpc_id" {
  description = "The VPC ID"
  value       = module.vpc.vpc_id
}

output "eks_cluster_name" {
  description = "EKS Cluster name"
  value       = aws_eks_cluster.eks_cluster.name

}

# output "ecr_repository_url" {
#   description = "ECR repository URL"
#   value       = aws_ecr_repository.ecr_repo.repository_url
# }

output "redis_endpoint" {
  description = "Redis Endpoint"
  value       = aws_elasticache_replication_group.redis_cluster.primary_endpoint_address
}