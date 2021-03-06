---
page_title: "taikun_showback_credential Resource - terraform-provider-taikun"
subcategory: ""
description: |-   Taikun Showback Credential
---

# taikun_showback_credential (Resource)

Taikun Showback Credential

~> **Role Requirement** To use the `taikun_showback_credential` resource, you need a Manager or Partner account.

-> **Organization ID** `organization_id` can be specified for the Partner role, it otherwise defaults to the user's organization.

## Example Usage

```terraform
resource "taikun_showback_credential" "foo" {
  name     = "foo"
  password = "password"
  url      = "url"
  username = "username"

  organization_id = "42"
  lock            = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the showback credential.
- `password` (String, Sensitive) The Prometheus password or other credential.
- `url` (String) URL of the source.
- `username` (String) The Prometheus username or other credential.

### Optional

- `lock` (Boolean) Indicates whether to lock the showback credential. Defaults to `false`.
- `organization_id` (String) The ID of the organization which owns the showback credential.

### Read-Only

- `created_by` (String) The creator of the showback credential.
- `id` (String) The ID of the showback credential.
- `last_modified` (String) Time of last modification.
- `last_modified_by` (String) The last user who modified the showback credential.
- `organization_name` (String) The name of the organization which owns the showback credential.

## Import

Import is supported using the following syntax:

```shell
terraform import taikun_showback_credential.mycredential 42
```
