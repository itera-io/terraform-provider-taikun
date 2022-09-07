resource "taikun_cloud_credential_azure" "foo" {
  name              = "foo"
  location          = "northeurope"
  availability_zone = "1"
}

data "taikun_images_azure" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_azure.foo.id
  publisher           = "Canonical"
  offer               = "UbuntuServer"
  sku                 = "19.04"
  latest              = true
}
