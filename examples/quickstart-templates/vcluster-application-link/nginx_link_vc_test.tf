# This test should be run using manager, connected with access keys in staging before we promote.

terraform {
  required_providers {
    taikun = {
      # source = "itera-io/taikun"
      source = "itera-io/dev/taikun"
    }
  }
}

// Define connection to Taikun
provider "taikun" {
}


// Create a new virtual cluster
resource "taikun_virtual_cluster" "foo" {
  name      = "tf-acc-vc03"
  parent_id = 42 // << Fill this
}

# // Make sure the taikun-managed-apps repository is enabled
resource "taikun_repository" "bar" {
  name              = "taikun-managed-apps"
  organization_name = "taikun"
  private           = false
  enabled           = true
}

// Create a catalog, bind one app and the project above
resource "taikun_catalog" "cat01" {
  name        = "vc-catalog03"
  description = "Created by terraform for mock deployment."

  application {
    name       = "nginx"
    repository = "taikun-managed-apps"
  }

  // The repository must be enabled at this point.
  depends_on = [
    taikun_repository.bar
  ]
}

// Bind the created project to the created catalog
resource "taikun_catalog_project_binding" "bind01" {
  catalog_name = "taikun-managed-apps"
  project_id   = taikun_virtual_cluster.foo.id
  is_bound     = true
}

// Finally create app instance
resource "taikun_app_instance" "app01" {
  name           = "vc-link-01"
  namespace      = "nginx-ns"
  project_id     = taikun_virtual_cluster.foo.id // The virtual cluster project above
  catalog_app_id = local.apache_app_id           // The app selected below, from the catalog above

  depends_on = [
    taikun_catalog_project_binding.bind01
  ]

  autosync    = false
  taikun_link = true
  parameters_base64 = base64encode(
    file("nginx.yaml")
  )
}

// You can also use the datasource to list apps in organization
data "taikun_app_instances" "foo" {}

// Selecting the app (get id of app bound to the catalog from name and org)
locals {
  apache_app_id = [for app in tolist(taikun_catalog.cat01.application) :
  app.id if app.name == "nginx" && app.repository == "taikun-managed-apps"][0]
}

// Print the URL of the app
output "taikun_link_url" {
  description = "The Taikun link URL of the deployed app instance."
  value       = resource.taikun_app_instance.app01.taikun_link_url
}