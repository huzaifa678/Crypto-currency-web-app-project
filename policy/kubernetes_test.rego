package kubernetes

import rego.v1

ecr_sha := "533267178572.dkr.ecr.us-east-1.amazonaws.com/crypto-ecr-repo:27d0b18"

ecr_latest := "533267178572.dkr.ecr.us-east-1.amazonaws.com/crypto-ecr-repo:latest"

test_allow_tagged_ecr if {
	obj := {"kind": "Deployment", "metadata": {"name": "api"}, "spec": {"template": {"spec": {"containers": [{"image": ecr_sha}]}}}}
	count(deny) == 0 with input as obj
	count(warn) == 0 with input as obj
}

test_warn_but_allow_ecr_latest if {
	obj := {"kind": "Deployment", "metadata": {"name": "api"}, "spec": {"template": {"spec": {"containers": [{"image": ecr_latest}]}}}}
	count(deny) == 0 with input as obj
	warn with input as obj
}

test_deny_untagged if {
	deny with input as {"kind": "Deployment", "metadata": {"name": "x"}, "spec": {"template": {"spec": {"containers": [{"image": "nginx"}]}}}}
}

test_deny_privileged if {
	deny with input as {"kind": "Deployment", "metadata": {"name": "x"}, "spec": {"template": {"spec": {"containers": [{"image": ecr_sha, "securityContext": {"privileged": true}}]}}}}
}

test_deny_hostpath if {
	deny with input as {"kind": "Deployment", "metadata": {"name": "x"}, "spec": {"template": {"spec": {"volumes": [{"hostPath": {"path": "/"}}], "containers": [{"image": ecr_sha}]}}}}
}

test_allow_test_hook_untagged if {
	obj := {
		"kind": "Pod",
		"metadata": {"name": "test-connection", "annotations": {"helm.sh/hook": "test"}},
		"spec": {"containers": [{"image": "busybox"}]},
	}
	count(deny) == 0 with input as obj
	count(warn) == 0 with input as obj
}

test_warn_unapproved_registry if {
	obj := {"kind": "Deployment", "metadata": {"name": "x"}, "spec": {"template": {"spec": {"containers": [{"image": "docker.io/library/redis:7"}]}}}}
	warn with input as obj
	count(deny) == 0 with input as obj
}
