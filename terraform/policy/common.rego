# Shared helpers for the Conftest / OPA policy set.
#
# Policies run against the JSON representation of a Terraform plan, produced by:
#   terraform show -json <planfile> > tfplan.json
#   conftest test tfplan.json --policy policy --all-namespaces
package main

import rego.v1

# Resources of `kind` that will still exist after the plan is applied
# (i.e. not being destroyed). Replacements ["delete","create"] are included.
managed_resources(kind) := [resource |
	some resource in input.resource_changes
	resource.type == kind
	some action in resource.change.actions
	action != "delete"
]

# Well-known ports that must never be reachable from the public internet.
sensitive_ports := {22, 3389, 5432, 6379, 3306, 27017, 9200, 11211}

# CIDRs that mean "the whole internet".
open_world := {"0.0.0.0/0", "::/0"}

# True when the port range [from, to] covers a sensitive port, or is a
# fully-open range (0-0 with protocol -1, or a 0-65535 span).
covers_sensitive_port(from_port, to_port) if {
	some port in sensitive_ports
	from_port <= port
	to_port >= port
}

covers_sensitive_port(from_port, to_port) if {
	from_port == 0
	to_port == 0
}
