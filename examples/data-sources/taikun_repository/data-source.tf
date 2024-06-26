data "taikun_repository" "foo" {
  name              = "private-repo-test"
  organization_name = "Org1"
  private           = "true"
}