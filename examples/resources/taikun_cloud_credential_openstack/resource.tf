resource "taikun_cloud_credential_openstack" "foo" {
  name = "foo"

  user     = "user"
  password = "password"
  url      = "url"
  domain   = "domain"

  project_name        = "project_name"
  public_network_name = "public_network_name"
  region              = "region"

  availability_zone          = "availability_zone"
  volume_type_name           = "volume_type_name"
  imported_network_subnet_id = "imported_network_subnet_id"

  organization_id = "42"
  is_locked       = false
}
