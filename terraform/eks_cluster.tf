resource "aws_eks_cluster" "eks_cluster" {
  name     = var.cluster_name
  role_arn = aws_iam_role.eks_cluster_role.arn
  version  = var.kubernetes_version

  vpc_config {
    subnet_ids = module.vpc.private_subnets
    endpoint_private_access = true
    endpoint_public_access  = true
    public_access_cidrs     = var.cluster_endpoint_public_access_cidrs
  }

  access_config {
    authentication_mode                         = "API"
    bootstrap_cluster_creator_admin_permissions = true
  }

  tags = {
    cluster = "demo"
  }
}

data "aws_eks_cluster" "this" {
  name = aws_eks_cluster.eks_cluster.name
}

resource "aws_iam_openid_connect_provider" "eks" {
  url             = data.aws_eks_cluster.this.identity[0].oidc[0].issuer
  client_id_list  = ["sts.amazonaws.com"]
  thumbprint_list = ["9e99a48a9960b14926bb7f3b02e22da0ecd4e0a4"] 
}

resource "aws_security_group" "eks_nodes" {
  name        = "eks-nodes-sg"
  description = "Security group for EKS worker nodes"
  vpc_id      = module.vpc.vpc_id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "eks-nodes-sg"
  }
}

resource "aws_security_group" "rds_sg" {
  name   = "rds-postgres-sg"
  vpc_id = module.vpc.vpc_id

  tags = {
    Name = "rds-postgres-sg"
  }
}

resource "aws_security_group" "redis_sg" {
  name   = "redis-sg"
  vpc_id = module.vpc.vpc_id

  tags = {
    Name = "redis-sg"
  }
}

resource "aws_security_group_rule" "allow_node_to_core_dns_tcp" {
  type                     = "ingress"
  from_port                = 53
  to_port                  = 53
  protocol                 = "tcp"
  security_group_id        = aws_eks_cluster.eks_cluster.vpc_config[0].cluster_security_group_id
  source_security_group_id = aws_security_group.eks_nodes.id

  depends_on = [aws_eks_cluster.eks_cluster]
}

resource "aws_security_group_rule" "allow_node_to_core_dns_udp" {
  type                     = "ingress"
  from_port                = 53
  to_port                  = 53
  protocol                 = "udp"
  security_group_id        = aws_eks_cluster.eks_cluster.vpc_config[0].cluster_security_group_id
  source_security_group_id = aws_security_group.eks_nodes.id

  depends_on = [aws_eks_cluster.eks_cluster]
}

resource "aws_security_group_rule" "allow_node_to_control_plane" {
  type                     = "ingress"
  from_port                = 443
  to_port                  = 443
  protocol                 = "tcp"
  security_group_id        = aws_eks_cluster.eks_cluster.vpc_config[0].cluster_security_group_id
  source_security_group_id = aws_security_group.eks_nodes.id

  depends_on = [aws_eks_cluster.eks_cluster]
}

resource "aws_security_group_rule" "allow_control_plane_to_nodes" {
  type                     = "ingress"
  from_port                = 443
  to_port                  = 443
  protocol                 = "tcp"
  security_group_id        = aws_security_group.eks_nodes.id
  source_security_group_id = aws_eks_cluster.eks_cluster.vpc_config[0].cluster_security_group_id

  depends_on = [aws_eks_cluster.eks_cluster]
}

resource "aws_security_group_rule" "allow_eks_to_rds" {
  type                     = "ingress"
  from_port                = 5432
  to_port                  = 5432
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.eks_nodes.id
  security_group_id        = aws_security_group.rds_sg.id
}

resource "aws_security_group_rule" "allow_eks_to_redis" {
  type                     = "ingress"
  from_port                = 6379
  to_port                  = 6379
  protocol                 = "tcp"
  source_security_group_id = aws_security_group.eks_nodes.id
  security_group_id        = aws_security_group.redis_sg.id
}

