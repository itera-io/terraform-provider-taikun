---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "taikun_catalog Data Source - terraform-provider-taikun"
subcategory: ""
description: |-
  Get an Catalog by its name.
---

# taikun_catalog (Data Source)

Get an Catalog by its name.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the catalog.

### Read-Only

- `application` (Set of Object) Bound Applications. (see [below for nested schema](#nestedatt--application))
- `default` (Boolean) Indicates whether to the catalog is the default catalog.
- `description` (String) The description of the catalog.
- `id` (String) The ID of the catalog.
- `lock` (Boolean) Indicates whether to lock the catalog.
- `projects` (Set of String) List of projects bound to the catalog.

<a id="nestedatt--application"></a>
### Nested Schema for `application`

Read-Only:

- `id` (String)
- `name` (String)
- `repository` (String)