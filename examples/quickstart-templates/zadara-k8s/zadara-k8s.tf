terraform {
  required_providers {
    taikun = {
      source = "itera-io/taikun"
      version = "1.9.1"
    }
  }
}

// Define connection to Taikun
provider "taikun" {
  email    = "test@example.com"
  password = "asdfghjkl"
}

// Define the Cloud credential that connects to ZADARA
resource "taikun_cloud_credential_zadara" "foo" {
  name		    = "terraform-zadara-cc"
  access_key_id     = "1234567890"
  secret_access_key = "asdfghjkl"
  volume_type       = "standard"
  url               = "example.com"
  organization_id   = "42"
  region            = "symphony"
}

// Download available flavors and save them as a local array
data "taikun_flavors" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_zadara.foo.id
  min_cpu = 2
  max_cpu = 2
  min_ram = 4
  max_ram = 8
}
locals {
  flavors = [for flavor in data.taikun_flavors.foo.flavors: flavor.name]
}

// Create a Taikun CloudWorks k8s project with two workers
resource "taikun_project" "foo" {
  name = "terraform-zadara-project"
  cloud_credential_id = resource.taikun_cloud_credential_zadara.foo.id
  flavors = local.flavors

  server_bastion {
    name = "zadara-bastion"
    disk_size = 30
    flavor = local.flavors[0]
  }
  server_kubemaster {
    name = "zadara-master"
    disk_size = 30
    flavor = local.flavors[0]
  }
  server_kubeworker {
    name = "zadara-worker-1"
    disk_size = 30
    flavor = local.flavors[0]
  }
  server_kubeworker {
    name = "zadara-worker-2"
    disk_size = 30
    flavor = local.flavors[0]
  }
}
