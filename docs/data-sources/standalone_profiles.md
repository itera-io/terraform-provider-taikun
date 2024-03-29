---
page_title: "taikun_standalone_profiles Data Source - terraform-provider-taikun"
subcategory: ""
description: |-   Retrieve all standalone profiles.
---

# taikun_standalone_profiles (Data Source)

Retrieve all standalone profiles.

~> **Role Requirement** To use the `taikun_standalone_profiles` data source, you need a Manager or Partner account.

-> **Organization ID** `organization_id` can be specified for the Partner role, it otherwise defaults to the user's organization.

## Example Usage

```terraform
data "taikun_standalone_profiles" "foo" {
  organization_id = "42"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `organization_id` (String) Organization ID filter.

### Read-Only

- `id` (String) The ID of this resource.
- `standalone_profiles` (List of Object) List of retrieved standalone profiles. (see [below for nested schema](#nestedatt--standalone_profiles))

<a id="nestedatt--standalone_profiles"></a>
### Nested Schema for `standalone_profiles`

Read-Only:

- `id` (String)
- `lock` (Boolean)
- `name` (String)
- `organization_id` (String)
- `organization_name` (String)
- `public_key` (String)
- `security_group` (List of Object) (see [below for nested schema](#nestedobjatt--standalone_profiles--security_group))

<a id="nestedobjatt--standalone_profiles--security_group"></a>
### Nested Schema for `standalone_profiles.security_group`

Read-Only:

- `cidr` (String)
- `from_port` (Number)
- `id` (String)
- `ip_protocol` (String)
- `name` (String)
- `to_port` (Number)


