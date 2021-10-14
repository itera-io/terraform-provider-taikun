resource "taikun_user" "foo" {
  # Required
  user_name = "foo"
  email     = "email@domain.fr"
  role      = "User"

  # Optional
  display_name        = "Foo"
  organization_id     = "42" # Optional for Partner and Admin
  user_disabled       = true
  approved_by_partner = true
}