package main

import rego.v1

bad_plan := {"resource_changes": [{
	"address": "aws_db_instance.postgres",
	"type": "aws_db_instance",
	"change": {"actions": ["create"], "after": {
		"publicly_accessible": true,
		"storage_encrypted": false,
		"skip_final_snapshot": true,
		"backup_retention_period": 0,
	}},
}]}

good_plan := {"resource_changes": [{
	"address": "aws_db_instance.postgres",
	"type": "aws_db_instance",
	"change": {"actions": ["create"], "after": {
		"publicly_accessible": false,
		"storage_encrypted": true,
		"skip_final_snapshot": false,
		"backup_retention_period": 7,
		"tags": {"Name": "cryptodb"},
	}},
}]}

test_public_rds_denied if {
	some msg in deny with input as bad_plan
	contains(msg, "must not be publicly accessible")
}

test_unencrypted_rds_denied if {
	some msg in deny with input as bad_plan
	contains(msg, "storage_encrypted = true")
}

test_hardened_rds_allowed if {
	count(deny) == 0 with input as good_plan
}

test_deleted_rds_ignored if {
	count(deny) == 0 with input as {"resource_changes": [{
		"address": "aws_db_instance.postgres",
		"type": "aws_db_instance",
		"change": {"actions": ["delete"], "after": null},
	}]}
}
