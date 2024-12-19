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

// Define the Cloud credential that connects to ZADARA
resource "taikun_cloud_credential_zadara" "foo" {
  name              = "terraform-zadara-cc"
  access_key_id     = "1234567890"
  secret_access_key = "asdfghjkl"
  volume_type       = "standard"
  url               = "example.com"
  organization_id   = "42"
  region            = "symphony"
}

// Get flavors for this Cloud credential
data "taikun_flavors" "foo" {
  cloud_credential_id = taikun_cloud_credential_zadara.foo.id
  min_cpu             = 2
  max_cpu             = 2
  min_ram             = 4
  max_ram             = 8
}

locals {
  flavors = [for flavor in data.taikun_flavors.foo.flavors : flavor.name]
}

// Your public key - you are able to connect to the bastion with this key
resource "taikun_standalone_profile" "foo" {
  name       = "tftest-zadara"
  public_key = "asdfghjkl"
}

// Create a CloudWorks project with one VM
resource "taikun_project" "foo" {
  name                = "tftest-zadara"
  cloud_credential_id = taikun_cloud_credential_zadara.foo.id
  flavors             = local.flavors
  images              = ["ami-688a0d59ac354bfe9081e9d34532ff25"]

  vm {
    name                  = "tftest-zadara"
    flavor                = local.flavors[0]
    image_id              = "ami-688a0d59ac354bfe9081e9d34532ff25"
    standalone_profile_id = taikun_standalone_profile.foo.id
    volume_size           = 30
  }
}
