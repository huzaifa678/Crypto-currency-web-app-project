package actions

import rego.v1

uses_refs contains u if {
	walk(input, [path, value])
	path[count(path) - 1] == "uses"
	is_string(value)
	u := value
}

pinnable(u) if {
	not startswith(u, "./")
	not startswith(u, "docker://")
}

ref_of(u) := ref if {
	contains(u, "@")
	parts := split(u, "@")
	ref := parts[count(parts) - 1]
}

is_sha(ref) if regex.match(`^[0-9a-f]{40}$`, ref)

mutable_refs := {"main", "master", "latest", "head", "develop"}

first_party(u) if startswith(u, "actions/")

first_party(u) if startswith(u, "github/")

deny contains msg if {
	some u in uses_refs
	pinnable(u)
	not contains(u, "@")
	msg := sprintf("action %q is not pinned to a version/ref", [u])
}

deny contains msg if {
	some u in uses_refs
	pinnable(u)
	lower(ref_of(u)) in mutable_refs
	msg := sprintf("action %q is pinned to a mutable ref %q — pin a tag or SHA", [u, ref_of(u)])
}

deny contains msg if {
	trigger_present("pull_request_target")
	msg := "workflow triggers on pull_request_target (privileged context) — avoid, or never checkout the untrusted PR head"
}

trigger_present(name) if input["true"][name]

trigger_present(name) if input["true"][_] == name

trigger_present(name) if input["true"] == name

warn contains msg if {
	some u in uses_refs
	pinnable(u)
	not first_party(u)
	contains(u, "@")
	not is_sha(ref_of(u))
	msg := sprintf("third-party action %q is tag-pinned, not SHA-pinned", [u])
}

warn contains msg if {
	not input.permissions
	not any_job_permissions
	msg := "workflow sets no permissions: — GITHUB_TOKEN defaults to a broad scope; set least-privilege permissions"
}

any_job_permissions if {
	some job in input.jobs
	job.permissions
}
