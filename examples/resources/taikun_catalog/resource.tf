resource "taikun_catalog" "foo" {
  name        = "new_catalog"
  description = "Created by Terraform"
  projects    = ["37415"]

  application {
    name       = "wordpress"
    repository = "taikun-managed-apps"
  }
}