---
page_title: "taikun_user Resource - terraform-provider-taikun"
subcategory: ""
description: |-   Taikun User
---

# taikun_user (Resource)

Taikun User

~> **Role Requirement** To use the `taikun_user` resource, you need a Manager or Partner account.

-> **Organization ID** `organization_id` can be specified for the Partner role, it otherwise defaults to the user's organization.

## Example Usage

```terraform
resource "taikun_user" "foo" {
  user_name = "foo"
  email     = "email@domain.fr"
  role      = "User"

  display_name     = "Foo"
  organization_id  = "42"
  disable          = true
  partner_approval = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `email` (String) The email of the user.
- `role` (String) The role of the user: `Manager` or `User`.
- `user_name` (String) The name of the user.

### Optional

- `display_name` (String) The user's display name. Defaults to ` `.
- `organization_id` (String) The ID of the user's organization.

### Read-Only

- `email_confirmed` (Boolean) Indicates whether the email of the user has been confirmed.
- `email_notification_enabled` (Boolean) Indicates whether the user has enabled notifications on their email.
- `id` (String) The UUID of the user.
- `is_approved_by_partner` (Boolean) Indicates whether the user account is approved by its Partner. If it isn't, the user won't be able to login.
- `is_csm` (Boolean) Indicates whether the user is a Customer Success Manager.
- `is_disabled` (Boolean) Indicates whether the user is locked.
- `is_owner` (Boolean) Indicates whether the user is the Owner of their organization.
- `organization_name` (String) The name of the user's organization.

## Import

Import is supported using the following syntax:

```shell
terraform import taikun_user.myuser 00000000-0000-0000-0000-000000000000
```
