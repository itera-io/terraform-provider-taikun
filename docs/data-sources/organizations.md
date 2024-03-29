---
page_title: "taikun_organizations Data Source - terraform-provider-taikun"
subcategory: ""
description: |-   Retrieve all organizations.
---

# taikun_organizations (Data Source)

Retrieve all organizations.

~> **Role Requirement** To use the `taikun_organizations` data source, you need a Partner account.

## Example Usage

```terraform
data "taikun_organizations" "all" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `organizations` (List of Object) List of retrieved organizations. (see [below for nested schema](#nestedatt--organizations))

<a id="nestedatt--organizations"></a>
### Nested Schema for `organizations`

Read-Only:

- `address` (String)
- `billing_email` (String)
- `city` (String)
- `cloud_credentials` (Number)
- `country` (String)
- `created_at` (String)
- `discount_rate` (Number)
- `email` (String)
- `full_name` (String)
- `id` (String)
- `is_read_only` (Boolean)
- `lock` (Boolean)
- `managers_can_change_subscription` (Boolean)
- `name` (String)
- `partner_id` (String)
- `partner_name` (String)
- `phone` (String)
- `projects` (Number)
- `servers` (Number)
- `users` (Number)
- `vat_number` (String)


