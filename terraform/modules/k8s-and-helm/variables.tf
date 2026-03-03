variable "cluster_name" {
  description = "EKS cluster name"
  type        = string
}

variable "cluster_endpoint" {
  description = "EKS cluster endpoint"
  type        = string
}

variable "cluster_ca" {
  description = "EKS cluster CA certificate"
  type        = string
}

variable "eks_node_group" {
  description = "EKS node group name"
  type        = any
}

variable "vpc" {
  description = "VPC object"
  type        = any
}

variable "cert_manager_irsa_role_arn" {
  description = "IAM role ARN for cert-manager service account"
  type        = string
}

variable "external_dns_irsa_role_arn" {
  description = "IAM role ARN for external-dns service account"
  type        = string
}

variable "aws_lb_controller_irsa_role_arn" {
  description = "IAM role ARN for AWS Load Balancer Controller"
  type        = string
}

variable "vpc_id" {
  description = "VPC ID, required for some Helm charts like aws-load-balancer-controller"
  type        = string
}

variable "region" {
  description = "AWS region for Helm charts"
  type        = string
}

variable "environment" {
  description = "Environment (test / post-test)"
  type        = string
}