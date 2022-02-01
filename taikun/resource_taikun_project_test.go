package taikun

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient/client/projects"
	"github.com/itera-io/taikungoclient/client/servers"
	"github.com/itera-io/taikungoclient/models"
)

func init() {
	resource.AddTestSweepers("taikun_project", &resource.Sweeper{
		Name:         "taikun_project",
		Dependencies: []string{},
		F: func(r string) error {

			meta, err := sharedConfig()
			if err != nil {
				return err
			}
			apiClient := meta.(*apiClient)

			params := projects.NewProjectsListParams().WithV(ApiVersion)

			var projectList []*models.ProjectListForUIDto

			for {
				response, err := apiClient.client.Projects.ProjectsList(params, apiClient)
				if err != nil {
					return err
				}
				projectList = append(projectList, response.GetPayload().Data...)
				if len(projectList) == int(response.GetPayload().TotalCount) {
					break
				}
				offset := int32(len(projectList))
				params = params.WithOffset(&offset)
			}

			var result *multierror.Error

			for _, e := range projectList {
				if shouldSweep(e.Name) {

					readParams := servers.NewServersDetailsParams().WithV(ApiVersion).WithProjectID(e.ID)
					response, err := apiClient.client.Servers.ServersDetails(readParams, apiClient)
					if err != nil {
						result = multierror.Append(result, err)
						continue
					}

					if response.Payload.Project.IsLocked {
						unlockedMode := getLockMode(false)
						unlockParams := projects.NewProjectsLockManagerParams().WithV(ApiVersion).WithID(&e.ID).WithMode(&unlockedMode)
						if _, err := apiClient.client.Projects.ProjectsLockManager(unlockParams, apiClient); err != nil {
							result = multierror.Append(result, err)
							continue
						}
					}

					serverIds := make([]int32, 0)
					for _, e := range response.Payload.Data {
						serverIds = append(serverIds, e.ID)
					}
					if len(serverIds) != 0 {
						deleteServerBody := &models.DeleteServerCommand{
							ProjectID: e.ID,
							ServerIds: serverIds,
						}
						deleteServerParams := servers.NewServersDeleteParams().WithV(ApiVersion).WithBody(deleteServerBody)
						_, _, err := apiClient.client.Servers.ServersDelete(deleteServerParams, apiClient)
						if err != nil {
							result = multierror.Append(result, err)
							continue
						}

						if err := resourceTaikunProjectWaitForStatus(context.Background(), []string{"Ready"}, []string{"PendingPurge", "Purging"}, apiClient, e.ID); err != nil {
							result = multierror.Append(result, err)
							continue
						}
					}

					params := projects.NewProjectsDeleteParams().WithV(ApiVersion).WithBody(&models.DeleteProjectCommand{ProjectID: e.ID})
					_, _, err = apiClient.client.Projects.ProjectsDelete(params, apiClient)
					if err != nil {
						result = multierror.Append(result, err)
					}
				}
			}

			return result.ErrorOrNil()
		},
	})
}

const testAccResourceTaikunProjectConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
  availability_zone = "%s"
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
					os.Getenv("AWS_AVAILABILITY_ZONE"),
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
					os.Getenv("AWS_AVAILABILITY_ZONE"),
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
					os.Getenv("AWS_AVAILABILITY_ZONE"),
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
  availability_zone = "%s"
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
  availability_zone = "%s"
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
  availability_zone = "%s"
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
					os.Getenv("AWS_AVAILABILITY_ZONE"),
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
					os.Getenv("AWS_AVAILABILITY_ZONE"),
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
					os.Getenv("AWS_AVAILABILITY_ZONE"),
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
					os.Getenv("AWS_AVAILABILITY_ZONE"),
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

const testAccResourceTaikunProjectQuotaConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
  availability_zone = "%s"
}
resource "taikun_project" "foo" {
  name = "%s"
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id
  %s
}
`

func TestAccResourceTaikunProjectQuota(t *testing.T) {
	cloudCredentialName := randomTestName()
	projectName := shortRandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectQuotaConfig,
					cloudCredentialName,
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					projectName,
					`
					quota_cpu_units = 500
					quota_disk_size = 200
					`,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttrSet("taikun_project.foo", "auto_upgrade"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "monitoring"),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "quota_id"),
					resource.TestCheckResourceAttr("taikun_project.foo", "quota_cpu_units", "500"),
					resource.TestCheckResourceAttr("taikun_project.foo", "quota_disk_size", "200"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectQuotaConfig,
					cloudCredentialName,
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					projectName,
					`
					quota_cpu_units = 501
					quota_ram_size = 200
					`,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttrSet("taikun_project.foo", "auto_upgrade"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "monitoring"),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "quota_id"),
					resource.TestCheckResourceAttr("taikun_project.foo", "quota_cpu_units", "501"),
					resource.TestCheckResourceAttr("taikun_project.foo", "quota_ram_size", "200"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectQuotaConfig,
					cloudCredentialName,
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					projectName,
					"",
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttrSet("taikun_project.foo", "auto_upgrade"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "monitoring"),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "quota_id"),
					resource.TestCheckResourceAttr("taikun_project.foo", "quota_cpu_units", "0"),
					resource.TestCheckResourceAttr("taikun_project.foo", "quota_ram_size", "0"),
					resource.TestCheckResourceAttr("taikun_project.foo", "quota_disk_size", "0"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectQuotaConfig,
					cloudCredentialName,
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					projectName,
					`
					quota_cpu_units = 502
					quota_disk_size = 201
					quota_ram_size = 201
					`,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttrSet("taikun_project.foo", "auto_upgrade"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "monitoring"),
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "quota_id"),
					resource.TestCheckResourceAttr("taikun_project.foo", "quota_cpu_units", "502"),
					resource.TestCheckResourceAttr("taikun_project.foo", "quota_ram_size", "201"),
					resource.TestCheckResourceAttr("taikun_project.foo", "quota_disk_size", "201"),
				),
			},
		},
	})
}

const testAccResourceTaikunProjectLockConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
  availability_zone = "%s"
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
					os.Getenv("AWS_AVAILABILITY_ZONE"),
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
					os.Getenv("AWS_AVAILABILITY_ZONE"),
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
					os.Getenv("AWS_AVAILABILITY_ZONE"),
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
	apiClient := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_project" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := projects.NewProjectsListParams().WithV(ApiVersion).WithID(&id)

		response, err := apiClient.client.Projects.ProjectsList(params, apiClient)
		if err != nil || len(response.Payload.Data) != 1 {
			return fmt.Errorf("project doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunProjectDestroy(state *terraform.State) error {
	apiClient := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_project" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			id, _ := atoi32(rs.Primary.ID)
			params := projects.NewProjectsListParams().WithV(ApiVersion).WithID(&id)

			response, err := apiClient.client.Projects.ProjectsList(params, apiClient)
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
