---
page_title: "taikun_billing_credentials Data Source - terraform-provider-taikun"
subcategory: ""
description: |-   Retrieve all billing credentials.
---

# taikun_billing_credentials (Data Source)

Retrieve all billing credentials.

~> **Role Requirement** To use the `taikun_billing_credentials` data source, you need a Partner account.

## Example Usage

```terraform
data "taikun_billing_credentials" "foo" {
  organization_id = "42"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `organization_id` (String) Organization ID filter.

### Read-Only

- `billing_credentials` (List of Object) List of retrieved billing credentials. (see [below for nested schema](#nestedatt--billing_credentials))
- `id` (String) The ID of this resource.

<a id="nestedatt--billing_credentials"></a>
### Nested Schema for `billing_credentials`

Read-Only:

- `created_by` (String)
- `id` (String)
- `is_default` (Boolean)
- `last_modified` (String)
- `last_modified_by` (String)
- `lock` (Boolean)
- `name` (String)
- `organization_id` (String)
- `organization_name` (String)
- `prometheus_password` (String)
- `prometheus_url` (String)
- `prometheus_username` (String)


