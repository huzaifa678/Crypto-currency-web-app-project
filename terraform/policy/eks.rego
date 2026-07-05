package main

import rego.v1

warn contains msg if {
	some cluster in managed_resources("aws_eks_cluster")
	some vpc_config in cluster.change.after.vpc_config
	vpc_config.endpoint_public_access == true
	some cidr in vpc_config.public_access_cidrs
	cidr == "0.0.0.0/0"
	msg := sprintf("EKS cluster %q exposes its public API endpoint to 0.0.0.0/0", [cluster.address])
}
