resource "aws_elasticache_subnet_group" "redis_subnet_group" {
  name       = "redis-subnet-group"
  subnet_ids = module.vpc.private_subnets

  tags = {
    Name = "redis-subnet-group"
  }
}


resource "aws_elasticache_replication_group" "redis_cluster" {
  replication_group_id          = "my-redis-cluster"
  description                   = "My Redis ElastiCache Cluster"
  engine                        = "redis"
  engine_version                = "7.0"
  node_type                     = "cache.t3.micro"
  num_cache_clusters            = 1 
  automatic_failover_enabled    = false
  multi_az_enabled              = false
  subnet_group_name             = aws_elasticache_subnet_group.redis_subnet_group.name
  security_group_ids            = [aws_security_group.redis_sg.id]
  port                          = 6379
  parameter_group_name          = "default.redis7"
  apply_immediately             = true
}