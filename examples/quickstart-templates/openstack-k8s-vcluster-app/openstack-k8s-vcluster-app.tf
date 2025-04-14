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
  # email    = "test@example.com"
  # password = "asdfghjkl"
}

resource "taikun_organization" "foo" {
  name          = "tf-acc-manual01"
  full_name     = "tf-acc-manual01"
  discount_rate = 100
}

// Define a Cloud credential (Connection to Openstack)
resource "taikun_cloud_credential_openstack" "foo" {
  name            = "terraform-openstack-cc"
  organization_id = taikun_organization.foo.id

  # user                = "zaphod"
  # password            = "bEEblebrox4prez"
  # url                 = "example.com"
  # domain              = "domain"
  # project_name        = "project_name"
  # public_network_name = "public_network_name"
  # region              = "region"
  # lock                = false
}

// Create a BMW ccluster
resource "taikun_project" "proj" {
  name                = "tf-openstack-vc"
  cloud_credential_id = taikun_cloud_credential_openstack.foo.id
  organization_id     = taikun_organization.foo.id
  flavors             = concat(local.flavors_small, local.flavors_big)

  server_bastion {
    name      = "tf-bastion"
    disk_size = 30
    flavor    = local.flavors_small[0]
  }
  server_kubemaster {
    name      = "tf-master"
    disk_size = 30
    flavor    = local.flavors_small[0]
  }
  server_kubeworker {
    name      = "tf-worker"
    disk_size = 100
    flavor    = local.flavors_big[0]
  }
}

// Create a new virtual cluster
resource "taikun_virtual_cluster" "foo" {
  name      = "tf-acc-vc01"
  parent_id = taikun_project.proj.id
}

# // Make sure the taikun-managed-apps repository is enabled
resource "taikun_repository" "bar" {
  name              = "taikun-managed-apps"
  organization_name = "taikun"
  organization_id   = taikun_organization.foo.id
  private           = false
  enabled           = true
}

// Create a catalog, bind one app and the project above
resource "taikun_catalog" "cat01" {
  name            = "mock-catalog"
  description     = "Created by terraform for mock deployment."
  organization_id = taikun_organization.foo.id

  application {
    name       = "apache"
    repository = "taikun-managed-apps"
  }

  // The repository must be enabled at this point.
  depends_on = [
    taikun_repository.bar
  ]
}

// Bind the created project to the created catalog
resource "taikun_catalog_project_binding" "bind01" {
  catalog_name    = "taikun-managed-apps"
  project_id      = taikun_virtual_cluster.foo.id
  organization_id = taikun_organization.foo.id
  is_bound        = true
}

// Finally create app instance
resource "taikun_app_instance" "app01" {
  name           = "terraform-apache01"
  namespace      = "apache-ns"
  project_id     = taikun_virtual_cluster.foo.id // The virtual cluster project above
  catalog_app_id = local.apache_app_id           // The app selected below, from the catalog above
  depends_on = [
    taikun_catalog_project_binding.bind01
  ]
}

// Selecting the app (get id of app bound to the catalog from name and org)
locals {
  apache_app_id = [for app in tolist(taikun_catalog.cat01.application) :
  app.id if app.name == "apache" && app.repository == "taikun-managed-apps"][0]
}

// Flavors configuration. Filtering data sources
data "taikun_flavors" "small" {
  cloud_credential_id = taikun_cloud_credential_openstack.foo.id
  min_cpu             = 2
  max_cpu             = 2
  min_ram             = 4
  max_ram             = 4
}
data "taikun_flavors" "big" {
  cloud_credential_id = taikun_cloud_credential_openstack.foo.id
  min_cpu             = 4
  max_cpu             = 4
  min_ram             = 8
  max_ram             = 8
}
locals {
  flavors_small = [for flavor in data.taikun_flavors.small.flavors : flavor.name]
  flavors_big   = [for flavor in data.taikun_flavors.big.flavors : flavor.name]
}