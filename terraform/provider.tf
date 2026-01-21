terraform {
  required_version = ">= 1.3.0"
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.25.2"
    }

    helm = {
      source  = "hashicorp/helm"
      version = "~> 3.1.1"
    }

    aws = {
      source  = "hashicorp/aws"
      version = ">=5.95.0"
    }

    local = {
      source  = "hashicorp/local"
      version = "~> 2.4.0"
    }

    null = {
      source  = "hashicorp/null"
      version = "~> 3.2.2"
    }

    cloudinit = {
      source  = "hashicorp/cloudinit"
      version = "~> 2.3.4"
    }

  }

  backend "s3" {
    bucket         = "my-terraform-state-bucket-1742982420"
    key            = "crypto-system/terraform.tfstate"
    region         = "us-west-2"
    encrypt        = false
  }
}

provider "aws" {
  region = var.region
}

data "aws_eks_cluster" "eks" {
  name = aws_eks_cluster.eks_cluster.name
}

data "aws_eks_cluster_auth" "eks" {
  name = aws_eks_cluster.eks_cluster.name
}

provider "kubernetes" {
  host                   = data.aws_eks_cluster.eks.endpoint
  cluster_ca_certificate = base64decode(
    data.aws_eks_cluster.eks.certificate_authority[0].data
  )
  token = data.aws_eks_cluster_auth.eks.token
}

provider "helm" {
  kubernetes {
    host                   = data.aws_eks_cluster.eks.endpoint
    cluster_ca_certificate = base64decode(
      data.aws_eks_cluster.eks.certificate_authority[0].data
    )
    token = data.aws_eks_cluster_auth.eks.token
  }
}
