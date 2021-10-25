package taikun

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunShowbackRuleConfig = `
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
}

data "taikun_showback_rule" "foo" {
  id = resource.taikun_showback_rule.foo.id
}
`

func TestAccDataSourceTaikunShowbackRule(t *testing.T) {
	showbackRuleName := randomTestName()
	price := math.Round(rand.Float64()*10000) / 100
	metricName := randomString()
	typeS := []string{"Count", "Sum"}[rand.Int()%2]
	kind := "External"
	projectLimit := rand.Int31()
	globalLimit := rand.Int31()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunShowbackRuleConfig,
					showbackRuleName,
					price,
					metricName,
					typeS,
					kind,
					projectLimit,
					globalLimit,
				),
				Check: checkDataSourceStateMatchesResourceState(
					"data.taikun_showback_rule.foo",
					"taikun_showback_rule.foo",
				),
			},
		},
	})
}
