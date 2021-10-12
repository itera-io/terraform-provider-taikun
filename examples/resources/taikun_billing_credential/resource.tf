resource "taikun_billing_credential" "foo" {
  # Required
  name                = "foo"
  prometheus_password = "password"
  prometheus_url      = "url"
  prometheus_username = "username"

  # Optional
  organization_id = "42" # Optional for Partner and Admin
  is_locked       = true
}