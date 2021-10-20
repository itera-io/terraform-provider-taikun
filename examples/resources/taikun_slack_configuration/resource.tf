resource "taikun_slack_configuration" "foo" {
  name    = "foo"
  channel = "ci"
  url     = "https://hooks.myapp.example/ci"
  type    = "Alert" // or "General"

  organization_id = "42"
}
