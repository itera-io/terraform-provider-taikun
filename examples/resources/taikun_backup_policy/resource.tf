resource "taikun_backup_policy" "foo" {
  name       = "foo"
  project_id = taikun_project.foo.id

  cron_period      = "0 0 * * 0"
  retention_period = "72h"

  included_namespaces = ["name"]
}