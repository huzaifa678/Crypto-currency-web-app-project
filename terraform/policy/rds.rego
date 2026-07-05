package main

import rego.v1

deny contains msg if {
	some rds in managed_resources("aws_db_instance")
	rds.change.after.publicly_accessible == true
	msg := sprintf("RDS instance %q must not be publicly accessible (set publicly_accessible = false)", [rds.address])
}

deny contains msg if {
	some rds in managed_resources("aws_db_instance")
	not rds.change.after.storage_encrypted
	msg := sprintf("RDS instance %q must have storage_encrypted = true", [rds.address])
}

warn contains msg if {
	some rds in managed_resources("aws_db_instance")
	object.get(rds.change.after, "backup_retention_period", 0) == 0
	msg := sprintf("RDS instance %q has no automated backups (backup_retention_period = 0)", [rds.address])
}

warn contains msg if {
	some rds in managed_resources("aws_db_instance")
	rds.change.after.skip_final_snapshot == true
	msg := sprintf("RDS instance %q skips its final snapshot (skip_final_snapshot = true)", [rds.address])
}
