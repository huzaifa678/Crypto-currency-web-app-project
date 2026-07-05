variable "github_owner" {
  description = "GitHub organisation or user that owns the repository"
  type        = string
  default     = "huzaifa678"
}

variable "github_repo" {
  description = "GitHub repository name allowed to assume the CI role"
  type        = string
  default     = "Crypto-currency-web-app-project"
}

variable "github_oidc_subjects" {
  description = "Allowed `sub` claims (branches / environments / tags) that may assume the role"
  type        = list(string)
  default = [
    "ref:refs/heads/main",
    "environment:main",
  ]
}

resource "aws_iam_openid_connect_provider" "github" {
  url             = "https://token.actions.githubusercontent.com"
  client_id_list  = ["sts.amazonaws.com"]
  thumbprint_list = ["6938fd4d98bab03faadb97b34396831e3780aea1"]
}

data "aws_iam_policy_document" "github_actions_assume" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRoleWithWebIdentity"]

    principals {
      type        = "Federated"
      identifiers = [aws_iam_openid_connect_provider.github.arn]
    }

    condition {
      test     = "StringEquals"
      variable = "token.actions.githubusercontent.com:aud"
      values   = ["sts.amazonaws.com"]
    }

    condition {
      test     = "StringLike"
      variable = "token.actions.githubusercontent.com:sub"
      values   = [for sub in var.github_oidc_subjects : "repo:${var.github_owner}/${var.github_repo}:${sub}"]
    }
  }
}

resource "aws_iam_role" "github_actions" {
  name                 = "github-actions-crypto-system"
  description          = "Role assumed by GitHub Actions via OIDC for Terraform + ECR deployments"
  assume_role_policy   = data.aws_iam_policy_document.github_actions_assume.json
  max_session_duration = 3600
}

variable "tf_state_bucket" {
  description = "S3 bucket holding the Terraform remote state (see backend config in provider.tf)"
  type        = string
  default     = "my-terraform-state-bucket-1742982420"
}

variable "tf_state_key_prefix" {
  description = "Key prefix within the state bucket that CI is allowed to read/write"
  type        = string
  default     = "crypto-system"
}


data "aws_iam_policy_document" "github_actions_ci" {
  # --- Terraform remote state (tightly scoped to the one bucket + prefix) ---
  statement {
    sid       = "TerraformStateBucketList"
    effect    = "Allow"
    actions   = ["s3:ListBucket", "s3:GetBucketLocation"]
    resources = ["arn:aws:s3:::${var.tf_state_bucket}"]
  }

  statement {
    sid    = "TerraformStateObjects"
    effect = "Allow"
    actions   = ["s3:GetObject", "s3:PutObject", "s3:DeleteObject"]
    resources = ["arn:aws:s3:::${var.tf_state_bucket}/${var.tf_state_key_prefix}/*"]
  }

\  statement {
    sid    = "ComputeNetworking"
    effect = "Allow"
    actions = [
      "ec2:*",
      "autoscaling:*",
    ]
    resources = ["*"]
  }

  statement {
    sid       = "LoadBalancerRead"
    effect    = "Allow"
    actions   = ["elasticloadbalancing:Describe*"]
    resources = ["*"]
  }

  statement {
    sid    = "ManagedServices"
    effect = "Allow"
    actions = [
      "eks:*",
      "ecr:*",
      "rds:*",
      "elasticache:*",
      "secretsmanager:*",
      "ssm:*",
      "kms:CreateKey",
      "kms:CreateAlias",
      "kms:DeleteAlias",
      "kms:CreateGrant",
      "kms:TagResource",
      "kms:UntagResource",
      "kms:ScheduleKeyDeletion",
      "kms:EnableKeyRotation",
      "kms:PutKeyPolicy",
      "kms:Describe*",
      "kms:Get*",
      "kms:List*",
      "kms:Encrypt",
      "kms:Decrypt",
      "kms:GenerateDataKey*",
      "logs:CreateLogGroup",
      "logs:DeleteLogGroup",
      "logs:DescribeLogGroups",
      "logs:PutRetentionPolicy",
      "logs:TagResource",
      "logs:ListTagsForResource",
    ]
    resources = ["*"]
  }

  statement {
    sid    = "IamManagement"
    effect = "Allow"
    actions = [
      "iam:CreateRole",
      "iam:DeleteRole",
      "iam:UpdateRole",
      "iam:GetRole",
      "iam:ListRoles",
      "iam:PassRole",
      "iam:TagRole",
      "iam:UntagRole",
      "iam:CreatePolicy",
      "iam:DeletePolicy",
      "iam:GetPolicy",
      "iam:ListPolicies",
      "iam:CreatePolicyVersion",
      "iam:DeletePolicyVersion",
      "iam:GetPolicyVersion",
      "iam:ListPolicyVersions",
      "iam:AttachRolePolicy",
      "iam:DetachRolePolicy",
      "iam:ListAttachedRolePolicies",
      "iam:PutRolePolicy",
      "iam:DeleteRolePolicy",
      "iam:GetRolePolicy",
      "iam:ListRolePolicies",
      "iam:CreateOpenIDConnectProvider",
      "iam:DeleteOpenIDConnectProvider",
      "iam:GetOpenIDConnectProvider",
      "iam:ListOpenIDConnectProviders",
      "iam:TagOpenIDConnectProvider",
      "iam:UpdateOpenIDConnectProviderThumbprint",
      "iam:CreateInstanceProfile",
      "iam:DeleteInstanceProfile",
      "iam:GetInstanceProfile",
      "iam:AddRoleToInstanceProfile",
      "iam:RemoveRoleFromInstanceProfile",
      "iam:CreateServiceLinkedRole",
      "iam:ListInstanceProfilesForRole",
      "iam:CreateUser",
      "iam:DeleteUser",
      "iam:GetUser",
      "iam:TagUser",
      "iam:AttachUserPolicy",
      "iam:DetachUserPolicy",
      "iam:ListAttachedUserPolicies",
    ]
    resources = ["*"]
  }

  statement {
    sid       = "CallerIdentity"
    effect    = "Allow"
    actions   = ["sts:GetCallerIdentity"]
    resources = ["*"]
  }

  statement {
    sid    = "DenySelfMutation"
    effect = "Deny"
    actions = [
      "iam:*Role*",
      "iam:*RolePolicy*",
      "iam:PassRole",
    ]

    resources = [
      aws_iam_role.github_actions.arn,
      "arn:aws:iam::${data.aws_caller_identity.current.account_id}:policy/github-actions-crypto-system-ci",
    ]
  }
}

resource "aws_iam_policy" "github_actions_ci" {
  name        = "github-actions-crypto-system-ci"
  description = "Least-privilege permissions for the GitHub Actions OIDC CI role"
  policy      = data.aws_iam_policy_document.github_actions_ci.json
}

resource "aws_iam_role_policy_attachment" "github_actions_ci" {
  role       = aws_iam_role.github_actions.name
  policy_arn = aws_iam_policy.github_actions_ci.arn
}

output "github_actions_role_arn" {
  description = "Set this as the AWS_ROLE_ARN GitHub Actions secret"
  value       = aws_iam_role.github_actions.arn
}
