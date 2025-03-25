variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "profile" {
  description = "AWS CLI profile"
  type        = string
  default     = "default"
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnets" {
  description = "Public subnet CIDR blocks"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable "private_subnets" {
  description = "Private subnet CIDR blocks"
  type        = list(string)
  default     = ["10.0.3.0/24", "10.0.4.0/24"]
}

variable "cluster_name" {
  description = "EKS cluster name"
  type        = string
  default     = "crypto-system-eks-cluster"
}

variable "kubernetes_version" {
  default     = 1.28
  description = "kubernetes version"
}

variable "desired_capacity" {
  description = "Desired capacity of worker nodes"
  type        = number
  default     = 2
}

variable "max_size" {
  description = "Max size of worker nodes"
  type        = number
  default     = 3
}

variable "min_size" {
  description = "Min size of worker nodes"
  type        = number
  default     = 1
}

variable "ecr_repo_name" {
  description = "ECR repository name"
  type        = string
  default     = "crypto-ecr-repo"
}

variable "rds_db_name" {
    description = "RDS DB name"
    type        = string
    default     = "cryptodb"
}

variable "rds_db_username" {
    description = "username credential for RDS DB"
    type        = string
    default     = "root"
}

variable "rds_db_password" {
    description = "password credential for RDS DB"
    type        = string
    default     = "secret1234"
}
