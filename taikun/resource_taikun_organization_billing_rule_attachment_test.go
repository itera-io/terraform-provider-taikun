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

const testAccResourceTaikunOrganizationBillingRuleAttachmentConfig = `
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

resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = %f
}

resource "taikun_organization_billing_rule_attachment" "foo" {
  billing_rule_id = resource.taikun_billing_rule.foo.id
  organization_id = resource.taikun_organization.foo.id
  discount_rate   = %f
}
`

func TestAccResourceTaikunOrganizationBillingRuleAttachment(t *testing.T) {
	credName := randomTestName()
	orgName := randomTestName()
	ruleName := randomTestName()
	fullOrgName := randomString()
	globalDiscountRate := math.Round(rand.Float64()*10000) / 100
	ruleDiscountRate := math.Round(rand.Float64()*10000) / 100

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunOrganizationBillingRuleAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunOrganizationBillingRuleAttachmentConfig,
					credName,
					os.Getenv("PROMETHEUS_PASSWORD"),
					os.Getenv("PROMETHEUS_URL"),
					os.Getenv("PROMETHEUS_USERNAME"),
					ruleName,
					orgName,
					fullOrgName,
					globalDiscountRate,
					ruleDiscountRate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunOrganizationBillingRuleAttachmentExists,
					resource.TestCheckResourceAttrSet("taikun_organization_billing_rule_attachment.foo", "billing_rule_id"),
					resource.TestCheckResourceAttrSet("taikun_organization_billing_rule_attachment.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_organization_billing_rule_attachment.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_organization_billing_rule_attachment.foo", "discount_rate"),
				),
			},
		},
	})
}

func testAccCheckTaikunOrganizationBillingRuleAttachmentExists(state *terraform.State) error {
	apiClient := testAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_organization_billing_rule_attachment" {
			continue
		}

		organizationId, billingRuleId, err := parseOrganizationBillingRuleAttachmentId(rs.Primary.ID)
		if err != nil {
			return err
		}

		response, _, err := apiClient.Client.PrometheusRulesApi.PrometheusrulesList(context.TODO()).Id(billingRuleId).Execute()
		if err != nil {
			return err
		}
		if len(response.GetData()) != 1 {
			return fmt.Errorf("billing rule with ID %d not found", billingRuleId)
		}

		rawBillingRule := response.GetData()[0]

		for _, e := range rawBillingRule.BoundOrganizations {
			if e.GetId() == organizationId {
				return nil
			}
		}

		return fmt.Errorf("organization_billing_rule_attachment doesn't exist (id = %s)", rs.Primary.ID)
	}

	return nil
}

func testAccCheckTaikunOrganizationBillingRuleAttachmentDestroy(state *terraform.State) error {
	apiClient := testAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_organization_billing_rule_attachment" {
			continue
		}

		organizationId, billingRuleId, err := parseOrganizationBillingRuleAttachmentId(rs.Primary.ID)
		if err != nil {
			return err
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			response, _, err := apiClient.Client.PrometheusRulesApi.PrometheusrulesList(context.TODO()).Id(billingRuleId).Execute()
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if len(response.GetData()) != 1 {
				return nil
			}

			rawBillingRule := response.GetData()[0]

			for _, e := range rawBillingRule.BoundOrganizations {
				if e.GetId() == organizationId {
					return resource.RetryableError(errors.New("organization_billing_rule_attachment still exists"))
				}
			}
			return nil
		})
		if timedOut(retryErr) {
			return errors.New("organization_billing_rule_attachment still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
