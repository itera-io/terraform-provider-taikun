resource "taikun_showback_credential" "foo" {
  name     = "foo"
  password = "password"
  url      = "url"
  username = "username"

  organization_id = "42" # Optional for Partner and Admin
  is_locked       = true
}
