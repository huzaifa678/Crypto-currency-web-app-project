package main

import rego.v1

deny contains msg if {
	some rule in managed_resources("aws_security_group_rule")
	rule.change.after.type == "ingress"
	some cidr in rule.change.after.cidr_blocks
	open_world[cidr]
	covers_sensitive_port(rule.change.after.from_port, rule.change.after.to_port)
	msg := sprintf(
		"Security group rule %q opens port range %d-%d to %s",
		[rule.address, rule.change.after.from_port, rule.change.after.to_port, cidr],
	)
}

deny contains msg if {
	some sg in managed_resources("aws_security_group")
	some ingress in sg.change.after.ingress
	some cidr in ingress.cidr_blocks
	open_world[cidr]
	covers_sensitive_port(ingress.from_port, ingress.to_port)
	msg := sprintf(
		"Security group %q has an inline ingress rule opening port range %d-%d to %s",
		[sg.address, ingress.from_port, ingress.to_port, cidr],
	)
}
