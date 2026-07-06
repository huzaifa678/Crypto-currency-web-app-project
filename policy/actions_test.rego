package actions

import rego.v1

test_clean_workflow if {
	wf := {
		"true": {"push": {"branches": ["main"]}},
		"permissions": {"contents": "read"},
		"jobs": {"a": {"runs-on": "ubuntu-latest", "steps": [{"uses": "actions/checkout@v7"}]}},
	}
	count(deny) == 0 with input as wf
	count(warn) == 0 with input as wf
}

test_deny_unpinned if {
	deny with input as {"jobs": {"a": {"steps": [{"uses": "actions/checkout"}]}}}
}

test_deny_mutable_main if {
	deny with input as {"jobs": {"a": {"steps": [{"uses": "some/action@main"}]}}}
}

test_deny_prt_map if {
	deny with input as {"true": {"pull_request_target": {}}, "jobs": {}}
}

test_deny_prt_list if {
	deny with input as {"true": ["push", "pull_request_target"], "jobs": {}}
}

test_allow_local_action if {
	count(deny) == 0 with input as {"permissions": {"contents": "read"}, "jobs": {"a": {"steps": [{"uses": "./.github/actions/x"}]}}}
}

test_warn_thirdparty_not_sha if {
	wf := {"permissions": {"contents": "read"}, "jobs": {"a": {"steps": [{"uses": "aws-actions/configure-aws-credentials@v6"}]}}}
	warn with input as wf
	count(deny) == 0 with input as wf
}

test_no_warn_sha_pinned if {
	sha := "aws-actions/configure-aws-credentials@0123456789abcdef0123456789abcdef01234567"
	wf := {"permissions": {"contents": "read"}, "jobs": {"a": {"steps": [{"uses": sha}]}}}
	count(warn) == 0 with input as wf
}

test_warn_missing_permissions if {
	warn with input as {"jobs": {"a": {"runs-on": "ubuntu-latest", "steps": [{"uses": "actions/checkout@v7"}]}}}
}

test_no_warn_job_permissions if {
	wf := {"jobs": {"a": {"permissions": {"contents": "read"}, "steps": [{"uses": "actions/checkout@v7"}]}}}
	count([m | some m in warn; m == "workflow sets no permissions: — GITHUB_TOKEN defaults to a broad scope; set least-privilege permissions"]) == 0 with input as wf
}
