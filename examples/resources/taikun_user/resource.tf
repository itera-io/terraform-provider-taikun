resource "taikun_user" "foo" {
  user_name = "foo"
  email     = "email@domain.fr"
  role      = "User"

  display_name     = "Foo"
  organization_id  = "42"
  disable          = true
  partner_approval = true
}
