resource "taikun_cloud_credential_aws" "foo" {
  name = "foo"
}

data "taikun_images_aws" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id
  latest              = true
  owners              = ["Canonical"]
}
