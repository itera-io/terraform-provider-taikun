package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceTaikunProjectToggleMonitoring(t *testing.T) {
	cloudCredentialName := randomTestName()
	projectName := randomTestName()
	enableAutoUpgrade := true
	enableMonitoring := true
	disableMonitoring := false
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
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfig,
					cloudCredentialName,
					os.Getenv("AWS_AVAILABILITY_ZONE"),
					projectName,
					enableAutoUpgrade,
					disableMonitoring,
					expirationDate),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "auto_upgrade", fmt.Sprint(enableAutoUpgrade)),
					resource.TestCheckResourceAttr("taikun_project.foo", "monitoring", fmt.Sprint(disableMonitoring)),
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
		},
	})
}

const testAccResourceTaikunProjectToggleBackupConfig = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
}

resource "taikun_backup_credential" "foo" {
  name            = "%s"

  s3_endpoint = "%s"
  s3_region   = "%s"
}

resource "taikun_backup_credential" "foo2" {
  name            = "%s"

  s3_endpoint = "%s"
  s3_region   = "%s"
}

resource "taikun_project" "foo" {
  name = "%s"
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id

  %s
}
`

func TestAccResourceTaikunProjectToggleBackup(t *testing.T) {
	cloudCredentialName := randomTestName()
	backupCredentialName := randomTestName()
	backupCredentialName2 := randomTestName()
	projectName := shortRandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t); testAccPreCheckS3(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectToggleBackupConfig,
					cloudCredentialName,
					backupCredentialName,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
					backupCredentialName2,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
					projectName,
					"backup_credential_id = resource.taikun_backup_credential.foo.id",
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
					resource.TestCheckResourceAttrPair("taikun_project.foo", "backup_credential_id", "taikun_backup_credential.foo", "id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectToggleBackupConfig,
					cloudCredentialName,
					backupCredentialName,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
					backupCredentialName2,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
					projectName,
					"backup_credential_id = resource.taikun_backup_credential.foo2.id",
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
					resource.TestCheckResourceAttrPair("taikun_project.foo", "backup_credential_id", "taikun_backup_credential.foo2", "id"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectToggleBackupConfig,
					cloudCredentialName,
					backupCredentialName,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
					backupCredentialName2,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
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
					resource.TestCheckResourceAttr("taikun_project.foo", "backup_credential_id", ""),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectToggleBackupConfig,
					cloudCredentialName,
					backupCredentialName,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
					backupCredentialName2,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
					projectName,
					"backup_credential_id = resource.taikun_backup_credential.foo.id",
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
					resource.TestCheckResourceAttrPair("taikun_project.foo", "backup_credential_id", "taikun_backup_credential.foo", "id"),
				),
			},
		},
	})
}

const testAccResourceTaikunProjectConfigWithFlavors = `
resource "taikun_cloud_credential_aws" "foo" {
  name = "%s"
}
data "taikun_flavors" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id
  min_cpu = %d
  max_cpu = %d
}
locals {
  flavors = [for flavor in data.taikun_flavors.foo.flavors: flavor.name]
}
resource "taikun_project" "foo" {
  name = "%s"
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id
  flavors = local.flavors
}
`

func TestAccResourceTaikunProjectModifyFlavors(t *testing.T) {
	cloudCredentialName := randomTestName()
	projectName := randomTestName()
	cpuCount := 2
	newCpuCount := 8
	checkFunc := resource.ComposeAggregateTestCheckFunc(
		testAccCheckTaikunProjectExists,
		resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
		resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
		resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
		resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
		resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
		resource.TestCheckResourceAttrPair("taikun_project.foo", "flavors.#", "data.taikun_flavors.foo", "flavors.#"),
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfigWithFlavors,
					cloudCredentialName,
					cpuCount, cpuCount,
					projectName),
				Check: checkFunc,
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfigWithFlavors,
					cloudCredentialName,
					newCpuCount, newCpuCount,
					projectName),
				Check: checkFunc,
			},
		},
	})
}

const testAccResourceTaikunProjectMinimal = `
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
  cloud_credential_id = resource.taikun_cloud_credential_openstack.foo.id
  flavors = local.flavors
  backup_credential_id = resource.taikun_backup_credential.foo.id
  policy_profile_id = resource.taikun_policy_profile.foo.id

  quota_cpu_units = 64
  quota_ram_size = 256
  quota_disk_size = 512

  server_bastion {
     name = "b"
     disk_size = 30
     flavor = local.flavors[0]
  }
  server_kubeworker {
     name = "w"
     disk_size = 30
     flavor = local.flavors[0]
  }
  server_kubemaster {
     name = "m"
     disk_size = 30
     flavor = local.flavors[0]
  }
}

