# Policy as Code (OPA / Conftest)

Machine-checkable guardrails for the Terraform in this repo, written in
[Rego](https://www.openpolicyagent.org/docs/latest/policy-language/) and run
with [Conftest](https://www.conftest.dev/). They evaluate the **plan**, so they
see the values Terraform will actually apply (including module-expanded and
computed attributes) rather than raw HCL.

## How it runs

**Enforcement (blocks apply)** — in [`.github/workflows/cd.yml`](../../.github/workflows/cd.yml),
between `terraform plan` and `terraform apply`:

```bash
terraform show -json tfplan > tfplan.json
conftest test tfplan.json --policy policy --all-namespaces
```

A `deny` rule fails the job and the apply never happens. A `warn` rule prints
but does not block.

**Unit tests (every PR)** — in
[`.github/workflows/terraform-static-analysis.yml`](../../.github/workflows/terraform-static-analysis.yml),
the `opa-policy` job runs `conftest verify` and `conftest fmt --check`; no AWS
access required.

## Run locally

```bash
cd terraform

# Unit-test the policies
conftest verify --policy policy

# Evaluate a real plan
terraform plan -out tfplan && terraform show -json tfplan > tfplan.json
conftest test tfplan.json --policy policy --all-namespaces
```

## Rules

| File | Level | Rule |
|------|-------|------|
| `rds.rego` | deny | RDS must not be `publicly_accessible` |
| `rds.rego` | deny | RDS must set `storage_encrypted = true` |
| `rds.rego` | warn | RDS should keep automated backups / take a final snapshot |
| `security_groups.rego` | deny | No ingress from `0.0.0.0/0` (or `::/0`) to a sensitive port |
| `ecr.rego` | deny | ECR repositories must enable `scan_on_push` |
| `ecr.rego` | warn | ECR should use `IMMUTABLE` tags |
| `eks.rego` | warn | EKS public API endpoint open to `0.0.0.0/0` |
| `tags.rego` | warn | Long-lived resources should carry a `Name` tag |

`common.rego` holds shared helpers (`managed_resources`, `sensitive_ports`,
`open_world`, `covers_sensitive_port`). Tests live in `*_test.rego`.

## Known current violations

The existing `aws_db_instance.postgres` ([../rds.tf](../rds.tf)) does **not**
set `storage_encrypted`, so the encryption `deny` rule will fail the pipeline
until it is added:

```hcl
resource "aws_db_instance" "postgres" {
  # ...
  storage_encrypted = true
}
```
