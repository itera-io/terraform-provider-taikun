---
page_title: "taikun_images_aws Data Source - terraform-provider-taikun"
subcategory: ""
description: |-   Retrieve images for a given AWS cloud credential.
---

# taikun_images_aws (Data Source)

Retrieve images for a given AWS cloud credential.

~> **Role Requirement** To use the `taikun_images_aws` data source, you need a Manager or Partner account.

## Example Usage

```terraform
resource "taikun_cloud_credential_aws" "foo" {
  name              = "foo"
  availability_zone = "eu-central-1"
}

data "taikun_images_aws" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id
  latest              = true
  owners              = ["Canonical"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cloud_credential_id` (String) AWS cloud credential ID.

### Optional

- `latest` (Boolean) Retrieve latest AWS images. Defaults to `false`.
- `owners` (Set of String) List of AWS image owners

### Read-Only

- `id` (String) The ID of this resource.
- `images` (List of Object) List of retrieved AWS images. (see [below for nested schema](#nestedatt--images))

<a id="nestedatt--images"></a>
### Nested Schema for `images`

Read-Only:

- `id` (String)
- `name` (String)

