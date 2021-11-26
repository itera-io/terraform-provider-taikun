resource "taikun_policy_profile" "foo" {
  name = "foo"

  forbid_node_port        = true
  forbid_http_ingress     = true
  require_probe           = true
  unique_ingress          = true
  unique_service_selector = true

  allowed_repos = [
    "repo"
  ]
  forbidden_tags = [
    "tag"
  ]
  ingress_whitelist = [
    "ingress"
  ]

  organization_id = "42"
  lock            = true
}
