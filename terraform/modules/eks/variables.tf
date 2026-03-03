variable "cluster_name" {
  type        = string
  description = "EKS cluster name"
}

variable "kubernetes_version" {
  type        = string
  description = "Kubernetes version for the cluster"
}

variable "private_subnets" {
  type        = list(string)
  description = "Private subnets for the cluster"
}

variable "cluster_endpoint_public_access_cidrs" {
  type        = list(string)
  description = "CIDRs allowed to access the EKS public endpoint"
}

variable "vpc_id" {
  type        = string
  description = "VPC ID where EKS will be deployed"
}

variable "environment" {
  type        = string
  description = "Deployment environment (dev/post-test/prod)"
}

variable "region" {
  type        = string
  default     = "us-east-1"
  description = "AWS region"
}