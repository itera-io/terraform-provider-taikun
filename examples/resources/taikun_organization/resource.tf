resource "taikun_organization" "foo" {
  name          = "foo"
  full_name     = "Foo Organization"
  discount_rate = 42

  vat_number                       = "CZ4495374355"
  email                            = "contact@foo.org"
  billing_email                    = "billing@foo.org"
  phone                            = "+420123456789"
  address                          = "Foo 42"
  city                             = "Praha"
  country                          = "Czechia"
  managers_can_change_subscription = true
}
