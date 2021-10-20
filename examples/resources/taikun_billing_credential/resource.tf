resource "taikun_billing_credential" "foo" {
  name                = "foo"
  prometheus_password = "password"
  prometheus_url      = "url"
  prometheus_username = "username"

  organization_id = "42"
  is_locked       = true
}
