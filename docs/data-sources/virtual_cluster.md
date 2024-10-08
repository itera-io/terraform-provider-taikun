---
page_title: "taikun_virtual_cluster Data Source - terraform-provider-taikun"
subcategory: ""
description: |-   Retrieve a Virtual Project by its ID.
---

# taikun_virtual_cluster (Data Source)

Retrieve a Virtual Project by its ID.

~> **Role Requirement** To use the `taikun_virtual_cluster` data source, you need a Manager or Partner account.

## Example Usage

```terraform
data "taikun_virtual_cluster" "foo" {
  id = "42"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The ID of the Virtual cluster project.

### Read-Only

- `alerting_profile_id` (String) The id of the alerting profile that will be used for the virtual cluster.
- `delete_on_expiration` (Boolean) If enabled, the virtual project will be deleted on the expiration date and it will not be possible to recover it.
- `expiration_date` (String) Virtual project's expiration date in the format: 'dd/mm/yyyy'.
- `hostname` (String) The hostname that will be used for the virtual cluster. If left empty, you are assigned a hostname based on your IP an virtual cluster name.
- `hostname_generated` (String) IP-based resolvable hostname generated by Taikun.
- `name` (String) The name of the virtual cluster.
- `parent_id` (String) The ID of the parent of the virtual cluster.
- `status` (String) Do not set. Used for tracking remote virtual cluster failures.
