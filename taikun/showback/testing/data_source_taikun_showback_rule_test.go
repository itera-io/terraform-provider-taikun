package testing

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"testing"

	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceTaikunShowbackRuleConfig = `
resource "taikun_showback_credential" "foo" {
  name            = "%s"

  password = "%s"
  url = "%s"
  username = "%s"
}


resource "taikun_showback_rule" "foo" {
  name = "%s"
  price = %f
  metric_name = "%s"
  type = "%s"
  kind = "%s"
  label {
    key = "key"
    value = "value"
  }
  project_alert_limit = %d
  global_alert_limit = %d
  showback_credential_id = resource.taikun_showback_credential.foo.id
}

data "taikun_showback_rule" "foo" {
  id = resource.taikun_showback_rule.foo.id
}
`

func TestAccDataSourceTaikunShowbackRule(t *testing.T) {
	showbackCredentialName := utils.RandomTestName()
	showbackRuleName := utils.RandomTestName()
	price := math.Round(rand.Float64()*10000) / 100
	metricName := utils.RandomString()
	typeS := []string{"Count", "Sum"}[rand.Int()%2]
	kind := "External"
	projectLimit := rand.Int31()
	globalLimit := rand.Int31()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunShowbackRuleConfig,
					showbackCredentialName,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
					showbackRuleName,
					price,
					metricName,
					typeS,
					kind,
					projectLimit,
					globalLimit,
				),
				Check: utils_testing.CheckDataSourceStateMatchesResourceState(
					"data.taikun_showback_rule.foo",
					"taikun_showback_rule.foo",
				),
			},
		},
	})
}
