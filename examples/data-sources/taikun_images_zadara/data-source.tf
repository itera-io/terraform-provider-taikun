resource "taikun_cloud_credential_zadara" "foo" {
  name = "foo"

  access_key_id     = "access_key_id"
  secret_access_key = "secret_access_key"
  region            = "region"
  volume_type       = "standard"
  url               = "example.com"

  organization_id = "42"
  lock            = false
}

data "taikun_images_zadara" "foo" {
  cloud_credential_id = taikun_cloud_credential_zadara.foo.id
  latest              = false
}
