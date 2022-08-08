package taikun

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/prometheus"
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
	credName := randomTestName()
	ruleName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
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
	credName := randomTestName()
	ruleName := randomTestName()
	ruleNameNew := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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
	credName := randomTestName()
	ruleName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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
	client := testAccProvider.Meta().(*taikungoclient.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_billing_rule" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := prometheus.NewPrometheusListOfRulesParams().WithV(ApiVersion).WithID(&id)

		response, err := client.Client.Prometheus.PrometheusListOfRules(params, client)
		if err != nil || response.Payload.TotalCount != 1 {
			return fmt.Errorf("billing rule doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunBillingRuleDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*taikungoclient.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_billing_rule" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			id, _ := atoi32(rs.Primary.ID)
			params := prometheus.NewPrometheusListOfRulesParams().WithV(ApiVersion).WithID(&id)

			response, err := client.Client.Prometheus.PrometheusListOfRules(params, client)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if response.Payload.TotalCount != 0 {
				return resource.RetryableError(errors.New("billing rule still exists ()"))
			}
			return nil
		})
		if timedOut(retryErr) {
			return errors.New("billing rule still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
