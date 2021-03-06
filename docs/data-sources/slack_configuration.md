---
page_title: "taikun_slack_configuration Data Source - terraform-provider-taikun"
subcategory: ""
description: |-   Get a Slack configuration by its ID.
---

# taikun_slack_configuration (Data Source)

Get a Slack configuration by its ID.

~> **Role Requirement** To use the `taikun_slack_configuration` data source, you need a Manager or Partner account.

## Example Usage

```terraform
data "taikun_slack_configuration" "foo" {
  id = "42"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The Slack configuration's ID.

### Read-Only

- `channel` (String) Slack channel for notifications.
- `name` (String) The Slack configuration's name.
- `organization_id` (String) The ID of the organization which owns the Slack configuration.
- `organization_name` (String) The name of the organization which owns the Slack configuration.
- `type` (String) The type of notifications to receive: `Alert` (only alert-type notifications) or `General` (all notifications).
- `url` (String) Webhook URL from Slack app.


