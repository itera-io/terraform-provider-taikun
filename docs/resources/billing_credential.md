---
page_title: "taikun_billing_credential Resource - terraform-provider-taikun"
subcategory: ""
description: |-   Taikun Billing Credential
---

# taikun_billing_credential (Resource)

Taikun Billing Credential

~> **Role Requirement** To use the `taikun_billing_credential` resource, you need a Partner account.

-> **Organization ID** `organization_id` can be specified for the Partner role, it otherwise defaults to the user's organization.

## Example Usage

```terraform
resource "taikun_billing_credential" "foo" {
  name                = "foo"
  prometheus_password = "password"
  prometheus_url      = "url"
  prometheus_username = "username"

  organization_id = "42"
  lock            = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the billing credential.
- `prometheus_password` (String, Sensitive) The Prometheus password.
- `prometheus_url` (String) The Prometheus URL.
- `prometheus_username` (String) The Prometheus username.

### Optional

- `lock` (Boolean) Indicates whether to lock the billing credential. Defaults to `false`.
- `organization_id` (String) The ID of the organization which owns the billing credential.

### Read-Only

- `created_by` (String) The creator of the billing credential.
- `id` (String) The ID of the billing credential.
- `is_default` (Boolean) Indicates whether the billing credential is the organization's default.
- `last_modified` (String) Time and date of last modification.
- `last_modified_by` (String) The last user to have modified the billing credential.
- `organization_name` (String) The name of the organization which owns the billing credential.

## Import

Import is supported using the following syntax:

```shell
terraform import taikun_billing_credential.mycredential 42
```
