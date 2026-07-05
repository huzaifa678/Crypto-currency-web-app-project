resource "aws_security_group" "eks_nodes" {
  name        = "eks-nodes-sg"
  description = "Security group for EKS worker nodes"
  vpc_id      = var.vpc_id

  tags = {
    Name = "eks-nodes-sg"
  }
}

resource "aws_security_group" "rds_sg" {
  name        = "rds-postgres-sg"
  description = "Security group for the RDS PostgreSQL instance"
  vpc_id      = var.vpc_id

  tags = {
    Name = "rds-postgres-sg"
  }
}

resource "aws_security_group" "redis_sg" {
  name        = "redis-sg"
  description = "Security group for the ElastiCache Redis cluster"
  vpc_id      = var.vpc_id

  tags = {
    Name = "redis-sg"
  }
}

# Scoped egress for worker nodes: HTTPS (AWS APIs, ECR, S3) and DNS only.
resource "aws_security_group_rule" "eks_nodes_egress_dns_tcp" {
  description       = "Allow nodes DNS resolution over TCP"
  type              = "egress"
  from_port         = 53
  to_port           = 53
  protocol          = "tcp"
  security_group_id = aws_security_group.eks_nodes.id
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "eks_nodes_egress_dns_udp" {
  description       = "Allow nodes DNS resolution over UDP"
  type              = "egress"
  from_port         = 53
  to_port           = 53
  protocol          = "udp"
  security_group_id = aws_security_group.eks_nodes.id
  cidr_blocks       = ["0.0.0.0/0"]
}


resource "aws_security_group_rule" "allow_nlb_to_nodes_https" {
  count                    = var.environment == "post-test" ? 1 : 0
  description              = "Allow NLB health checks / traffic to nodes on HTTPS"
  type                     = "ingress"
  from_port                = 443
  to_port                  = 443
  protocol                 = "tcp"
  security_group_id        = aws_security_group.eks_nodes.id
  source_security_group_id = aws_security_group.eks_nodes.id
}

resource "aws_security_group_rule" "eks_nodes_egress_https" {
  description       = "Allow nodes HTTPS egress to AWS APIs, ECR and S3"
  type              = "egress"
  from_port         = 443
  to_port           = 443
  protocol          = "tcp"
  security_group_id = aws_security_group.eks_nodes.id
  cidr_blocks       = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "allow_node_to_core_dns_tcp" {
  description              = "Allow nodes to reach CoreDNS on the control plane (TCP)"
  type                     = "ingress"
  from_port                = 53
  to_port                  = 53
  protocol                 = "tcp"
  security_group_id        = aws_eks_cluster.eks_cluster.vpc_config[0].cluster_security_group_id
  source_security_group_id = aws_security_group.eks_nodes.id

  depends_on = [aws_eks_cluster.eks_cluster]
}

resource "aws_security_group_rule" "allow_node_to_core_dns_udp" {
  description              = "Allow nodes to reach CoreDNS on the control plane (UDP)"
  type                     = "ingress"
  from_port                = 53
  to_port                  = 53
  protocol                 = "udp"
  security_group_id        = aws_eks_cluster.eks_cluster.vpc_config[0].cluster_security_group_id
  source_security_group_id = aws_security_group.eks_nodes.id

  depends_on = [aws_eks_cluster.eks_cluster]
}

resource "aws_security_group_rule" "allow_node_to_node_dns_tcp" {
  description              = "Allow node-to-node DNS (TCP)"
  type                     = "ingress"
  from_port                = 53
  to_port                  = 53
  protocol                 = "tcp"
  security_group_id        = aws_security_group.eks_nodes.id
  source_security_group_id = aws_security_group.eks_nodes.id
}

resource "aws_security_group_rule" "allow_node_to_node_dns_udp" {
  description              = "Allow node-to-node DNS (UDP)"
  type                     = "ingress"
  from_port                = 53
  to_port                  = 53
  protocol                 = "udp"
  security_group_id        = aws_security_group.eks_nodes.id
  source_security_group_id = aws_security_group.eks_nodes.id
}

resource "aws_security_group_rule" "allow_node_to_control_plane" {
  description              = "Allow nodes to reach the control plane API (HTTPS)"
  type                     = "ingress"
  from_port                = 443
  to_port                  = 443
  protocol                 = "tcp"
  security_group_id        = aws_eks_cluster.eks_cluster.vpc_config[0].cluster_security_group_id
  source_security_group_id = aws_security_group.eks_nodes.id

  depends_on = [aws_eks_cluster.eks_cluster]
}

resource "aws_security_group_rule" "allow_control_plane_to_nodes" {
  description              = "Allow control plane to reach nodes (HTTPS)"
  type                     = "ingress"
  from_port                = 443
  to_port                  = 443
  protocol                 = "tcp"
  security_group_id        = aws_security_group.eks_nodes.id
  source_security_group_id = aws_eks_cluster.eks_cluster.vpc_config[0].cluster_security_group_id

  depends_on = [aws_eks_cluster.eks_cluster]
}

resource "aws_security_group_rule" "allow_eks_to_rds" {
  description              = "Allow EKS nodes to reach RDS PostgreSQL"
  type                     = "ingress"
  from_port                = 5432
  to_port                  = 5432
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.eks_nodes.id
  security_group_id        = aws_security_group.rds_sg.id
}

resource "aws_security_group_rule" "allow_eks_to_redis" {
  description              = "Allow EKS nodes to reach ElastiCache Redis"
  type                     = "ingress"
  from_port                = 6379
  to_port                  = 6379
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.eks_nodes.id
  security_group_id        = aws_security_group.redis_sg.id
}