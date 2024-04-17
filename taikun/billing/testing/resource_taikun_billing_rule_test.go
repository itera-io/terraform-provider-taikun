package testing

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tk "github.com/itera-io/taikungoclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceTaikunBillingRuleConfig = `
resource "taikun_billing_credential" "foo" {
  name            = "%s"
  lock       = false

  prometheus_password = "%s"
  prometheus_url = "%s"
  prometheus_username = "%s"
}

resource "taikun_billing_rule" "foo" {
  name            = "%s"
  metric_name     =   "coredns_forward_request_duration_seconds"
  price = 1
  type = "Sum"
  billing_credential_id = resource.taikun_billing_credential.foo.id
  label {
    key = "key"
    value = "value"
  }
}
`

func TestAccResourceTaikunBillingRule(t *testing.T) {
	credName := utils.RandomTestName()
	ruleName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunBillingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunBillingRuleConfig,
					credName,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
					ruleName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunBillingRuleExists,
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "name", ruleName),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "metric_name", "coredns_forward_request_duration_seconds"),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "type", "Sum"),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "price", "1"),
					resource.TestCheckResourceAttrSet("taikun_billing_rule.foo", "billing_credential_id"),
				),
			},
			{
				ResourceName:      "taikun_billing_rule.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceTaikunBillingRuleRename(t *testing.T) {
	credName := utils.RandomTestName()
	ruleName := utils.RandomTestName()
	ruleNameNew := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunBillingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunBillingRuleConfig,
					credName,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
					ruleName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunBillingRuleExists,
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "name", ruleName),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "metric_name", "coredns_forward_request_duration_seconds"),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "type", "Sum"),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "price", "1"),
					resource.TestCheckResourceAttrSet("taikun_billing_rule.foo", "billing_credential_id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunBillingRuleConfig,
					credName,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
					ruleNameNew,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunBillingRuleExists,
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "name", ruleNameNew),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "metric_name", "coredns_forward_request_duration_seconds"),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "type", "Sum"),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "price", "1"),
					resource.TestCheckResourceAttrSet("taikun_billing_rule.foo", "billing_credential_id"),
				),
			},
		},
	})
}

const testAccResourceTaikunBillingRuleConfigUpdateLabels = `
resource "taikun_billing_credential" "foo" {
  name            = "%s"
  lock       = false

  prometheus_password = "%s"
  prometheus_url = "%s"
  prometheus_username = "%s"
}

resource "taikun_billing_rule" "foo" {
  name            = "%s"
  metric_name     =   "coredns_forward_request_duration_seconds"
  price = 1
  type = "Sum"
  billing_credential_id = resource.taikun_billing_credential.foo.id
  label {
    key = "key1"
    value = "value1"
  }
  label {
    key = "key2"
    value = "value2"
  }
  label {
    key = "key3"
    value = "value3"
  }
  label {
    key = "key4"
    value = "value4"
  }
}
`

func TestAccResourceTaikunBillingRuleUpdateLabels(t *testing.T) {
	credName := utils.RandomTestName()
	ruleName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunBillingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunBillingRuleConfig,
					credName,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
					ruleName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunBillingRuleExists,
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "name", ruleName),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "metric_name", "coredns_forward_request_duration_seconds"),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "type", "Sum"),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "price", "1"),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "label.#", "1"),
					resource.TestCheckResourceAttrSet("taikun_billing_rule.foo", "billing_credential_id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunBillingRuleConfigUpdateLabels,
					credName,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
					ruleName,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunBillingRuleExists,
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "name", ruleName),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "metric_name", "coredns_forward_request_duration_seconds"),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "type", "Sum"),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "price", "1"),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "label.#", "4"),
					resource.TestCheckResourceAttrSet("taikun_billing_rule.foo", "billing_credential_id"),
				),
			},
		},
	})
}

func testAccCheckTaikunBillingRuleExists(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_billing_rule" {
			continue
		}

		id, _ := utils.Atoi32(rs.Primary.ID)

		response, _, err := client.Client.PrometheusRulesAPI.PrometheusrulesList(context.TODO()).Id(id).Execute()
		if err != nil || response.GetTotalCount() != 1 {
			return fmt.Errorf("billing rule doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunBillingRuleDestroy(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_billing_rule" {
			continue
		}

		retryErr := retry.RetryContext(context.Background(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
			id, _ := utils.Atoi32(rs.Primary.ID)

			response, _, err := client.Client.PrometheusRulesAPI.PrometheusrulesList(context.TODO()).Id(id).Execute()
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if response.GetTotalCount() != 0 {
				return retry.RetryableError(errors.New("billing rule still exists ()"))
			}
			return nil
		})
		if utils.TimedOut(retryErr) {
			return errors.New("billing rule still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
