resource "aws_elasticache_subnet_group" "redis_subnet_group" {
  name       = "redis-subnet-group"
  subnet_ids = module.vpc.private_subnets

  tags = {
    Name = "redis-subnet-group"
  }
}

resource "random_password" "redis_auth" {
  length           = 32
  special          = true
  override_special = "!&#$^<>-"
}

resource "aws_elasticache_replication_group" "redis_cluster" {
  replication_group_id       = "my-redis-cluster"
  description                = "My Redis ElastiCache Cluster"
  engine                     = "redis"
  engine_version             = "7.0"
  node_type                  = "cache.t3.micro"
  num_node_groups            = 1
  replicas_per_node_group    = 1
  automatic_failover_enabled = true
  multi_az_enabled           = true
  subnet_group_name          = aws_elasticache_subnet_group.redis_subnet_group.name
  security_group_ids         = [module.eks.redis_sg_id]
  port                       = 6379
  parameter_group_name       = "default.redis7"
  apply_immediately          = true

  at_rest_encryption_enabled = true
  transit_encryption_enabled = true
  kms_key_id                 = aws_kms_key.elasticache.arn
  auth_token                 = random_password.redis_auth.result
}