resource "aws_iam_role" "eks_cluster_role" {
  name = "eks-cluster-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action    = "sts:AssumeRole"
        Effect    = "Allow"
        Principal = {
          Service = "eks.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "eks_cluster_creator_admin" {
  policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess"
  role       = aws_iam_role.eks_cluster_role.name
}

data "aws_ssm_parameter" "eks_al2023_ami" {
  name = "/aws/service/eks/optimized-ami/1.32/amazon-linux-2023/x86_64/standard/recommended/image_id"
}

resource "aws_launch_template" "eks_nodes" {
  name_prefix   = "eks-nodes-lt"
  instance_type = "t3.small"

  vpc_security_group_ids = [
    aws_security_group.eks_nodes.id,
    aws_eks_cluster.eks_cluster.vpc_config[0].cluster_security_group_id
  ]

  image_id = data.aws_ssm_parameter.eks_al2023_ami.value

  user_data = base64encode(templatefile("${path.module}/user_data.tpl", {
    cluster_name     = aws_eks_cluster.eks_cluster.name
    cluster_endpoint = data.aws_eks_cluster.this.endpoint
    cluster_ca       = data.aws_eks_cluster.this.certificate_authority[0].data
    cidr             = "10.100.0.0/16"
  }))

  tag_specifications {
    resource_type = "instance"
    tags = {
      Name = "eks-node"
    }
  }
}

resource "aws_eks_node_group" "eks_node_group" {
  cluster_name    = aws_eks_cluster.eks_cluster.name
  node_group_name = "node-group"
  node_role_arn   = aws_iam_role.eks_node_role.arn
  subnet_ids      = module.vpc.private_subnets

  scaling_config {
    desired_size = 2
    min_size     = 0
    max_size     = 2
  }

  ami_type = "CUSTOM"

  launch_template {
    id      = aws_launch_template.eks_nodes.id
    version = aws_launch_template.eks_nodes.latest_version
  }

  depends_on = [aws_security_group_rule.allow_node_to_control_plane]
}

resource "aws_iam_role" "eks_node_role" {
  name = "eks-node-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action    = "sts:AssumeRole"
        Effect    = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy" "eks_cluster_creator_admin_policy" {
  name = "eks-cluster-creator-admin-policy"
  role = aws_iam_role.eks_cluster_role.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "eks:DescribeCluster"
        Effect = "Allow"
        Resource = "*"
      },
      {
        Action = "eks:ListClusters"
        Effect = "Allow"
        Resource = "*"
      },
      {
        Action = "eks:CreateCluster"
        Effect = "Allow"
        Resource = "*"
      },
      {
        Action = "eks:UpdateClusterConfig"
        Effect = "Allow"
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "eks_node_ssm" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
  role       = aws_iam_role.eks_node_role.name
}

resource "aws_iam_role_policy_attachment" "eks_node_policy_attachment" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  role       = aws_iam_role.eks_node_role.name
}

resource "aws_iam_role_policy_attachment" "eks_node_ecr_readonly" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  role       = aws_iam_role.eks_node_role.name
}

resource "aws_iam_role_policy_attachment" "eks_node_cni_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  role       = aws_iam_role.eks_node_role.name
}

resource "aws_iam_policy" "custom_route53_policy" {
  name        = "eks-node-custom-route53"
  description = "Allow nodes to manage Route53 records"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "route53:ChangeResourceRecordSets",
          "route53:ListHostedZones",
          "route53:ListResourceRecordSets",
          "route53:GetChange"
        ],
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role" "cert_manager_irsa_role" {
  name = "cert-manager-irsa-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Federated = aws_iam_openid_connect_provider.eks.arn
        },
        Action = "sts:AssumeRoleWithWebIdentity",
        Condition = {
          StringEquals = {
            "${replace(data.aws_eks_cluster.this.identity[0].oidc[0].issuer, "https://", "")}:sub" = "system:serviceaccount:cert-manager:cert-manager"
          }
        }
      }
    ]
  })
}

resource "aws_iam_role" "external_dns_irsa_role" {
  name = "external-dns-irsa-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Federated = aws_iam_openid_connect_provider.eks.arn
        }
        Action = "sts:AssumeRoleWithWebIdentity"
        Condition = {
          StringEquals = {
            "${replace(data.aws_eks_cluster.this.identity[0].oidc[0].issuer, "https://", "")}:sub" = "system:serviceaccount:external-dns:external-dns"
          }
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "eks_node_custom_route53_policy_attachment" {
  policy_arn = aws_iam_policy.custom_route53_policy.arn
  role       = aws_iam_role.eks_node_role.name
}

resource "aws_iam_role_policy_attachment" "cert_manager_irsa_route53_attach" {
  role       = aws_iam_role.cert_manager_irsa_role.name
  policy_arn = aws_iam_policy.custom_route53_policy.arn
}

resource "aws_iam_role_policy_attachment" "external_dns_route53_attach" {
  role       = aws_iam_role.external_dns_irsa_role.name
  policy_arn = aws_iam_policy.custom_route53_policy.arn
} 