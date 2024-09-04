resource "taikun_virtual_cluster" "foo" {
  name                 = "test-virtual-cluster-42"
  parent_id            = 424242
  expiration_date      = "20/01/2050"
  delete_on_expiration = "true"
}
