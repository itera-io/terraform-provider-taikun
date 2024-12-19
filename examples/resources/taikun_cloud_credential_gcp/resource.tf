resource "taikun_cloud_credential_gcp" "foo" {
  name = "foo"

  config_file    = "./config.json"
  region         = "europe-north1"
  import_project = true

  organization_id = "42"
  lock            = false
}

data "taikun_images_gcp" "foo" {
  cloud_credential_id = taikun_cloud_credential_gcp.foo.id
  type                = "ubuntu"
  latest              = true
}

locals {
  images = [for image in data.taikun_images_gcp.foo.images : image.name] // GCP uses image names, not image ids.
}


resource "taikun_project" "foo" {
  name                = "mock-project"
  cloud_credential_id = taikun_cloud_credential_gcp.foo.id
  images              = local.images
}