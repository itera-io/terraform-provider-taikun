resource "taikun_billing_credential" "foo" {
  name                = "foo"
  prometheus_password = "password"
  prometheus_url      = "url"
  prometheus_username = "username"
}

resource "taikun_billing_rule" "foo" {
  name        = "foo"
  metric_name = "coredns_forward_request_duration_seconds"
  price       = 1
  type        = "Sum"

  billing_credential_id = taikun_billing_credential.foo.id

  label {
    key   = "key"
    value = "value"
  }
}

resource "taikun_organization" "foo" {
  name          = "foo"
  full_name     = "foo"
  discount_rate = 100
}

resource "taikun_organization_billing_rule_attachment" "foo" {
  billing_rule_id = taikun_billing_rule.foo.id
  organization_id = taikun_organization.foo.id

  discount_rate = 100
}
