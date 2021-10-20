resource "taikun_backup_credential" "foo" {
  name                 = "foo"
  s3_access_key_id     = "s3_access_key_id"
  s3_secret_access_key = "s3_secret_access_key"
  s3_endpoint          = "s3_endpoint"
  s3_region            = "s3_region"

  organization_id = "42"
  is_locked       = true
}
