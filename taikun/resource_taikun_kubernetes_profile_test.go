package taikun

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient/client/kubernetes_profiles"
)

const testAccResourceTaikunKubernetesProfileConfig = `
resource "taikun_kubernetes_profile" "foo" {
	name = "%s"
    lock = %t
}
`

func TestAccResourceTaikunKubernetesProfile(t *testing.T) {
	firstName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunKubernetesProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunKubernetesProfileConfig, firstName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunKubernetesProfileExists,
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "load_balancing_solution", "Octavia"),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "schedule_on_master", "false"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "bastion_proxy"),
				),
			},
			{
				ResourceName:      "taikun_kubernetes_profile.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceTaikunKubernetesProfileLock(t *testing.T) {
	firstName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunKubernetesProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunKubernetesProfileConfig, firstName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunKubernetesProfileExists,
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "load_balancing_solution", "Octavia"),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "schedule_on_master", "false"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "bastion_proxy"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunKubernetesProfileConfig, firstName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunKubernetesProfileExists,
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "lock", "true"),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "load_balancing_solution", "Octavia"),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "schedule_on_master", "false"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "bastion_proxy"),
				),
			},
		},
	})
}

func testAccCheckTaikunKubernetesProfileExists(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_kubernetes_profile" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := kubernetes_profiles.NewKubernetesProfilesListParams().WithV(ApiVersion).WithID(&id)

		response, err := client.client.KubernetesProfiles.KubernetesProfilesList(params, client)
		if err != nil || response.Payload.TotalCount != 1 {
			return fmt.Errorf("kubernetes profile doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunKubernetesProfileDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_kubernetes_profile" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			id, _ := atoi32(rs.Primary.ID)
			params := kubernetes_profiles.NewKubernetesProfilesListParams().WithV(ApiVersion).WithID(&id)

			response, err := client.client.KubernetesProfiles.KubernetesProfilesList(params, client)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if response.Payload.TotalCount != 0 {
				return resource.RetryableError(errors.New("kubernetes profile still exists"))
			}
			return nil
		})
		if timedOut(retryErr) {
			return errors.New("kubernetes profile still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
