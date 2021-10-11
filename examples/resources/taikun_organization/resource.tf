# Only valid for roles with permission to create/delete/modify Organizations
# i.e. Partner and Admin roles
resource "taikun_organization" "foo" {
  # Required
  name          = "foo"
  full_name     = "Foo Organization"
  discount_rate = 42

  # Optional
  vat_number                       = "CZ4495374355"
  email                            = "contact@foo.org"
  billing_email                    = "billing@foo.org"
  phone                            = "+420123456789"
  address                          = "Foo 42"
  city                             = "Praha"
  country                          = "Czechia"
  is_locked                        = false
  let_managers_change_subscription = true
}
