---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "taikun_showback_rules Data Source - terraform-provider-taikun"
subcategory: ""
description: |-
  Get the list of showback rules, optionally filtered by organization.
---

# taikun_showback_rules (Data Source)

Get the list of showback rules, optionally filtered by organization.

## Example Usage

```terraform
data "taikun_showback_rules" "foo" {
  # Optional for Partner and Admin
  organization_id = "42"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- **id** (String) The ID of this resource.
- **organization_id** (String) Organization id filter.

### Read-Only

- **showback_rules** (List of Object) (see [below for nested schema](#nestedatt--showback_rules))

<a id="nestedatt--showback_rules"></a>
### Nested Schema for `showback_rules`

Read-Only:

- **created_by** (String)
- **global_alert_limit** (Number)
- **id** (String)
- **kind** (String)
- **label** (List of Object) (see [below for nested schema](#nestedobjatt--showback_rules--label))
- **last_modified** (String)
- **last_modified_by** (String)
- **metric_name** (String)
- **name** (String)
- **organization_id** (String)
- **organization_name** (String)
- **price** (Number)
- **project_alert_limit** (Number)
- **showback_credential_id** (String)
- **showback_credential_name** (String)
- **type** (String)

<a id="nestedobjatt--showback_rules--label"></a>
### Nested Schema for `showback_rules.label`

Read-Only:

- **key** (String)
- **value** (String)

