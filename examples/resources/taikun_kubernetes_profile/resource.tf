resource "taikun_kubernetes_profile" "foo" {
  # Required
  name = "foo"

  # Optional
  organization_id         = "42"
  load_balancing_solution = "Taikun"
  bastion_proxy_enabled   = true
  is_locker               = true
}