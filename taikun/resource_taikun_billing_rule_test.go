package taikun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient/client/prometheus"
	"github.com/itera-io/taikungoclient/models"
	"strings"
	"testing"
)

func init() {
	resource.AddTestSweepers("taikun_billing_rule", &resource.Sweeper{
		Name: "taikun_billing_rule",
		F: func(r string) error {

			meta, err := sharedConfig()
			if err != nil {
				return err
			}
			apiClient := meta.(*apiClient)

			params := prometheus.NewPrometheusListOfRulesParams().WithV(ApiVersion)

			var billingRulesList []*models.PrometheusRuleListDto
			for {
				response, err := apiClient.client.Prometheus.PrometheusListOfRules(params, apiClient)
				if err != nil {
					return err
				}
				billingRulesList = append(billingRulesList, response.GetPayload().Data...)
				if len(billingRulesList) == int(response.GetPayload().TotalCount) {
					break
				}
				offset := int32(len(billingRulesList))
				params = params.WithOffset(&offset)
			}

			for _, e := range billingRulesList {
				if strings.HasPrefix(e.Name, testNamePrefix) {
					params := prometheus.NewPrometheusDeleteParams().WithV(ApiVersion).WithID(e.ID)
					_, err = apiClient.client.Prometheus.PrometheusDelete(params, apiClient)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	})
}

const testAccResourceTaikunBillingRule = `
resource "taikun_billing_rule" "foo" {
  name            = "%s"
  metric_name     =   "coredns_forward_request_duration_seconds"
  price = 1
  type = "Sum"
  billing_credential_id = "89"
  label {
    label = "label"
    value = "value"
  }
}
`

func TestAccResourceTaikunBillingRule(t *testing.T) {
	firstName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunBillingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunBillingRule, firstName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunBillingRuleExists,
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "metric_name", "coredns_forward_request_duration_seconds")),
			},
		},
	})
}

//func TestAccResourceTaikunBillingRuleRename(t *testing.T) {
//	firstName := randomTestName()
//	secondName := randomTestName()
//
//	resource.ParallelTest(t, resource.TestCase{
//		PreCheck:          func() { testAccPreCheck(t) },
//		ProviderFactories: testAccProviderFactories,
//		CheckDestroy:      testAccCheckTaikunBillingRuleDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config: fmt.Sprintf(testAccResourceTaikunBillingRule, firstName, false),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckTaikunBillingRuleExists,
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "name", firstName),
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "is_locked", "false"),
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "organization_id", "638"),
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "dns_server.#", "2"),
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "dns_server.0.address", "8.8.8.8"),
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "dns_server.1.address", "8.8.4.4"),
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "ntp_server.#", "2"),
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "ntp_server.0.address", "time.windows.com"),
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "ntp_server.1.address", "ntp.pool.org"),
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "ssh_user.#", "1"),
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "ssh_user.0.name", "oui oui"),
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "ssh_user.0.public_key", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"),
//				),
//			},
//			{
//				Config: fmt.Sprintf(testAccResourceTaikunBillingRule, secondName, true),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckTaikunBillingRuleExists,
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "name", secondName),
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "is_locked", "true"),
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "organization_id", "638"),
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "dns_server.#", "2"),
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "dns_server.0.address", "8.8.8.8"),
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "dns_server.1.address", "8.8.4.4"),
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "ntp_server.#", "2"),
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "ntp_server.0.address", "time.windows.com"),
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "ntp_server.1.address", "ntp.pool.org"),
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "ssh_user.#", "1"),
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "ssh_user.0.name", "oui oui"),
//					resource.TestCheckResourceAttr("taikun_billing_rule.foo", "ssh_user.0.public_key", "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"),
//				),
//			},
//		},
//	})
//}

func testAccCheckTaikunBillingRuleExists(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_billing_rule" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := prometheus.NewPrometheusListOfRulesParams().WithV(ApiVersion).WithID(&id)

		response, err := client.client.Prometheus.PrometheusListOfRules(params, client)
		if err != nil || response.Payload.TotalCount != 1 {
			return fmt.Errorf("access profile doesn't exist")
		}
	}

	return nil
}

func testAccCheckTaikunBillingRuleDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_billing_rule" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := prometheus.NewPrometheusListOfRulesParams().WithV(ApiVersion).WithID(&id)

		response, err := client.client.Prometheus.PrometheusListOfRules(params, client)
		if err == nil && response.Payload.TotalCount != 0 {
			return fmt.Errorf("access profile still exists")
		}
	}

	return nil
}
