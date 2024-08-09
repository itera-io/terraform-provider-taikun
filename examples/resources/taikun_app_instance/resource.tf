resource "taikun_catalog" "foo" {
  name        = "new-catalog"
  description = "Created by Terraform"
  projects    = ["37415"]

  application {
    name       = "wordpress"
    repository = "taikun-managed-apps"
  }
}

resource "taikun_app_instance" "foo" {
  name           = "wordpress01"
  namespace      = "wordpress01-ns"
  project_id     = "37415"
  catalog_app_id = local.app_id
}


locals {
  app_id = [for app in tolist(taikun_catalog.foo.application) : app.id if app.name == "wordpress" && app.repository == "taikun-managed-apps"][0]
}