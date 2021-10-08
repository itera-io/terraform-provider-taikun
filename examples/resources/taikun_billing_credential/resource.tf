resource "taikun_billing_credential" "foo" {
  name            = "foo"
  organization_id = "42"

  prometheus_password = "password"
  prometheus_url      = "url"
  prometheus_username = "username"
}