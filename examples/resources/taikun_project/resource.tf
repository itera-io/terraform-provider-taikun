resource "taikun_cloud_credential_aws" "foo" {
  name              = "foo-credentials"
  availability_zone = "a"
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

data "taikun_flavors" "small" {
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id
  min_cpu             = 2
  max_cpu             = 8
}

locals {
  flavors = [for flavor in data.taikun_flavors.small.flavors : flavor.name]
}

resource "taikun_project" "foobar" {
  name                = "foobar"
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id

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
