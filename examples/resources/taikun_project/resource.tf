resource "taikun_cloud_credential_aws" "foo" {
  name              = "foo-credentials"
  availability_zone = "a"
}

resource "taikun_access_profile" "foo" {
  name = "foo"
}

resource "taikun_alerting_profile" "foo" {
  name     = "foo"
  reminder = "Daily"
}

resource "taikun_kubernetes_profile" "foo" {
  name = "foo"
}

resource "taikun_project" "foobar" {
  name                = "foobar"
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id

  access_profile_id     = resource.taikun_access_profile.foo.id
  alerting_profile_id   = resource.taikun_alerting_profile.foo.id
  kubernetes_profile_id = resource.taikun_kubernetes_profile.foo.id

  expiration_date     = "21/12/2012"
  enable_auto_upgrade = true
  enable_monitoring   = true
}
