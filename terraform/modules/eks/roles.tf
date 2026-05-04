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

resource "aws_iam_role_policy" "eks_lb_policy" {
  name = "eks-loadbalancer-policy"
  role = aws_iam_role.eks_cluster_role.name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "ec2:DescribeSecurityGroups",
          "ec2:DescribeSubnets",
          "ec2:DescribeVpcs",
          "ec2:CreateSecurityGroup",
          "ec2:DeleteSecurityGroup",
          "ec2:AuthorizeSecurityGroupIngress",
          "ec2:RevokeSecurityGroupIngress",
          "elasticloadbalancing:*",
          "iam:PassRole"
        ]
        Resource = "*"
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

resource "aws_iam_role_policy_attachment" "eks_node_ssm" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
  role       = aws_iam_role.eks_node_role.name
}

resource "aws_iam_role_policy_attachment" "eks_node_worker_policy" {
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

resource "aws_iam_role_policy_attachment" "eks_node_custom_route53_policy_attachment" {
  policy_arn = aws_iam_policy.custom_route53_policy.arn
  role       = aws_iam_role.eks_node_role.name
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

resource "aws_iam_role_policy_attachment" "cert_manager_irsa_route53_attach" {
  role       = aws_iam_role.cert_manager_irsa_role.name
  policy_arn = aws_iam_policy.custom_route53_policy.arn
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

resource "aws_iam_role_policy_attachment" "external_dns_route53_attach" {
  role       = aws_iam_role.external_dns_irsa_role.name
  policy_arn = aws_iam_policy.custom_route53_policy.arn
}

resource "aws_iam_role" "aws_lb_controller_irsa_role" {
  name = "aws-lb-controller-irsa-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Federated = aws_iam_openid_connect_provider.eks.arn
        }
        Action = "sts:AssumeRoleWithWebIdentity"
        Condition = {
          StringEquals = {
            "${replace(data.aws_eks_cluster.this.identity[0].oidc[0].issuer, "https://", "")}:sub" = "system:serviceaccount:kube-system:aws-load-balancer-controller"
          }
        }
      }
    ]
  })
}

resource "aws_iam_policy" "aws_lb_controller_policy" {
  name        = "aws-lb-controller-policy"
  description = "IAM policy for AWS Load Balancer Controller"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "ec2:DescribeAccountAttributes",
          "ec2:DescribeAddresses",
          "ec2:DescribeInstances",
          "ec2:DescribeSecurityGroups",
          "ec2:DescribeSubnets",
          "ec2:DescribeVpcs",
          "ec2:CreateTags",
          "ec2:DeleteTags",
          "ec2:DescribeTags",
          "ec2:DescribeRouteTables",   
          "ec2:DescribeAvailabilityZones",           
          "ec2:ModifyNetworkInterfaceAttribute",
          "ec2:CreateSecurityGroup",
          "ec2:DeleteSecurityGroup",
          "ec2:AuthorizeSecurityGroupIngress",
          "ec2:RevokeSecurityGroupIngress",
          "elasticloadbalancing:*",
          "iam:PassRole"
        ],
        Resource = "*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "aws_lb_controller_attach" {
  policy_arn = aws_iam_policy.aws_lb_controller_policy.arn
  role       = aws_iam_role.aws_lb_controller_irsa_role.name
}