package main

import rego.v1

taggable := {
	"aws_db_instance",
	"aws_db_subnet_group",
	"aws_eks_cluster",
	"aws_ecr_repository",
	"aws_elasticache_replication_group",
	"aws_elasticache_subnet_group",
}

warn contains msg if {
	some resource in input.resource_changes
	taggable[resource.type]
	some action in resource.change.actions
	action != "delete"
	tags := object.get(resource.change.after, "tags", {})
	not tags.Name
	msg := sprintf("Resource %q is missing a Name tag", [resource.address])
}
