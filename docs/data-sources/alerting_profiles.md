---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "taikun_alerting_profiles Data Source - terraform-provider-taikun"
subcategory: ""
description: |-
  Get the list of alerting profiles for your organizations, or filter by organization if Partner or Admin
---

# taikun_alerting_profiles (Data Source)

Get the list of alerting profiles for your organizations, or filter by organization if Partner or Admin

## Example Usage

```terraform
data "taikun_alerting_profiles" "foo" {
  # Optional
  organization_id = "42"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- **id** (String) The ID of this resource.
- **organization_id** (String) organization ID filter (for Partner and Admin roles)

### Read-Only

- **alerting_profiles** (List of Object) (see [below for nested schema](#nestedatt--alerting_profiles))

<a id="nestedatt--alerting_profiles"></a>
### Nested Schema for `alerting_profiles`

Read-Only:

- **created_by** (String)
- **emails** (List of String)
- **id** (String)
- **is_locked** (Boolean)
- **last_modified** (String)
- **last_modified_by** (String)
- **name** (String)
- **organization_id** (String)
- **organization_name** (String)
- **reminder** (String)
- **slack_configuration_id** (String)
- **slack_configuration_name** (String)
- **webhook** (Set of Object) (see [below for nested schema](#nestedobjatt--alerting_profiles--webhook))

<a id="nestedobjatt--alerting_profiles--webhook"></a>
### Nested Schema for `alerting_profiles.webhook`

Read-Only:

- **header** (Set of Object) (see [below for nested schema](#nestedobjatt--alerting_profiles--webhook--header))
- **url** (String)

<a id="nestedobjatt--alerting_profiles--webhook--header"></a>
### Nested Schema for `alerting_profiles.webhook.header`

Read-Only:

- **key** (String)
- **value** (String)

