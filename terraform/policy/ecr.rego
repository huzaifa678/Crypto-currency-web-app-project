package main

import rego.v1

deny contains msg if {
	some repo in managed_resources("aws_ecr_repository")
	not repo.change.after.image_scanning_configuration[0].scan_on_push
	msg := sprintf("ECR repository %q must enable scan_on_push", [repo.address])
}

warn contains msg if {
	some repo in managed_resources("aws_ecr_repository")
	repo.change.after.image_tag_mutability != "IMMUTABLE"
	msg := sprintf("ECR repository %q should set image_tag_mutability = \"IMMUTABLE\"", [repo.address])
}
