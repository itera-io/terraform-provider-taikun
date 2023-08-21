resource "taikun_user" "foo" {
  user_name = "foo"
  email     = "email"
  role      = "User"
}

resource "taikun_cloud_credential_aws" "foo" {
  name              = "foo"
}

resource "taikun_project" "foo" {
  name                = "foo"
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id
}

resource "taikun_project_user_attachment" "foo" {
  project_id = resource.taikun_project.foo.id
  user_id    = resource.taikun_user.foo.id
}
