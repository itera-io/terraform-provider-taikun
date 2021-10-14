resource "taikun_showback_rule" "foo" {
  # Required
  name        = "foo"
  price       = 1000
  metric_name = "my_metric_name"
  type        = "Sum"
  kind        = "General"

  # Optional
  label {
    key   = "key"
    value = "value"
  }
  project_alert_limit    = 42
  global_alert_limit     = 42
  organization_id        = 42
  showback_credential_id = 42
}