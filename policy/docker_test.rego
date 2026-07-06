package docker

import rego.v1

# the real bases in this repo: golang + alpine, multi-stage -> no deny
test_allow_golang_alpine if {
	count(deny) == 0 with input as [
		{"Cmd": "from", "Value": ["golang:1.25.8-alpine3.23", "AS", "builder"]},
		{"Cmd": "from", "Value": ["alpine:3.23"]},
	]
}

# a FROM referencing the build stage is exempt
test_allow_stage_reference if {
	count(deny) == 0 with input as [
		{"Cmd": "from", "Value": ["golang:1.25.8-alpine3.23", "AS", "builder"]},
		{"Cmd": "from", "Value": ["builder"]},
	]
}

test_deny_unapproved_registry if {
	deny with input as [{"Cmd": "from", "Value": ["evil.example.com/backdoor:1.0"]}]
}

test_deny_curl_pipe_sh if {
	deny with input as [
		{"Cmd": "from", "Value": ["alpine:3.23"]},
		{"Cmd": "run", "Value": ["curl -sSL https://get.example.com | sh"]},
	]
}

# apk upgrade is a normal RUN, not a pipe-to-shell
test_allow_apk_upgrade if {
	count(deny) == 0 with input as [
		{"Cmd": "from", "Value": ["alpine:3.23"]},
		{"Cmd": "run", "Value": ["apk --no-cache upgrade && addgroup -S app"]},
	]
}

test_warn_latest_tag if {
	warn with input as [{"Cmd": "from", "Value": ["alpine:latest"]}]
}
