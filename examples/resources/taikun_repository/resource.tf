// Import and enable new private repository
resource "taikun_repository" "foo" {
  name              = "new-private-repo"
  organization_name = "Org01"
  url               = "https://example.org/Charts/"
  private           = true
  enabled           = true
}

// Enable public repository already present in Taikun
resource "taikun_repository" "bar" {
  name              = "argo"
  organization_name = "argoproj"
  private           = false
  enabled           = true
}
