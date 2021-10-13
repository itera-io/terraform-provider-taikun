resource "taikun_showback_credential" "foo" {
  # Required
  name     = "foo"
  password = "password"
  url      = "url"
  username = "username"

  # Optional
  organization_id = "42" # Optional for Partner and Admin
  is_locked       = true
}