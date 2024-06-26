package testing

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tk "github.com/itera-io/taikungoclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceTaikunKubernetesProfileConfig = `
resource "taikun_kubernetes_profile" "foo" {
	name = "%s"
    lock = %t
}
`

func TestAccResourceTaikunKubernetesProfile(t *testing.T) {
	firstName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
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
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "nvidia_gpu_operator"),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "nvidia_gpu_operator", "false"),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "proxmox_storage", "NFS"),
				),
			},
			{
				ResourceName: "taikun_kubernetes_profile.foo",
				ImportState:  true,
			},
		},
	})
}

const testAccResourceTaikunKubernetesProfileNoUniqueClusterNameConfig = `
resource "taikun_kubernetes_profile" "foo" {
	name = "%s"
	unique_cluster_name = false
}
`

func TestAccResourceTaikunKubernetesProfileNoUniqueClusterName(t *testing.T) {
	firstName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunKubernetesProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunKubernetesProfileNoUniqueClusterNameConfig, firstName),
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
				ResourceName: "taikun_kubernetes_profile.foo",
				ImportState:  true,
			},
		},
	})
}

const TestAccResourceTaikunKubernetesProfileNvidiaGpuEnableConfig = `
resource "taikun_kubernetes_profile" "foo" {
	name = "%s"
	nvidia_gpu_operator = true
}
`

func TestAccResourceTaikunKubernetesProfileNvidiaGpuEnable(t *testing.T) {
	firstName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunKubernetesProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(TestAccResourceTaikunKubernetesProfileNvidiaGpuEnableConfig, firstName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunKubernetesProfileExists,
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "load_balancing_solution", "Octavia"),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "schedule_on_master", "false"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "bastion_proxy"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "nvidia_gpu_operator"),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "nvidia_gpu_operator", "true"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "wasm"),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "wasm", "false"),
				),
			},
			{
				ResourceName: "taikun_kubernetes_profile.foo",
				ImportState:  true,
			},
		},
	})
}

func TestAccResourceTaikunKubernetesProfileLock(t *testing.T) {
	firstName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckPrometheus(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
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
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "nvidia_gpu_operator"),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "nvidia_gpu_operator", "false"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "unique_cluster_name"),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "unique_cluster_name", "false"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "wasm"),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "wasm", "false"),
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
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "nvidia_gpu_operator"),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "nvidia_gpu_operator", "false"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "unique_cluster_name"),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "unique_cluster_name", "false"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "wasm"),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "wasm", "false"),
				),
			},
		},
	})
}

func testAccCheckTaikunKubernetesProfileExists(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_kubernetes_profile" {
			continue
		}

		id, _ := utils.Atoi32(rs.Primary.ID)

		response, _, err := client.Client.KubernetesProfilesAPI.KubernetesprofilesList(context.TODO()).Id(id).Execute()
		if err != nil || response.GetTotalCount() != 1 {
			return fmt.Errorf("kubernetes profile doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunKubernetesProfileDestroy(state *terraform.State) error {
	client := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_kubernetes_profile" {
			continue
		}

		retryErr := retry.RetryContext(context.Background(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
			id, _ := utils.Atoi32(rs.Primary.ID)

			response, _, err := client.Client.KubernetesProfilesAPI.KubernetesprofilesList(context.TODO()).Id(id).Execute()
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if response.GetTotalCount() != 0 {
				return retry.RetryableError(errors.New("kubernetes profile still exists"))
			}
			return nil
		})
		if utils.TimedOut(retryErr) {
			return errors.New("kubernetes profile still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
