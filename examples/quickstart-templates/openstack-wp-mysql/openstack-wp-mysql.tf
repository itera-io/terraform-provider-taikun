terraform {
  required_providers {
    taikun = {
      source  = "itera-io/taikun"
      version = "1.9.1" // Use the latest version
    }
  }
}

// Define connection to Taikun
provider "taikun" {
  email    = "example@taikun.cloud"
  password = "your-password"
}

// Define a Cloud credential (Connection to Openstack)
resource "taikun_cloud_credential_openstack" "wordpress_credentials" {
  name                = "wordpress-openstack-cc"
  user                = "your-username"
  password            = "your-password"
  url                 = "openstack-url"
  domain              = "your-domain"
  project_name        = "your-project"
  public_network_name = "public-network"
  region              = "your-region"
  lock                = false
}
// Create a kubernetes cluster
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

// Create a catalog with WordPress and MySQL applications
resource "taikun_catalog" "cat01" {
  name        = "oneclick-catalog"
  description = "Created by terraform for oneclick deployment."

  application {
    name       = "wordpress"
    repository = "taikun-managed-apps"
  }
  application {
    name       = "mysql"
    repository = "taikun-managed-apps"
  }
}

// Bind the created project to the created catalog
resource "taikun_catalog_project_binding" "bind01" {
  catalog_name = taikun_catalog.cat01.name
  project_id   = taikun_project.proj1.id
  is_bound     = true
}

// Create MySQL instance first
resource "taikun_app_instance" "mysql01" {
  name           = "terraform-mysql01"
  namespace      = "wordpress-ns"
  project_id     = taikun_project.proj1.id
  catalog_app_id = local.mysql_app_id
  depends_on     = [taikun_project.proj1]
}

// Create WordPress instance with MySQL connection
resource "taikun_app_instance" "wordpress01" {
  name           = "terraform-wp01"
  namespace      = "wordpress-ns"
  project_id     = taikun_project.proj1.id
  catalog_app_id = local.wp_app_id
  depends_on     = [taikun_project.proj1, taikun_app_instance.mysql01]
}

// App ID selectors
locals {
  wp_app_id = [for app in tolist(taikun_catalog.cat01.application) :
  app.id if app.name == "wordpress" && app.repository == "taikun-managed-apps"][0]

  mysql_app_id = [for app in tolist(taikun_catalog.cat01.application) :
  app.id if app.name == "mysql" && app.repository == "taikun-managed-apps"][0]
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