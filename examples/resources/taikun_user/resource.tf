resource "taikun_user" "foo" {
  user_name        = "foo"
  email            = "email@domain.fr"
  display_name     = "Foo"
  account_id       = "42"
  is_account_admin = false
}
