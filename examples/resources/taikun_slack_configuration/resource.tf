resource "taikun_slack_configuration" "foo" {
  # Required
  name    = "foo"
  channel = "ci"
  url     = "https://hooks.myapp.example/ci"
  type    = "Alert" // or "General"

  # Optional for Partner and Admin
  organization_id = "42"
}
