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

const testAccResourceTaikunProjectConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
}

resource "taikun_project" "foo" {
  name = "%s"
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id

  monitoring = %t
  expiration_date = "%s"
}
`

func TestAccResourceTaikunProject(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	projectName := utils.RandomTestName()
	enableMonitoring := false
	expirationDate := "01/04/2999"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAWS(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfig,
					cloudCredentialName,
					projectName,
					enableMonitoring,
					expirationDate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "monitoring", fmt.Sprint(enableMonitoring)),
					resource.TestCheckResourceAttr("taikun_project.foo", "expiration_date", expirationDate),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
				),
			},
			{
				ResourceName:      "taikun_project.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceTaikunProjectExtendLifetime(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	projectName := utils.RandomTestName()
	enableMonitoring := false
	expirationDate := "01/04/2999"
	newExpirationDate := "07/02/3000"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAWS(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfig,
					cloudCredentialName,
					projectName,
					enableMonitoring,
					expirationDate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "monitoring", fmt.Sprint(enableMonitoring)),
					resource.TestCheckResourceAttr("taikun_project.foo", "expiration_date", expirationDate),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfig,
					cloudCredentialName,
					projectName,
					enableMonitoring,
					newExpirationDate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "monitoring", fmt.Sprint(enableMonitoring)),
					resource.TestCheckResourceAttr("taikun_project.foo", "expiration_date", newExpirationDate),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
				),
			},
		},
	})
}

const testAccResourceTaikunProjectConfigWithAlertingProfile = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
}

resource "taikun_alerting_profile" "foo" {
  name = "%s"
  reminder = "Daily"
}

resource "taikun_project" "foo" {
  name = "%s"
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id

  alerting_profile_id = resource.taikun_alerting_profile.foo.id

  monitoring = %t
  expiration_date = "%s"
}
`

const testAccResourceTaikunProjectConfigWithAlertingProfileDetach = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
}

resource "taikun_alerting_profile" "foo" {
  name = "%s"
  reminder = "Daily"
}

resource "taikun_project" "foo" {
  name = "%s"
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id

  monitoring = %t
  expiration_date = "%s"
}
`

const testAccResourceTaikunProjectConfigWithAlertingProfiles = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
}

resource "taikun_alerting_profile" "%s" {
  name = "%s"
  reminder = "Daily"
}

resource "taikun_alerting_profile" "%s" {
  name = "%s"
  reminder = "Daily"
}

resource "taikun_project" "foo" {
  name = "%s"
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id

  alerting_profile_id = resource.taikun_alerting_profile.%s.id

  monitoring = %t
  expiration_date = "%s"
}
`

func TestAccResourceTaikunProjectModifyAlertingProfile(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	projectName := utils.RandomTestName()
	alertingProfileName := utils.RandomTestName()
	newAlertingProfileName := utils.RandomTestName()
	alertingProfileResourceName := utils.RandomTestName()
	newAlertingProfileResourceName := utils.RandomTestName()
	enableMonitoring := false
	expirationDate := "01/04/2999"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAWS(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfigWithAlertingProfiles,
					cloudCredentialName,
					alertingProfileResourceName,
					alertingProfileName,
					newAlertingProfileResourceName,
					newAlertingProfileName,
					projectName,
					alertingProfileResourceName,
					enableMonitoring,
					expirationDate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "alerting_profile_name", alertingProfileName),
					resource.TestCheckResourceAttr("taikun_project.foo", "monitoring", fmt.Sprint(enableMonitoring)),
					resource.TestCheckResourceAttr("taikun_project.foo", "expiration_date", expirationDate),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "alerting_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfigWithAlertingProfiles,
					cloudCredentialName,
					alertingProfileResourceName,
					alertingProfileName,
					newAlertingProfileResourceName,
					newAlertingProfileName,
					projectName,
					newAlertingProfileResourceName,
					enableMonitoring,
					expirationDate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "alerting_profile_name", newAlertingProfileName),
					resource.TestCheckResourceAttr("taikun_project.foo", "monitoring", fmt.Sprint(enableMonitoring)),
					resource.TestCheckResourceAttr("taikun_project.foo", "expiration_date", expirationDate),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
				),
			},
		},
	})
}