resource "taikun_kubeconfig" "view" {
  project_id = resource.taikun_project.foo.id

  name = "view-all"

  role = "view"
  access_scope = "all"
  namespace = "default"
}

resource "taikun_kubeconfig" "edit" {
  project_id = resource.taikun_project.foo.id

  name = "edit-all"

  role = "edit"
  access_scope = "all"
  validity_period = 1440
}

resource "taikun_kubeconfig" "admin" {
  project_id = resource.taikun_project.foo.id

  name = "admin-managers"

  role = "admin"
  access_scope = "managers"
}

resource "taikun_kubeconfig" "cluster_admin" {
  project_id = resource.taikun_project.foo.id

  name = "cluster-admin-personal"

  role = "cluster-admin"
  access_scope = "personal"
}

data "taikun_kubeconfigs" "foo" {
  depends_on = [
    taikun_kubeconfig.view,
    taikun_kubeconfig.edit,
    taikun_kubeconfig.admin,
    taikun_kubeconfig.cluster_admin,
  ]

  project_id = resource.taikun_project.foo.id
}

resource "taikun_backup_credential" "foo" {
  name            = "%s"

  s3_endpoint = "%s"
  s3_region   = "%s"
}

resource "taikun_backup_policy" "foo" {
  name = "%s"
  project_id = resource.taikun_project.foo.id
  cron_period = "0 0 * * 0"
  retention_period = "2h"
  included_namespaces = ["test"]
  excluded_namespaces = ["aled"]
}

resource "taikun_policy_profile" "foo" {
  name = "%s"

  forbid_node_port = %t
  forbid_http_ingress = %t
  require_probe = %t
  unique_ingress = %t
  unique_service_selector = %t

}
`

func TestAccResourceTaikunProjectMinimal(t *testing.T) {
	cloudCredentialName := randomTestName()
	backupCredentialName := randomTestName()
	backupPolicyName := randomTestName()
	projectName := shortRandomTestName()
	OPAProfileName := randomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckOpenStack(t); testAccPreCheckS3(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectMinimal,
					cloudCredentialName,
					projectName,
					backupCredentialName,
					os.Getenv("S3_ENDPOINT"),
					os.Getenv("S3_REGION"),
					backupPolicyName,
					OPAProfileName,
					true,
					false,
					true,
					false,
					true,
				),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunProjectExists,
					resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "auto_upgrade"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "monitoring"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
					resource.TestCheckResourceAttr("taikun_project.foo", "server_bastion.#", "1"),
					resource.TestCheckResourceAttr("taikun_project.foo", "server_kubeworker.#", "1"),
					resource.TestCheckResourceAttr("taikun_project.foo", "server_kubemaster.#", "1"),
					resource.TestCheckResourceAttr("data.taikun_kubeconfigs.foo", "kubeconfigs.#", "4"),
					resource.TestCheckResourceAttrSet("taikun_kubeconfig.view", "content"),
					resource.TestCheckResourceAttrSet("taikun_kubeconfig.edit", "content"),
					resource.TestCheckResourceAttrSet("taikun_kubeconfig.admin", "content"),
					resource.TestCheckResourceAttrSet("taikun_kubeconfig.cluster_admin", "content"),
					resource.TestCheckResourceAttr("taikun_project.foo", "quota_cpu_units", "64"),
					resource.TestCheckResourceAttr("taikun_project.foo", "quota_ram_size", "256"),
					resource.TestCheckResourceAttr("taikun_project.foo", "quota_disk_size", "512"),
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
