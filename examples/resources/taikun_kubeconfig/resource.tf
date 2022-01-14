resource "taikun_kubeconfig" "foo" {
  project_id = "1234"

  name         = "all-can-view"
  role         = "view"
  access_scope = "all"
}

resource "local_file" "kubeconfig-foo" {
  content         = taikun_kubeconfig.foo.content
  filename        = "${path.module}/${taikun_kubeconfig.foo.project_id}-kubeconfig.yaml"
  file_permission = "0644"
}
