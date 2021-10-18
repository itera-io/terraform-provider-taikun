package taikun

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient/client/kubernetes_profiles"
	"github.com/itera-io/taikungoclient/models"
	"strings"
	"testing"
)

func init() {
	resource.AddTestSweepers("taikun_kubernetes_profile", &resource.Sweeper{
		Name: "taikun_kubernetes_profile",
		F: func(r string) error {

			meta, err := sharedConfig()
			if err != nil {
				return err
			}
			apiClient := meta.(*apiClient)

			params := kubernetes_profiles.NewKubernetesProfilesListParams().WithV(ApiVersion)

			var kubernetesProfilesList []*models.KubernetesProfilesListDto
			for {
				response, err := apiClient.client.KubernetesProfiles.KubernetesProfilesList(params, apiClient)
				if err != nil {
					return err
				}
				kubernetesProfilesList = append(kubernetesProfilesList, response.GetPayload().Data...)
				if len(kubernetesProfilesList) == int(response.GetPayload().TotalCount) {
					break
				}
				offset := int32(len(kubernetesProfilesList))
				params = params.WithOffset(&offset)
			}

			for _, e := range kubernetesProfilesList {
				if strings.HasPrefix(e.Name, testNamePrefix) {
					params := kubernetes_profiles.NewKubernetesProfilesDeleteParams().WithV(ApiVersion).WithID(e.ID)
					_, _, err = apiClient.client.KubernetesProfiles.KubernetesProfilesDelete(params, apiClient)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	})
}

const testAccResourceTaikunKubernetesProfileConfig = `
resource "taikun_kubernetes_profile" "foo" {
	name = "%s"
    is_locked = %t
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunKubernetesProfileExists,
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "is_locked", "false"),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "load_balancing_solution", "Octavia"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "bastion_proxy_enabled"),
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunKubernetesProfileExists,
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "is_locked", "false"),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "load_balancing_solution", "Octavia"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "bastion_proxy_enabled"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunKubernetesProfileConfig, firstName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunKubernetesProfileExists,
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "name", firstName),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "is_locked", "true"),
					resource.TestCheckResourceAttr("taikun_kubernetes_profile.foo", "load_balancing_solution", "Octavia"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "organization_name"),
					resource.TestCheckResourceAttrSet("taikun_kubernetes_profile.foo", "bastion_proxy_enabled"),
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

		id, _ := atoi32(rs.Primary.ID)
		params := kubernetes_profiles.NewKubernetesProfilesListParams().WithV(ApiVersion).WithID(&id)

		response, err := client.client.KubernetesProfiles.KubernetesProfilesList(params, client)
		if err == nil && response.Payload.TotalCount != 0 {
			return fmt.Errorf("kubernetes profile still exists (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}
