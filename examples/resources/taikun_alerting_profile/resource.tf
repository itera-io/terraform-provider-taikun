resource "taikun_slack_configuration" "foo" {
  name    = "foo"
  channel = "ci"
  url     = "https://hooks.myapp.example/ci"
  type    = "Alert"

  organization_id = "42"
}

resource "taikun_alerting_profile" "foo" {
  name     = "foo"
  reminder = "None"

  emails = ["test@example.com", "test@example.org", "test@example.net"]

  lock = false

  slack_configuration_id = taikun_slack_configuration.foo.id

  organization_id = "42"

  webhook {
    url = "https://www.example.com"
  }

  webhook {
    header {
      key   = "key"
      value = "value"
    }
    url = "https://www.example.com"
  }

  webhook {
    header {
      key   = "key"
      value = "value"
    }
    header {
      key   = "key2"
      value = "value"
    }
    url = "https://www.example.org"
  }

  integration {
    type  = "Opsgenie"
    url   = "https://www.opsgenie.example"
    token = "secret_token"
  }
  integration {
    type = "MicrosoftTeams"
    url  = "https://www.teams.example"
  }
  integration {
    type  = "Splunk"
    url   = "https://www.splunk.example"
    token = "secret_token"
  }
}
