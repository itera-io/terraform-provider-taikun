package taikun

import (
	"context"
	"errors"
	"fmt"
	tk "github.com/chnyda/taikungoclient"
	"math"
	"math/rand"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccResourceTaikunShowbackRuleConfig = `
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
`

func TestAccResourceTaikunShowbackRule(t *testing.T) {
	name := randomTestName()
	price := math.Round(rand.Float64()*10000) / 100
	metricName := randomString()
	typeS := []string{"Count", "Sum"}[rand.Int()%2]
	kind := []string{"General", "External"}[rand.Int()%2]
	projectLimit := rand.Int31()
	globalLimit := rand.Int31()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunShowbackRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunShowbackRuleConfig,
					name,
					price,
					metricName,
					typeS,
					kind,
					projectLimit,
					globalLimit),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunShowbackRuleExists,
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "metric_name", metricName),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "price", fmt.Sprint(price)),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "type", typeS),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "kind", kind),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "project_alert_limit", fmt.Sprint(projectLimit)),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "global_alert_limit", fmt.Sprint(globalLimit)),
					resource.TestCheckNoResourceAttr("taikun_showback_rule.foo", "showback_credential_id"),
					resource.TestCheckNoResourceAttr("taikun_showback_rule.foo", "showback_credential_name"),
				),
			},
			{
				ResourceName:      "taikun_showback_rule.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceTaikunShowbackRuleUpdate(t *testing.T) {
	name := randomTestName()
	price := math.Round(rand.Float64()*10000) / 100
	metricName := randomString()
	typeS := "Count"
	kind := "General"
	projectLimit := rand.Int31()
	globalLimit := rand.Int31()

	newName := randomTestName()
	newPrice := math.Round(rand.Float64()*10000) / 100
	newMetricName := randomString()
	newTypeS := "Sum"
	newKind := "External"
	newProjectLimit := rand.Int31()
	newGlobalLimit := rand.Int31()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunShowbackRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunShowbackRuleConfig,
					name,
					price,
					metricName,
					typeS,
					kind,
					projectLimit,
					globalLimit,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunShowbackRuleExists,
					resource.TestCheckResourceAttrSet("taikun_showback_rule.foo", "id"),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "metric_name", metricName),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "price", fmt.Sprint(price)),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "type", typeS),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "kind", kind),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "project_alert_limit", fmt.Sprint(projectLimit)),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "global_alert_limit", fmt.Sprint(globalLimit)),
					resource.TestCheckNoResourceAttr("taikun_showback_rule.foo", "showback_credential_id"),
					resource.TestCheckNoResourceAttr("taikun_showback_rule.foo", "showback_credential_name"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunShowbackRuleConfig,
					newName,
					newPrice,
					newMetricName,
					newTypeS,
					newKind,
					newProjectLimit,
					newGlobalLimit,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunShowbackRuleExists,
					resource.TestCheckResourceAttrSet("taikun_showback_rule.foo", "id"),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "name", newName),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "metric_name", newMetricName),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "price", fmt.Sprint(newPrice)),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "type", newTypeS),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "kind", newKind),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "project_alert_limit", fmt.Sprint(newProjectLimit)),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "global_alert_limit", fmt.Sprint(newGlobalLimit)),
					resource.TestCheckNoResourceAttr("taikun_showback_rule.foo", "showback_credential_id"),
					resource.TestCheckNoResourceAttr("taikun_showback_rule.foo", "showback_credential_name"),
				),
			},
		},
	})
}

const testAccResourceTaikunShowbackRuleWithCredentialsConfig = `
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
`

func TestAccResourceTaikunShowbackRuleWithCredentials(t *testing.T) {
	showbackCredentialName := randomTestName()
	name := randomTestName()
	price := math.Round(rand.Float64()*10000) / 100
	metricName := randomString()
	typeS := []string{"Count", "Sum"}[rand.Int()%2]
	kind := "External"
	projectLimit := rand.Int31()
	globalLimit := rand.Int31()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunShowbackRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunShowbackRuleWithCredentialsConfig,
					showbackCredentialName,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
					name,
					price,
					metricName,
					typeS,
					kind,
					projectLimit,
					globalLimit,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunShowbackRuleExists,
					resource.TestCheckResourceAttrSet("taikun_showback_rule.foo", "id"),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "metric_name", metricName),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "price", fmt.Sprint(price)),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "type", typeS),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "kind", kind),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "project_alert_limit", fmt.Sprint(projectLimit)),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "global_alert_limit", fmt.Sprint(globalLimit)),
					resource.TestCheckResourceAttrSet("taikun_showback_rule.foo", "showback_credential_id"),
					resource.TestCheckResourceAttr("taikun_showback_rule.foo", "showback_credential_name", showbackCredentialName),
				),
			},
			{
				ResourceName:      "taikun_showback_rule.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckTaikunShowbackRuleExists(state *terraform.State) error {
	apiClient := testAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_showback_rule" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)

		response, _, err := apiClient.ShowbackClient.ShowbackRulesApi.ShowbackrulesList(context.TODO()).Id(id).Execute()
		if err != nil || response.GetTotalCount() != 1 {
			return fmt.Errorf("showback rule doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunShowbackRuleDestroy(state *terraform.State) error {
	apiClient := testAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_showback_rule" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			id, _ := atoi32(rs.Primary.ID)

			response, _, err := apiClient.ShowbackClient.ShowbackRulesApi.ShowbackrulesList(context.TODO()).Id(id).Execute()
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if response.GetTotalCount() != 0 {
				return resource.RetryableError(errors.New("showback rule still exists"))
			}
			return nil
		})
		if timedOut(retryErr) {
			return errors.New("showback rule still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
