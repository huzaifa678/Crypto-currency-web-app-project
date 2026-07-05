data "aws_caller_identity" "current" {}


data "aws_iam_policy_document" "kms_default" {
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

resource "aws_kms_key" "rds" {
  description             = "CMK for RDS PostgreSQL encryption"
  enable_key_rotation     = true
  deletion_window_in_days = 7
  policy                  = data.aws_iam_policy_document.kms_default.json

  tags = {
    Name = "crypto-rds"
  }
}

resource "aws_kms_alias" "rds" {
  name          = "alias/crypto-rds"
  target_key_id = aws_kms_key.rds.key_id
}

resource "aws_kms_key" "elasticache" {
  description             = "CMK for ElastiCache Redis encryption"
  enable_key_rotation     = true
  deletion_window_in_days = 7
  policy                  = data.aws_iam_policy_document.kms_default.json

  tags = {
    Name = "crypto-elasticache"
  }
}

resource "aws_kms_alias" "elasticache" {
  name          = "alias/crypto-elasticache"
  target_key_id = aws_kms_key.elasticache.key_id
}

resource "aws_kms_key" "secrets" {
  description             = "CMK for Secrets Manager encryption"
  enable_key_rotation     = true
  deletion_window_in_days = 7
  policy                  = data.aws_iam_policy_document.kms_default.json

  tags = {
    Name = "crypto-secrets"
  }
}

resource "aws_kms_alias" "secrets" {
  name          = "alias/crypto-secrets"
  target_key_id = aws_kms_key.secrets.key_id
}
