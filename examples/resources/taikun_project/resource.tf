resource "taikun_cloud_credential_openstack" "foo" {
  name = "foo-credentials"
}

resource "taikun_access_profile" "foo" {
  name = "foo"
}

resource "taikun_alerting_profile" "foo" {
  name     = "foo"
  reminder = "Daily"
}

resource "taikun_kubernetes_profile" "foo" {
  name = "foo"
}

resource "taikun_standalone_profile" "foo" {
  name       = "foo"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"
}

data "taikun_flavors" "small" {
  cloud_credential_id = resource.taikun_cloud_credential_openstack.foo.id
  min_cpu             = 2
  max_cpu             = 8
}

data "taikun_images" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_openstack.foo.id
}

locals {
  flavors = [for flavor in data.taikun_flavors.small.flavors : flavor.name]
  images  = [for image in data.taikun_images.foo.images : image.name]
}

resource "taikun_project" "foobar" {
  name                = "foobar"
  cloud_credential_id = resource.taikun_cloud_credential_openstack.foo.id

  access_profile_id     = resource.taikun_access_profile.foo.id
  alerting_profile_id   = resource.taikun_alerting_profile.foo.id
  kubernetes_profile_id = resource.taikun_kubernetes_profile.foo.id
  policy_profile_id     = resource.taikun_policy_profile.foo.id

  expiration_date = "21/12/2012"
  auto_upgrade    = true
  monitoring      = true

  quota_cpu_units = 64
  quota_disk_size = 1024
  quota_ram_size  = 256

  flavors = local.flavors
  images  = local.images

  vm {
    name        = "b"
    volume_size = 30

    flavor   = local.flavors[0]
    image_id = local.images[0]

    cloud_init            = ""
    standalone_profile_id = resource.taikun_standalone_profile.foo.id
    public_ip             = true

    // OpenStack
    volume_type = "ssd-2000iops"

    tag {
      key   = "key"
      value = "value"
    }

    disk {
      name = "name"
      size = 30

      // AWS
      device_name = "/dev/sda3"
      // OpenStack
      volume_type = "ssd-2000iops"
      // Azure
      lun_id = 3
    }
  }

  server_bastion {
    name   = "b"
    flavor = local.flavors[0]
  }
  server_kubemaster {
    name   = "m"
    flavor = local.flavors[0]
  }
  server_kubeworker {
    name      = "w"
    flavor    = local.flavors[0]
    disk_size = 30
  }
}
