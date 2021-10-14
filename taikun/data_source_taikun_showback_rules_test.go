package taikun

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataSourceTaikunShowbackRulesConfig = `
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

data "taikun_showback_rules" "all" {
   depends_on = [
    taikun_showback_rule.foo
  ]
}`

func TestAccDataSourceTaikunShowbackRules(t *testing.T) {
	showbackRuleName := randomTestName()
	price := math.Round(rand.Float64()*10000) / 100
	metricName := randomString()
	typeS := []string{"Count", "Sum"}[rand.Int()%2]
	kind := "External"
	projectLimit := rand.Int31()
	globalLimit := rand.Int31()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunShowbackRulesConfig,
					showbackRuleName,
					price,
					metricName,
					typeS,
					kind,
					projectLimit,
					globalLimit,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_showback_rules.all", "showback_rules.#"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_rules.all", "showback_rules.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_rules.all", "showback_rules.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_rules.all", "showback_rules.0.metric_name"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_rules.all", "showback_rules.0.price"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_rules.all", "showback_rules.0.type"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_rules.all", "showback_rules.0.kind"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_rules.all", "showback_rules.0.project_alert_limit"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_rules.all", "showback_rules.0.global_alert_limit"),
				),
			},
		},
	})
}

const testAccDataSourceTaikunShowbackRulesWithFilterConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "Foo"
  discount_rate = 42
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
  organization_id = resource.taikun_organization.foo.id
}

data "taikun_showback_rules" "all" {
  organization_id = resource.taikun_organization.foo.id

  depends_on = [
    taikun_showback_rule.foo
  ]
}`

func TestAccDataSourceTaikunShowbackRulesWithFilter(t *testing.T) {
	organizationName := randomTestName()
	showbackRuleName := randomTestName()
	price := math.Round(rand.Float64()*10000) / 100
	metricName := randomString()
	typeS := []string{"Count", "Sum"}[rand.Int()%2]
	kind := "External"
	projectLimit := rand.Int31()
	globalLimit := rand.Int31()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceTaikunShowbackRulesWithFilterConfig,
					organizationName,
					showbackRuleName,
					price,
					metricName,
					typeS,
					kind,
					projectLimit,
					globalLimit,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.taikun_showback_rules.all", "showback_rules.0.organization_name", organizationName),
					resource.TestCheckResourceAttrSet("data.taikun_showback_rules.all", "showback_rules.#"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_rules.all", "showback_rules.0.id"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_rules.all", "showback_rules.0.name"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_rules.all", "showback_rules.0.metric_name"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_rules.all", "showback_rules.0.price"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_rules.all", "showback_rules.0.type"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_rules.all", "showback_rules.0.kind"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_rules.all", "showback_rules.0.project_alert_limit"),
					resource.TestCheckResourceAttrSet("data.taikun_showback_rules.all", "showback_rules.0.global_alert_limit"),
				),
			},
		},
	})
}
