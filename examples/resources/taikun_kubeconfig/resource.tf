resource "taikun_kubeconfig" "foo" {
  project_id = "1234"

  name         = "foo"
  role         = "edit"
  access_scope = "managers"
}
