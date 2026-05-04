# resource "aws_ecr_repository" "ecr_repo" {
#   name                 = var.ecr_repo_name
#   image_tag_mutability = "MUTABLE"

#   encryption_configuration {
#     encryption_type = "AES256"
#   }

#   image_scanning_configuration {
#     scan_on_push = true
#   }

#   tags = {
#     Name = var.ecr_repo_name
#   }

# }
