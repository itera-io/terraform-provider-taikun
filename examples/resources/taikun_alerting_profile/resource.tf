resource "taikun_alerting_profile" "foo" {
  # Required
  name     = "foo"
  reminder = "None"

  # Optional
  emails = ["test@example.com", "test@example.org", "test@example.net"]

  is_locked = false

  organization_id = resource.taikun_organization.foo.id

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
