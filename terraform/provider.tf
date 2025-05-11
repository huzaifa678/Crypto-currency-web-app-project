terraform {
  required_version = ">= 1.3.0"
  required_providers {
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.25.2"
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