func TestAccResourceTaikunProjectDetachAlertingProfile(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	projectName := utils.RandomTestName()
	alertingProfileName := utils.RandomTestName()
	enableMonitoring := false
	expirationDate := "01/04/2999"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAWS(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfigWithAlertingProfile,
					cloudCredentialName,
					alertingProfileName,
					projectName,
					enableMonitoring,
					expirationDate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "monitoring", fmt.Sprint(enableMonitoring)),
					resource.TestCheckResourceAttr("taikun_project.foo", "expiration_date", expirationDate),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "alerting_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfigWithAlertingProfileDetach,
					cloudCredentialName,
					alertingProfileName,
					projectName,
					enableMonitoring,
					expirationDate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "alerting_profile_id", ""),
					resource.TestCheckResourceAttr("taikun_project.foo", "monitoring", fmt.Sprint(enableMonitoring)),
					resource.TestCheckResourceAttr("taikun_project.foo", "expiration_date", expirationDate),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
				),
			},
		},
	})
}

const testAccResourceTaikunProjectKubernetesVersionConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
}
resource "taikun_project" "foo" {
  name = "%s"
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id
  kubernetes_version = "%s"
}
`

func TestAccResourceTaikunProjectKubernetesVersion(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	projectName := utils.ShortRandomTestName()
	kubernetesVersion := "v1.30.4"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAWS(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectKubernetesVersionConfig,
					cloudCredentialName,
					projectName,
					kubernetesVersion,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "kubernetes_version", kubernetesVersion),
				),
			},
		},
	})
}

const testAccResourceTaikunProjectLockConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
}

resource "taikun_project" "foo" {
  name = "%s"
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id
  lock = %t
}
`

func TestAccResourceTaikunProjectToggleLock(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	projectName := utils.ShortRandomTestName()
	locked := true
	unlocked := false

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAWS(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectLockConfig,
					cloudCredentialName,
					projectName,
					locked,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttr("taikun_project.foo", "lock", fmt.Sprint(locked)),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "monitoring"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectLockConfig,
					cloudCredentialName,
					projectName,
					unlocked,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttr("taikun_project.foo", "lock", fmt.Sprint(unlocked)),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "monitoring"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectLockConfig,
					cloudCredentialName,
					projectName,
					locked,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttr("taikun_project.foo", "lock", fmt.Sprint(locked)),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "monitoring"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
				),
			},
		},
	})
}

