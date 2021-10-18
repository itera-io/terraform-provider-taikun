resource "taikun_cloud_credential_aws" "foo" {
  # Required
  name = "foo"

  access_key_id     = "access_key_id"
  secret_access_key = "secret_access_key"
  region            = "region"
  availability_zone = "availability_zone"

  # Optional
  organization_id = "42"
  is_locked       = false
}