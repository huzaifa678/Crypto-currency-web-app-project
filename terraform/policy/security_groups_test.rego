package main

import rego.v1

test_ssh_open_to_world_denied if {
	some msg in deny with input as {"resource_changes": [{
		"address": "aws_security_group_rule.ssh",
		"type": "aws_security_group_rule",
		"change": {"actions": ["create"], "after": {
			"type": "ingress",
			"from_port": 22,
			"to_port": 22,
			"cidr_blocks": ["0.0.0.0/0"],
		}},
	}]}
	contains(msg, "opens port range 22-22 to 0.0.0.0/0")
}

test_scoped_ssh_allowed if {
	count(deny) == 0 with input as {"resource_changes": [{
		"address": "aws_security_group_rule.ssh",
		"type": "aws_security_group_rule",
		"change": {"actions": ["create"], "after": {
			"type": "ingress",
			"from_port": 22,
			"to_port": 22,
			"cidr_blocks": ["10.0.0.0/16"],
		}},
	}]}
}

# Egress to 0.0.0.0/0 (or ingress from an SG, cidr_blocks null) must not trip.
test_sg_referenced_ingress_allowed if {
	count(deny) == 0 with input as {"resource_changes": [{
		"address": "aws_security_group_rule.db",
		"type": "aws_security_group_rule",
		"change": {"actions": ["create"], "after": {
			"type": "ingress",
			"from_port": 5432,
			"to_port": 5432,
			"cidr_blocks": null,
			"source_security_group_id": "sg-123",
		}},
	}]}
}