func testAccCheckTaikunProjectExists(state *terraform.State) error {
	apiClient := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_project" {
			continue
		}

		id, _ := utils.Atoi32(rs.Primary.ID)

		response, _, err := apiClient.Client.ProjectsAPI.ProjectsList(context.TODO()).Id(id).Execute()
		if err != nil || len(response.GetData()) != 1 {
			return fmt.Errorf("project doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunProjectDestroy(state *terraform.State) error {
	apiClient := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_project" {
			continue
		}

		retryErr := retry.RetryContext(context.Background(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
			id, _ := utils.Atoi32(rs.Primary.ID)

			response, _, err := apiClient.Client.ProjectsAPI.ProjectsList(context.TODO()).Id(id).Execute()
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if len(response.GetData()) != 0 {
				return retry.RetryableError(errors.New("project still exists"))
			}
			return nil
		})
		if utils.TimedOut(retryErr) {
			return errors.New("project still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}

const testAccResourceTaikunProjectAutoscalerOpenstackConfig = `
resource "taikun_cloud_credential_openstack" "foo" {
  name = "%s"
}

data "taikun_flavors" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_openstack.foo.id
  min_cpu = 2
  max_cpu = 2
  min_ram = 4
  max_ram = 8
}

locals {
  flavors = [for flavor in data.taikun_flavors.foo.flavors: flavor.name]
}

resource "taikun_project" "foo" {
  name = "%s"
  flavors = local.flavors
  cloud_credential_id = resource.taikun_cloud_credential_openstack.foo.id
  
  autoscaler_name = "scaler"
  autoscaler_flavor = local.flavors[0]
  autoscaler_min_size = 1
  autoscaler_max_size = 2
  autoscaler_disk_size = 30
}
`

func TestAccResourceTaikunAutoscalerOpenstackProject(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	projectName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAWS(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectAutoscalerOpenstackConfig,
					cloudCredentialName,
					projectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttr("taikun_project.foo", "autoscaler_name", "scaler"),
					resource.TestCheckResourceAttr("taikun_project.foo", "autoscaler_min_size", "1"),
					resource.TestCheckResourceAttr("taikun_project.foo", "autoscaler_max_size", "2"),
					resource.TestCheckResourceAttr("taikun_project.foo", "autoscaler_disk_size", "30"),
				),
			},
			{
				ResourceName:      "taikun_project.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccResourceTaikunProjectAutoscalerAwsConfig = `
resource "taikun_cloud_credential_aws" "foo" {
 name = "%s"
}

data "taikun_flavors" "foo" {
 cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id
 min_cpu = 2
 max_cpu = 2
 min_ram = 4
 max_ram = 8
}

locals {
 flavors = [for flavor in data.taikun_flavors.foo.flavors: flavor.name]
}

resource "taikun_project" "foo" {
 name = "%s"
 flavors = local.flavors
 cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id

 autoscaler_name = "scaler"
 autoscaler_flavor = local.flavors[0]
 autoscaler_min_size = "%d"
 autoscaler_max_size = "%d"
 autoscaler_disk_size = "%d"
 autoscaler_spot_enabled = "%t"

 spot_max_price="%d"
 spot_vms="%t"
 spot_worker="%t"
}
`

func TestAccResourceTaikunAutoscalerAwsProject(t *testing.T) {
	cloudCredentialName := utils.RandomTestName()
	projectName := utils.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t); utils_testing.TestAccPreCheckAWS(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectAutoscalerAwsConfig,
					cloudCredentialName,
					projectName,
					2,     // Autoscaler min
					3,     // Autoscaler max
					31,    // Autoscaler disk
					true,  // Autoscaler spot enabled
					10,    // Autoscaler spot max price
					true,  // Spot Vms enabled
					true), //Spot workers enabled
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttr("taikun_project.foo", "autoscaler_name", "scaler"),
					resource.TestCheckResourceAttr("taikun_project.foo", "autoscaler_min_size", "2"),
					resource.TestCheckResourceAttr("taikun_project.foo", "autoscaler_max_size", "3"),
					resource.TestCheckResourceAttr("taikun_project.foo", "autoscaler_disk_size", "31"),
					resource.TestCheckResourceAttr("taikun_project.foo", "autoscaler_spot_enabled", "true"),
					resource.TestCheckResourceAttr("taikun_project.foo", "spot_full", "false"),
					resource.TestCheckResourceAttr("taikun_project.foo", "spot_worker", "true"),
					resource.TestCheckResourceAttr("taikun_project.foo", "spot_vms", "true"),
					resource.TestCheckResourceAttr("taikun_project.foo", "spot_max_price", "10"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectAutoscalerAwsConfig,
					cloudCredentialName,
					projectName,
					1,     // Autoscaler min
					2,     // Autoscaler max
					31,    // Autoscaler disk
					false, // Autoscaler spot enabled
					10,    // Autoscaler spot max price
					false, // Spot Vms enabled
					true), //Spot workers enabled
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttr("taikun_project.foo", "autoscaler_name", "scaler"),
					resource.TestCheckResourceAttr("taikun_project.foo", "autoscaler_min_size", "1"),
					resource.TestCheckResourceAttr("taikun_project.foo", "autoscaler_max_size", "2"),
					resource.TestCheckResourceAttr("taikun_project.foo", "autoscaler_disk_size", "31"),
					resource.TestCheckResourceAttr("taikun_project.foo", "autoscaler_spot_enabled", "false"),
					resource.TestCheckResourceAttr("taikun_project.foo", "spot_full", "false"),
					resource.TestCheckResourceAttr("taikun_project.foo", "spot_worker", "true"),
					resource.TestCheckResourceAttr("taikun_project.foo", "spot_vms", "false"),
					resource.TestCheckResourceAttr("taikun_project.foo", "spot_max_price", "10"),
				),
			},
			{
				ResourceName:      "taikun_project.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
