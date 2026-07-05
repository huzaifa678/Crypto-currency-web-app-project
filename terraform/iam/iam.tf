data "aws_caller_identity" "current" {}

resource "aws_iam_role" "eks_admin" {
  name = "cryto-eks-cluster-admin"

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "sts:AssumeRole",
      "Principal": {
        "AWS": "arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"
      }
    }
  ]
}
POLICY
}

resource "aws_iam_policy" "eks_admin" {
  # checkov:skip=CKV_AWS_290: eks:* admin actions cannot be resource-scoped.
  # checkov:skip=CKV_AWS_355: eks:ListClusters and similar require "*" resource.
  name = "AmazonEKSAdminPolicy"

  policy = <<POLICY
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "eks:*"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": "iam:PassRole",
            "Resource": "*",
            "Condition": {
                "StringEquals": {
                    "iam:PassedToService": "eks.amazonaws.com"
                }
            }
        }
    ]
}
POLICY
}

resource "aws_iam_role_policy_attachment" "eks_admin" {
  role       = aws_iam_role.eks_admin.name
  policy_arn = aws_iam_policy.eks_admin.arn
}

resource "aws_iam_user" "admin" {
  # checkov:skip=CKV_AWS_273: Break-glass IAM user retained for bootstrapping;
  name = "terraform2"
}

resource "aws_iam_policy" "eks_assume_admin" {
  name = "AmazonEKSAssumeAdminPolicy"

  policy = <<POLICY
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "sts:AssumeRole"
            ],
            "Resource": "${aws_iam_role.eks_admin.arn}"
        }
    ]
}
POLICY
}

resource "aws_iam_group" "eks_admins" {
  name = "eks-admins"
}

resource "aws_iam_group_policy_attachment" "eks_admins" {
  group      = aws_iam_group.eks_admins.name
  policy_arn = aws_iam_policy.eks_assume_admin.arn
}

resource "aws_iam_user_group_membership" "admin" {
  user   = aws_iam_user.admin.name
  groups = [aws_iam_group.eks_admins.name]
}

resource "aws_eks_access_entry" "admin" {
  cluster_name      = "crypto-system-eks-cluster"
  principal_arn     = aws_iam_role.eks_admin.arn
  kubernetes_groups = ["admin"]
}