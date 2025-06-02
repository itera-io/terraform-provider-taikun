terraform {
  required_providers {
    taikun = {
      source = "itera-io/taikun"
    }
  }
}

// Define connection to Taikun, you can load from ENV variables with TAIKUN_EMAIL, TAIKUN_PASSWORD...
provider "taikun" {
  # email    = "test@example.com"
  # password = "asdfghjkl"
}

// Create a new virtual cluster
resource "taikun_virtual_cluster" "foo" {
  name      = "terraform-vc01"
  parent_id = 4216
}

# // Make sure the taikun-managed-apps repository is enabled
resource "taikun_repository" "bar" {
  name              = "taikun-managed-apps"
  organization_name = "taikun"
  organization_id   = 2160
  private           = false
  enabled           = true
}

// Create a catalog, bind one app and the project above
resource "taikun_catalog" "cat01" {
  name            = "catalog-for-vc"
  description     = "Created by terraform for mock deployment."
  organization_id = 2160

  application {
    name       = "apache"
    repository = "taikun-managed-apps"
  }

  depends_on = [
    taikun_repository.bar
  ]
}

// Bind the created project to the created catalog
resource "taikun_catalog_project_binding" "bind01" {
  catalog_name    = taikun_catalog.cat01.name
  project_id      = taikun_virtual_cluster.foo.id
  organization_id = 2160
  is_bound        = true
}

// Finally create app instance
resource "taikun_app_instance" "app01" {
  name           = "terraform-apache01"
  namespace      = "apache-ns"
  project_id     = taikun_virtual_cluster.foo.id // The project above
  catalog_app_id = local.apache_app_id           // The app selected below, from the catalog above
  depends_on = [
    taikun_catalog_project_binding.bind01
  ]
}

// Selecting the app (get id ofo app bound to the catalog from name and org)
locals {
  apache_app_id = [for app in tolist(taikun_catalog.cat01.application) :
  app.id if app.name == "apache" && app.repository == "taikun-managed-apps"][0]
}