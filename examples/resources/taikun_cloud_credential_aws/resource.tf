resource "taikun_cloud_credential_aws" "foo" {
  name = "foo"

  access_key_id     = "access_key_id"
  secret_access_key = "secret_access_key"
  region            = "region"

  organization_id = "42"
  lock            = false
}
