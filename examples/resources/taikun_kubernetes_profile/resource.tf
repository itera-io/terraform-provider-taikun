resource "taikun_kubernetes_profile" "foo" {
  name = "foo"

  organization_id         = "42"
  load_balancing_solution = "Taikun"
  bastion_proxy_enabled   = true
  schedule_on_master      = false
  is_locked               = true
}
