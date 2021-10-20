resource "taikun_kubernetes_profile" "foo" {
  name = "foo"

  organization_id         = "42"
  load_balancing_solution = "Taikun"
  bastion_proxy_enabled   = true
  is_locked               = true
}
