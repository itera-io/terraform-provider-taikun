package taikun

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/projects"
)

const testAccResourceTaikunProjectConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
}

resource "taikun_project" "foo" {
  name = "%s"
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id

  auto_upgrade = %t
  monitoring = %t
  expiration_date = "%s"
}
`

func TestAccResourceTaikunProject(t *testing.T) {
	cloudCredentialName := randomTestName()
	projectName := randomTestName()
	enableAutoUpgrade := true
	enableMonitoring := false
	expirationDate := "01/04/2999"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfig,
					cloudCredentialName,
					projectName,
					enableAutoUpgrade,
					enableMonitoring,
					expirationDate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "auto_upgrade", fmt.Sprint(enableAutoUpgrade)),
					resource.TestCheckResourceAttr("taikun_project.foo", "monitoring", fmt.Sprint(enableMonitoring)),
					resource.TestCheckResourceAttr("taikun_project.foo", "expiration_date", expirationDate),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
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
	cloudCredentialName := randomTestName()
	projectName := randomTestName()
	enableAutoUpgrade := true
	enableMonitoring := false
	expirationDate := "01/04/2999"
	newExpirationDate := "07/02/3000"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfig,
					cloudCredentialName,
					projectName,
					enableAutoUpgrade,
					enableMonitoring,
					expirationDate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "auto_upgrade", fmt.Sprint(enableAutoUpgrade)),
					resource.TestCheckResourceAttr("taikun_project.foo", "monitoring", fmt.Sprint(enableMonitoring)),
					resource.TestCheckResourceAttr("taikun_project.foo", "expiration_date", expirationDate),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfig,
					cloudCredentialName,
					projectName,
					enableAutoUpgrade,
					enableMonitoring,
					newExpirationDate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "auto_upgrade", fmt.Sprint(enableAutoUpgrade)),
					resource.TestCheckResourceAttr("taikun_project.foo", "monitoring", fmt.Sprint(enableMonitoring)),
					resource.TestCheckResourceAttr("taikun_project.foo", "expiration_date", newExpirationDate),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
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

  auto_upgrade = %t
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

  auto_upgrade = %t
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

  auto_upgrade = %t
  monitoring = %t
  expiration_date = "%s"
}
`

func TestAccResourceTaikunProjectModifyAlertingProfile(t *testing.T) {
	cloudCredentialName := randomTestName()
	projectName := randomTestName()
	alertingProfileName := randomTestName()
	newAlertingProfileName := randomTestName()
	alertingProfileResourceName := randomTestName()
	newAlertingProfileResourceName := randomTestName()
	enableAutoUpgrade := true
	enableMonitoring := false
	expirationDate := "01/04/2999"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
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
					enableAutoUpgrade,
					enableMonitoring,
					expirationDate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "alerting_profile_name", alertingProfileName),
					resource.TestCheckResourceAttr("taikun_project.foo", "auto_upgrade", fmt.Sprint(enableAutoUpgrade)),
					resource.TestCheckResourceAttr("taikun_project.foo", "monitoring", fmt.Sprint(enableMonitoring)),
					resource.TestCheckResourceAttr("taikun_project.foo", "expiration_date", expirationDate),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "alerting_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
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
					enableAutoUpgrade,
					enableMonitoring,
					expirationDate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "alerting_profile_name", newAlertingProfileName),
					resource.TestCheckResourceAttr("taikun_project.foo", "auto_upgrade", fmt.Sprint(enableAutoUpgrade)),
					resource.TestCheckResourceAttr("taikun_project.foo", "monitoring", fmt.Sprint(enableMonitoring)),
					resource.TestCheckResourceAttr("taikun_project.foo", "expiration_date", expirationDate),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
				),
			},
		},
	})
}

func TestAccResourceTaikunProjectDetachAlertingProfile(t *testing.T) {
	cloudCredentialName := randomTestName()
	projectName := randomTestName()
	alertingProfileName := randomTestName()
	enableAutoUpgrade := true
	enableMonitoring := false
	expirationDate := "01/04/2999"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfigWithAlertingProfile,
					cloudCredentialName,
					alertingProfileName,
					projectName,
					enableAutoUpgrade,
					enableMonitoring,
					expirationDate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "auto_upgrade", fmt.Sprint(enableAutoUpgrade)),
					resource.TestCheckResourceAttr("taikun_project.foo", "monitoring", fmt.Sprint(enableMonitoring)),
					resource.TestCheckResourceAttr("taikun_project.foo", "expiration_date", expirationDate),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "alerting_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfigWithAlertingProfileDetach,
					cloudCredentialName,
					alertingProfileName,
					projectName,
					enableAutoUpgrade,
					enableMonitoring,
					expirationDate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "alerting_profile_id", ""),
					resource.TestCheckResourceAttr("taikun_project.foo", "auto_upgrade", fmt.Sprint(enableAutoUpgrade)),
					resource.TestCheckResourceAttr("taikun_project.foo", "monitoring", fmt.Sprint(enableMonitoring)),
					resource.TestCheckResourceAttr("taikun_project.foo", "expiration_date", expirationDate),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
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
  auto_upgrade = false
  kubernetes_version = "%s"
}
`

func TestAccResourceTaikunProjectKubernetesVersion(t *testing.T) {
	cloudCredentialName := randomTestName()
	projectName := shortRandomTestName()
	kubernetesVersion := "v1.24.4"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
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
	cloudCredentialName := randomTestName()
	projectName := shortRandomTestName()
	locked := true
	unlocked := false

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
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
					resource.TestCheckResourceAttrSet("taikun_project.foo", "auto_upgrade"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "monitoring"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
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
					resource.TestCheckResourceAttrSet("taikun_project.foo", "auto_upgrade"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "monitoring"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
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
					resource.TestCheckResourceAttrSet("taikun_project.foo", "auto_upgrade"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "monitoring"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
				),
			},
		},
	})
}

func testAccCheckTaikunProjectExists(state *terraform.State) error {
	apiClient := testAccProvider.Meta().(*taikungoclient.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_project" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := projects.NewProjectsListParams().WithV(ApiVersion).WithID(&id)

		response, err := apiClient.Client.Projects.ProjectsList(params, apiClient)
		if err != nil || len(response.Payload.Data) != 1 {
			return fmt.Errorf("project doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunProjectDestroy(state *terraform.State) error {
	apiClient := testAccProvider.Meta().(*taikungoclient.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_project" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			id, _ := atoi32(rs.Primary.ID)
			params := projects.NewProjectsListParams().WithV(ApiVersion).WithID(&id)

			response, err := apiClient.Client.Projects.ProjectsList(params, apiClient)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if len(response.Payload.Data) != 0 {
				return resource.RetryableError(errors.New("project still exists"))
			}
			return nil
		})
		if timedOut(retryErr) {
			return errors.New("project still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
