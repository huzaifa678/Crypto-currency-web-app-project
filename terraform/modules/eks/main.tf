data "aws_caller_identity" "current" {}

data "aws_iam_policy_document" "eks_kms" {
  # checkov:skip=CKV_AWS_109: Standard KMS key policy delegating access to account IAM.
  # checkov:skip=CKV_AWS_111: kms:* to the account root is the AWS-default key policy.
  # checkov:skip=CKV_AWS_356: A key policy's "*" resource refers to the key itself; not scopable.
  statement {
    sid    = "EnableRootAccountAccess"
    effect = "Allow"
    principals {
      type        = "AWS"
      identifiers = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"]
    }
    actions   = ["kms:*"]
    resources = ["*"]
  }
}

resource "aws_kms_key" "eks" {
  description             = "CMK for EKS secrets envelope encryption"
  enable_key_rotation     = true
  deletion_window_in_days = 7
  policy                  = data.aws_iam_policy_document.eks_kms.json

  tags = {
    Name = "${var.cluster_name}-secrets"
  }
}

resource "aws_kms_alias" "eks" {
  name          = "alias/${var.cluster_name}-secrets"
  target_key_id = aws_kms_key.eks.key_id
}

resource "aws_eks_cluster" "eks_cluster" {
  # checkov:skip=CKV_AWS_38: Public endpoint is required for GitHub-hosted CI
  name     = var.cluster_name
  role_arn = aws_iam_role.eks_cluster_role.arn
  version  = var.kubernetes_version

  enabled_cluster_log_types = ["api", "audit", "authenticator", "controllerManager", "scheduler"]

  vpc_config {
    subnet_ids              = var.private_subnets
    endpoint_private_access = true
    endpoint_public_access  = true
    public_access_cidrs     = var.cluster_endpoint_public_access_cidrs
  }

  # Envelope-encrypt Kubernetes secrets with the CMK above.
  encryption_config {
    provider {
      key_arn = aws_kms_key.eks.arn
    }
    resources = ["secrets"]
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

data "aws_ssm_parameter" "eks_al2023_ami" {
  name = "/aws/service/eks/optimized-ami/1.32/amazon-linux-2023/x86_64/standard/recommended/image_id"
}

resource "aws_launch_template" "eks_nodes" {
  name_prefix   = "eks-nodes-lt"
  instance_type = "t3.small"

  metadata_options {
    http_endpoint               = "enabled"
    http_tokens                 = "required"
    http_put_response_hop_limit = 1
  }

  vpc_security_group_ids = [
    aws_security_group.eks_nodes.id,
    aws_eks_cluster.eks_cluster.vpc_config[0].cluster_security_group_id
  ]

  image_id = data.aws_ssm_parameter.eks_al2023_ami.value

  user_data = base64encode(templatefile("${path.module}/user_data.tpl", {
    cluster_name     = aws_eks_cluster.eks_cluster.name
    cluster_endpoint = data.aws_eks_cluster.this.endpoint
    cluster_ca       = data.aws_eks_cluster.this.certificate_authority[0].data
    cidr             = aws_eks_cluster.eks_cluster.kubernetes_network_config[0].service_ipv4_cidr
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
  subnet_ids      = var.private_subnets

  scaling_config {
    desired_size = 3
    min_size     = 1
    max_size     = 4
  }

  ami_type = "CUSTOM"

  launch_template {
    id      = aws_launch_template.eks_nodes.id
    version = aws_launch_template.eks_nodes.latest_version
  }

  depends_on = [aws_security_group_rule.allow_node_to_control_plane]
}

resource "aws_iam_openid_connect_provider" "eks" {
  url             = data.aws_eks_cluster.this.identity[0].oidc[0].issuer
  client_id_list  = ["sts.amazonaws.com"]
  thumbprint_list = ["9e99a48a9960b14926bb7f3b02e22da0ecd4e0a4"]
}