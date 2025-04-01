terraform {
  required_providers {
    taikun = {
      source = "itera-io/taikun"
    }
  }
}

// Define connection to Taikun
provider "taikun" {
  email    = "test@example.com"
  password = "asdfghjkl"
}

// Define a Cloud credential (Connection to Openstack)
resource "taikun_cloud_credential_openstack" "foo" {
  name = "terraform-openstack-cc"

  user                = "zaphod"
  password            = "bEEblebrox4prez"
  url                 = "example.com"
  domain              = "domain"
  project_name        = "project_name"
  public_network_name = "public_network_name"
  region              = "region"
  lock                = false
}


// Create a BMW ccluster
resource "taikun_project" "proj1" {
  name                = "tf-acc-oneclick"
  cloud_credential_id = taikun_cloud_credential_openstack.foo.id
  flavors             = concat(local.flavors_small, local.flavors_big)

  server_bastion {
    name      = "tf-acc-bastion"
    disk_size = 30
    flavor    = local.flavors_small[0]
  }
  server_kubemaster {
    name      = "tf-acc-master"
    disk_size = 30
    flavor    = local.flavors_small[0]
  }
  server_kubeworker {
    name      = "tf-acc-worker"
    disk_size = 100
    flavor    = local.flavors_big[0]
  }
}

// Create a catalog, bind one app and the project above
resource "taikun_catalog" "cat01" {
  name        = "oneclick-catalog"
  description = "Created by terraform for oneclick deployment."
  projects    = [taikun_project.proj1.id]

  application {
    name       = "wordpress"
    repository = "taikun-managed-apps"
  }
}

// Finally create app instance
resource "taikun_app_instance" "app01" {
  name           = "terraform-wp01"
  namespace      = "wordpress-ns"
  project_id     = taikun_project.proj1.id // The project above
  catalog_app_id = local.wp_app_id         // The app selected below, from the catalog above
}

// Selecting the app (get id ofo app bound to the catalog from name and org)
locals {
  wp_app_id = [for app in tolist(taikun_catalog.cat01.application) :
  app.id if app.name == "wordpress" && app.repository == "taikun-managed-apps"][0]
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