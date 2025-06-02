resource "taikun_catalog_project_binding" "foo" {
  catalog_name = "new-catalog"
  project_id   = 12345
  is_bound     = true
}