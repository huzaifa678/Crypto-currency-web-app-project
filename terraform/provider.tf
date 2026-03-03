terraform {
  required_version = ">= 1.3.0"
  required_providers {
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
  }

  backend "s3" {
    bucket         = "my-terraform-state-bucket-1742982420"
    key            = "crypto-system/terraform.tfstate"
    region         = "us-west-2"
    use_lockfile = true
    encrypt        = true
  }
}

provider "aws" {
  region = var.region
}

# data "aws_eks_cluster" "eks" {
#   name = aws_eks_cluster.eks_cluster.name
# }

# provider "kubectl" {
#   host                   = data.aws_eks_cluster.eks.endpoint
#   cluster_ca_certificate = base64decode(
#     data.aws_eks_cluster.eks.certificate_authority[0].data
#   )

#   exec {
#     api_version = "client.authentication.k8s.io/v1beta1"
#     command     = "aws"
#     args = [
#       "eks",
#       "get-token",
#       "--cluster-name",
#       var.cluster_name
#     ]
#   }
# }


# provider "kubernetes" {
#   alias                  = "eks"
#   host                   = data.aws_eks_cluster.eks.endpoint
#   cluster_ca_certificate = base64decode(
#     data.aws_eks_cluster.eks.certificate_authority[0].data
#   )

#   exec {
#     api_version = "client.authentication.k8s.io/v1beta1"
#     command     = "aws"
#     args = [
#       "eks",
#       "get-token",
#       "--cluster-name",
#       var.cluster_name
#     ]
#   }
# }


# provider "helm" {
#   kubernetes = {
#     host                   = data.aws_eks_cluster.eks.endpoint
#     cluster_ca_certificate = base64decode(data.aws_eks_cluster.eks.certificate_authority[0].data)
#     exec = {
#       api_version = "client.authentication.k8s.io/v1beta1"
#       args        = ["eks", "get-token", "--cluster-name", var.cluster_name]
#       command     = "aws"
#     }
#   }
# }
