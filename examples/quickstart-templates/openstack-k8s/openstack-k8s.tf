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

// Download available flavors and save them as a local array
data "taikun_flavors" "foo" {
  cloud_credential_id = taikun_cloud_credential_openstack.foo.id
  min_cpu             = 2
  max_cpu             = 2
  min_ram             = 4
  max_ram             = 8
}
locals {
  flavors = [for flavor in data.taikun_flavors.foo.flavors : flavor.name]
}

// Create a Taikun CloudWorks k8s project with two workers
resource "taikun_project" "foo" {
  name                = "terraform-openstack-project"
  cloud_credential_id = taikun_cloud_credential_openstack.foo.id
  flavors             = local.flavors

  server_bastion {
    name      = "os-bastion"
    disk_size = 30
    flavor    = local.flavors[0]
  }
  server_kubemaster {
    name      = "os-master"
    disk_size = 30
    flavor    = local.flavors[0]
  }
  server_kubeworker {
    name      = "os-worker-1"
    disk_size = 30
    flavor    = local.flavors[0]
  }
  server_kubeworker {
    name      = "os-worker-2"
    disk_size = 30
    flavor    = local.flavors[0]
  }
}
