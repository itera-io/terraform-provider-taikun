package testing

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tk "github.com/itera-io/taikungoclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/policy_profile"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceTaikunPolicyProfileConfig = `
resource "taikun_policy_profile" "foo" {
  name = "%s"
  lock = %t

  forbid_node_port = %t
  forbid_http_ingress = %t
  require_probe = %t
  unique_ingress = %t
  unique_service_selector = %t

  %s
}
`

func TestAccResourceTaikunPolicyProfile(t *testing.T) {
	firstName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunPolicyProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccResourceTaikunPolicyProfileConfig,
					firstName,
					false,
					false,
					false,
					false,
					false,
					false,
					"",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunPolicyProfileExists,
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "forbid_node_port", "false"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "forbid_http_ingress", "false"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "require_probe", "false"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "unique_ingress", "false"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "unique_service_selector", "false"),
					resource.TestCheckResourceAttrSet("taikun_policy_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_policy_profile.foo", "organization_name"),
				),
			},
			{
				ResourceName:      "taikun_policy_profile.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceTaikunPolicyProfileLock(t *testing.T) {
	firstName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunPolicyProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccResourceTaikunPolicyProfileConfig,
					firstName,
					false,
					true,
					true,
					true,
					true,
					true,
					"",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunPolicyProfileExists,
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "forbid_node_port", "true"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "forbid_http_ingress", "true"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "require_probe", "true"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "unique_ingress", "true"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "unique_service_selector", "true"),
					resource.TestCheckResourceAttrSet("taikun_policy_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_policy_profile.foo", "organization_name"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunPolicyProfileConfig,
					firstName,
					true,
					true,
					true,
					true,
					true,
					true,
					"",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunPolicyProfileExists,
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "lock", "true"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "forbid_node_port", "true"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "forbid_http_ingress", "true"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "require_probe", "true"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "unique_ingress", "true"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "unique_service_selector", "true"),
					resource.TestCheckResourceAttrSet("taikun_policy_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_policy_profile.foo", "organization_name"),
				),
			},
		},
	})
}

func TestAccResourceTaikunPolicyProfileUpdate(t *testing.T) {
	firstName := utils.RandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunPolicyProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccResourceTaikunPolicyProfileConfig,
					firstName,
					false,
					true,
					true,
					true,
					true,
					true,
					"forbidden_tags = [\"tag1\", \"tag2\"]",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunPolicyProfileExists,
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "forbid_node_port", "true"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "forbid_http_ingress", "true"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "require_probe", "true"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "unique_ingress", "true"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "unique_service_selector", "true"),
					resource.TestCheckResourceAttrSet("taikun_policy_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_policy_profile.foo", "organization_name"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "forbidden_tags.#", "2"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "forbidden_tags.0", "tag1"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "forbidden_tags.1", "tag2"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunPolicyProfileConfig,
					firstName,
					false,
					true,
					false,
					true,
					false,
					true,
					"forbidden_tags = [\"tag3\"]",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunPolicyProfileExists,
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "forbid_node_port", "true"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "forbid_http_ingress", "false"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "require_probe", "true"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "unique_ingress", "false"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "unique_service_selector", "true"),
					resource.TestCheckResourceAttrSet("taikun_policy_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_policy_profile.foo", "organization_name"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "forbidden_tags.#", "1"),
					resource.TestCheckResourceAttr("taikun_policy_profile.foo", "forbidden_tags.0", "tag3"),
				),
			},
		},
	})
}

func testAccCheckTaikunPolicyProfileExists(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_policy_profile" {
			continue
		}

		id, _ := utils.Atoi32(rs.Primary.ID)
		resource, err := policy_profile.ResourceTaikunPolicyProfileFind(id, client)
		if err != nil || resource == nil {
			return fmt.Errorf("policy profile doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunPolicyProfileDestroy(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_policy_profile" {
			continue
		}

		retryErr := retry.RetryContext(context.Background(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
			id, _ := utils.Atoi32(rs.Primary.ID)
			policyProfile, err := policy_profile.ResourceTaikunPolicyProfileFind(id, client)
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if policyProfile != nil {
				return retry.RetryableError(errors.New("policy profile still exists"))
			}
			return nil
		})
		if utils.TimedOut(retryErr) {
			return errors.New("policy profile still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
