package taikun

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/itera-io/taikungoclient"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
	firstName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
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
	firstName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
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
	firstName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
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
	client := testAccProvider.Meta().(*taikungoclient.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_policy_profile" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		resource, err := resourceTaikunPolicyProfileFind(id, client)
		if err != nil || resource == nil {
			return fmt.Errorf("policy profile doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunPolicyProfileDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*taikungoclient.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_policy_profile" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			id, _ := atoi32(rs.Primary.ID)
			policyProfile, err := resourceTaikunPolicyProfileFind(id, client)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if policyProfile != nil {
				return resource.RetryableError(errors.New("policy profile still exists"))
			}
			return nil
		})
		if timedOut(retryErr) {
			return errors.New("policy profile still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
