resource "taikun_cloud_credential_openstack" "foo" {
  name = "foo"
}

data "taikun_flavors" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_openstack.foo.id

  min_cpu = 8
  max_cpu = 16
  min_ram = 32
  max_ram = 256
}
