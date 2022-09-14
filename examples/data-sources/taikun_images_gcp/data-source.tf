resource "taikun_cloud_credential_gcp" "foo" {
  name               = "foo"
  config_file        = "./gcp.json"
  billing_account_id = "000000000000"
  folder_id          = "000000-000000-000000"
  region             = "asia-northeast1"
  zone               = "asia-northeast1-b"
}

data "taikun_images_gcp" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_gcp.foo.id
  type                = "windows"
}
