resource "taikun_showback_credential" "foo" {
  name     = "foo"
  password = "password"
  url      = "url"
  username = "username"

  organization_id = "42"
  lock            = true
}
