---
page_title: "taikun_policy_profile Resource - terraform-provider-taikun"
subcategory: ""
description: |-   Taikun Policy Profile
---

# taikun_policy_profile (Resource)

Taikun Policy Profile

~> **Role Requirement** In order to use the `taikun_policy_profile` resource you need a `Manager` or `Partner` account.

-> **Organization ID** `organization_id` can be specified for the Partner role, it otherwise defaults to the user's organization.

## Example Usage

```terraform
resource "taikun_policy_profile" "foo" {
  name = "foo"

  forbid_node_port        = true
  forbid_http_ingress     = true
  require_probe           = true
  unique_ingress          = true
  unique_service_selector = true

  allowed_repos = [
    "repo"
  ]
  forbidden_tags = [
    "tag"
  ]
  ingress_whitelist = [
    "ingress"
  ]

  organization_id = "42"
  lock            = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the Policy profile.

### Optional

- `allowed_repos` (Set of String) Requires container images to begin with a string from the specified list.
- `forbid_http_ingress` (Boolean) Requires Ingress resources to be HTTPS only. Defaults to `false`.
- `forbid_node_port` (Boolean) Disallows all Services with type NodePort. Defaults to `false`.
- `forbidden_tags` (Set of String) Container images must have an image tag different from the ones in the list.
- `ingress_whitelist` (Set of String) List of allowed Ingress rule hosts.
- `lock` (Boolean) Indicates whether to lock the Policy profile. Defaults to `false`.
- `organization_id` (String) The ID of the organization which owns the Policy profile.
- `require_probe` (Boolean) Requires Pods to have readiness and liveness probes. Defaults to `false`.
- `unique_ingress` (Boolean) Requires all Ingress rule hosts to be unique. Defaults to `false`.
- `unique_service_selector` (Boolean) Whether services must have globally unique service selectors or not. Defaults to `false`.

### Read-Only

- `id` (String) The ID of the Policy profile.
- `is_default` (Boolean) Indicates whether the Policy Profile is the default one.
- `organization_name` (String) The name of the organization which owns the Policy profile.

## Import

Import is supported using the following syntax:

```shell
terraform import taikun_policy_profile.myprofile 42
```
