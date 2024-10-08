---
page_title: "taikun_virtual_cluster Resource - terraform-provider-taikun"
subcategory: ""
description: |-   Virtual Cluster project in Taikun.
---

# taikun_virtual_cluster (Resource)

Virtual Cluster project in Taikun.

~> **Role Requirement** To use the `taikun_virtual_cluster` resource, you need a Manager or Partner account.

-> **Organization ID** `organization_id` cannot be specified. Virtual project is created in the same organization as the parent project.

## Example Usage

```terraform
resource "taikun_virtual_cluster" "foo" {
  name                 = "test-virtual-cluster-42"
  parent_id            = 424242
  expiration_date      = "20/01/2050"
  delete_on_expiration = "true"
}
```

Take a look at the **quickstart template** that deploys a new k8s cluster with virtual cluster inside - available in the [quickstart examples](https://github.com/itera-io/terraform-provider-taikun/tree/dev/examples/quickstart-templates) folder.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the virtual cluster.
- `parent_id` (String) The ID of the parent of the virtual cluster.

### Optional

- `alerting_profile_id` (String) The id of the alerting profile that will be used for the virtual cluster.
- `delete_on_expiration` (Boolean) If enabled, the virtual project will be deleted on the expiration date and it will not be possible to recover it. Defaults to `false`. Required with: `expiration_date`.
- `expiration_date` (String) Virtual project's expiration date in the format: 'dd/mm/yyyy'.
- `hostname` (String) The hostname that will be used for the virtual cluster. If left empty, you are assigned a hostname based on your IP an virtual cluster name. Defaults to ` `.
- `status` (String) Do not set. Used for tracking remote virtual cluster failures. Defaults to ` `.

### Read-Only

- `hostname_generated` (String) IP-based resolvable hostname generated by Taikun.
- `id` (String) The ID of the Virtual cluster project.

## Import

Import is supported using the following syntax:

```shell
terraform import taikun_virtual_cluster.myvirtualcluster 42
```
