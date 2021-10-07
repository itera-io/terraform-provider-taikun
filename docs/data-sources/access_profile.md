---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "taikun_access_profile Data Source - terraform-provider-taikun"
subcategory: ""
description: |-
  Get an access profiles by its id.
---

# taikun_access_profile (Data Source)

Get an access profiles by its id.

## Example Usage

```terraform
data "taikun_access_profile" "foo" {
  id = "42"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **id** (String) The id of the access profile.

### Read-Only

- **created_by** (String) The creator of the access profile.
- **dns_server** (List of Object) List of DNS servers. (see [below for nested schema](#nestedatt--dns_server))
- **http_proxy** (String) HTTP Proxy of the access profile.
- **is_locked** (Boolean) Indicates whether the access profile is locked or not.
- **last_modified** (String) Time of last modification.
- **last_modified_by** (String) The last user who modified the access profile.
- **name** (String) The name of the access profile.
- **ntp_server** (List of Object) List of NTP servers. (see [below for nested schema](#nestedatt--ntp_server))
- **organization_id** (String) The id of the organization which owns the access profile.
- **organization_name** (String) The name of the organization which owns the access profile.
- **projects** (List of Object) List of associated projects. (see [below for nested schema](#nestedatt--projects))
- **ssh_user** (List of Object) List of SSH Users. (see [below for nested schema](#nestedatt--ssh_user))

<a id="nestedatt--dns_server"></a>
### Nested Schema for `dns_server`

Read-Only:

- **address** (String)
- **id** (String)


<a id="nestedatt--ntp_server"></a>
### Nested Schema for `ntp_server`

Read-Only:

- **address** (String)
- **id** (String)


<a id="nestedatt--projects"></a>
### Nested Schema for `projects`

Read-Only:

- **id** (String)
- **name** (String)


<a id="nestedatt--ssh_user"></a>
### Nested Schema for `ssh_user`

Read-Only:

- **id** (String)
- **name** (String)
- **public_key** (String)

