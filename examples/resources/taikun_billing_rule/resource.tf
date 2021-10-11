resource "taikun_billing_credential" "foo" {
  name            = "foo"
  organization_id = "42"

  prometheus_password = "password"
  prometheus_url      = "url"
  prometheus_username = "username"
}

resource "taikun_billing_rule" "foo" {
  name        = "foo"
  metric_name = "coredns_forward_request_duration_seconds"
  price       = 1
  type        = "Sum"

  billing_credential_id = resource.taikun_billing_credential.foo.id

  label {
    key   = "label"
    value = "value"
  }
}