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

data "taikun_flavors" "foo" {
  cloud_credential_id = taikun_cloud_credential_openstack.foo.id
  min_cpu             = 2
  max_cpu             = 2
  min_ram             = 4
  max_ram             = 8
}

data "taikun_images_openstack" "foo" {
  cloud_credential_id = taikun_cloud_credential_openstack.foo.id
}

locals {
  images  = [for image in data.taikun_images_openstack.foo.images : image.id]
  flavors = [for flavor in data.taikun_flavors.foo.flavors : flavor.name]
}

resource "taikun_standalone_profile" "foo" {
  name       = "tf-acc-radekmanual-sp2-os"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIDVfC6srYMVRqcXWU8JYGWQFPO2UZGsE897S7vEubty8 radek.smid@taikun.cloud"
}

resource "taikun_project" "foo" {
  name                = "tf-acc-radekmanual-p1-vs"
  cloud_credential_id = taikun_cloud_credential_openstack.foo.id
  flavors             = local.flavors
  images              = local.images

  quota_vm_cpu_units   = 64
  quota_vm_ram_size    = 256
  quota_vm_volume_size = 512

  vm {
    name   = "tf-acc-vm"
    flavor = local.flavors[0]
    //image_id = local.images[0] // Parametrized selection of image
    image_id              = "9414e6a0-d67c-4c92-a47e-82c697067530" // Hardcoded selection of image
    standalone_profile_id = taikun_standalone_profile.foo.id
    volume_size           = 60

    disk {
      name = "tf-acc-disk"
      size = 30
    }
    disk {
      name        = "tf-acc-disk2"
      size        = 30
      volume_type = "ssd"
    }
    tag {
      key   = "key"
      value = "value"
    }
    tag {
      key   = "key2"
      value = "value"
    }
  }
}