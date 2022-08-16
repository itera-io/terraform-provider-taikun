resource "taikun_cloud_credential_openstack" "foo" {
  name = "foo"
}

data "taikun_images_openstack" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_openstack.foo.id
}
