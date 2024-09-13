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
    key   = "label"
    value = "value"
  }
}
