terraform {
  required_providers {
    taikun = {
      source = "itera-io/taikun"
    }
  }
}

// Define connection to Taikun, you can load from ENV variables with TAIKUN_EMAIL, TAIKUN_PASSWORD...
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
resource "taikun_project" "proj" {
  name                = "tf-openstack-vc"
  cloud_credential_id = taikun_cloud_credential_openstack.foo.id
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

resource "taikun_virtual_cluster" "foo" {
  name      = "terraform-vc01"
  parent_id = taikun_project.proj.id
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