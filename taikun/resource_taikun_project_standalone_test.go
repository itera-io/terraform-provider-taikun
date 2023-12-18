package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccResourceTaikunProjectConfigWithImages = `
resource "taikun_cloud_credential_openstack" "foo" {
  name = "%s"
}
data "taikun_images_openstack" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_openstack.foo.id
}
locals {
  images = [for image in data.taikun_images_openstack.foo.images: image.id]
}
resource "taikun_project" "foo" {
  name = "%s"
  cloud_credential_id = resource.taikun_cloud_credential_openstack.foo.id
  images = local.images
}
`

func TestAccResourceTaikunProjectModifyImages(t *testing.T) {
	cloudCredentialName := randomTestName()
	projectName := randomTestName()
	checkFunc := resource.ComposeAggregateTestCheckFunc(
		testAccCheckTaikunProjectExists,
		resource.TestCheckResourceAttr("taikun_project.foo", "name", projectName),
		resource.TestCheckResourceAttrSet("taikun_project.foo", "access_profile_id"),
		resource.TestCheckResourceAttrSet("taikun_project.foo", "cloud_credential_id"),
		resource.TestCheckResourceAttrSet("taikun_project.foo", "kubernetes_profile_id"),
		resource.TestCheckResourceAttrSet("taikun_project.foo", "organization_id"),
		resource.TestCheckResourceAttrPair("taikun_project.foo", "images.#", "data.taikun_images_openstack.foo", "images.#"),
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckOpenStack(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfigWithImages,
					cloudCredentialName,
					projectName),
				Check: checkFunc,
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectConfigWithImages,
					cloudCredentialName,
					projectName),
				Check: checkFunc,
			},
		},
	})
}

const testAccResourceTaikunProjectStandaloneOpenStackMinimal = `
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

data "taikun_images_openstack" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_openstack.foo.id
}

locals {
  images = [for image in data.taikun_images_openstack.foo.images: image.id]
  flavors = [for flavor in data.taikun_flavors.foo.flavors: flavor.name]
}

resource "taikun_standalone_profile" "foo" {
	name = "%s"
    public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"
}

resource "taikun_project" "foo" {
  name = "%s"
  cloud_credential_id = resource.taikun_cloud_credential_openstack.foo.id
  flavors = local.flavors
  images = local.images

  quota_vm_cpu_units = 64
  quota_vm_ram_size = 256
  quota_vm_volume_size = 512

  vm {
    name = "tf-acc-vm"
    flavor = local.flavors[%d]
    image_id = local.images[0]
    standalone_profile_id =  resource.taikun_standalone_profile.foo.id
    volume_size = 60
    %s
    disk {
      name = "tf-acc-disk"
      size = 30
    }
    disk {
      name = "tf-acc-disk2"
      size = 30
      volume_type = "ssd"
    }
    tag {
      key = "key"
      value = "value"
    }
    tag {
      key = "key2"
      value = "value"
    }
  }
}
`

func TestAccResourceTaikunProjectStandaloneOpenStackMinimal(t *testing.T) {
	cloudCredentialName := randomTestName()
	standaloneProfileName := randomTestName()
	projectName := shortRandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckOpenStack(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectStandaloneOpenStackMinimal,
					cloudCredentialName,
					standaloneProfileName,
					projectName,
					0,
					"",
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
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.#", "1"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.volume_size", "60"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.public_ip", "false"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.access_ip", ""),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.tag.#", "2"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.disk.#", "2"),
					resource.TestCheckResourceAttr("taikun_project.foo", "quota_vm_cpu_units", "64"),
					resource.TestCheckResourceAttr("taikun_project.foo", "quota_vm_ram_size", "256"),
					resource.TestCheckResourceAttr("taikun_project.foo", "quota_vm_volume_size", "512"),
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

func TestAccResourceTaikunProjectStandaloneOpenStackMinimalUpdateIP(t *testing.T) {
	cloudCredentialName := randomTestName()
	standaloneProfileName := randomTestName()
	projectName := shortRandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckOpenStack(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectStandaloneOpenStackMinimal,
					cloudCredentialName,
					standaloneProfileName,
					projectName,
					0,
					"",
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
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.#", "1"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.volume_size", "60"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.public_ip", "false"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.access_ip", ""),
					//resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.tag.#", "2"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.disk.#", "2"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectStandaloneOpenStackMinimal,
					cloudCredentialName,
					standaloneProfileName,
					projectName,
					0,
					"public_ip = true",
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
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.#", "1"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.volume_size", "60"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.public_ip", "true"),
					resource.TestCheckResourceAttrSet("taikun_project.foo", "vm.0.access_ip"),
					//resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.tag.#", "2"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.disk.#", "2"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectStandaloneOpenStackMinimal,
					cloudCredentialName,
					standaloneProfileName,
					projectName,
					0,
					"",
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
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.#", "1"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.volume_size", "60"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.public_ip", "false"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.access_ip", ""),
					//resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.tag.#", "2"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.disk.#", "2"),
				),
			},
		},
	})
}

func TestAccResourceTaikunProjectStandaloneOpenStackMinimalUpdateFlavor(t *testing.T) {
	cloudCredentialName := randomTestName()
	standaloneProfileName := randomTestName()
	projectName := shortRandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckOpenStack(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectStandaloneOpenStackMinimal,
					cloudCredentialName,
					standaloneProfileName,
					projectName,
					0,
					"",
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
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.#", "1"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.volume_size", "60"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.public_ip", "false"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.access_ip", ""),
					//resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.tag.#", "2"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.disk.#", "2"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectStandaloneOpenStackMinimal,
					cloudCredentialName,
					standaloneProfileName,
					projectName,
					1,
					"",
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
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.#", "1"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.volume_size", "60"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.volume_size", "60"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.public_ip", "false"),
					//resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.tag.#", "2"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.disk.#", "2"),
				),
			},
		},
	})
}

func TestAccResourceTaikunProjectStandaloneOpenStackMinimalWithVolumeType(t *testing.T) {
	cloudCredentialName := randomTestName()
	standaloneProfileName := randomTestName()
	projectName := shortRandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckOpenStack(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectStandaloneOpenStackMinimal,
					cloudCredentialName,
					standaloneProfileName,
					projectName,
					0,
					"volume_type = \"ssd\"",
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
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.#", "1"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.volume_size", "60"),
					//resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.tag.#", "2"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.disk.#", "2"),
				),
			},
		},
	})
}

const testAccResourceTaikunProjectStandaloneAWSMinimal = `
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

resource "taikun_standalone_profile" "foo" {
	name = "%s"
    public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"
}

resource "taikun_project" "foo" {
  name = "%s"
  cloud_credential_id = resource.taikun_cloud_credential_aws.foo.id
  flavors = local.flavors
  images = ["ami-0f94f6b47f28fb0e7"]

  vm {
    name = "tf-acc-vm"
    flavor = local.flavors[%d]
    image_id = "ami-0f94f6b47f28fb0e7"
    standalone_profile_id =  resource.taikun_standalone_profile.foo.id
    volume_size = 60
    disk {
      name = "tf-acc-disk"
      size = 30
    }
    disk {
      name = "tf-acc-disk2"
      size = 30
    }
    //tag {
    //  key = "key"
    //  value = "value"
    //}
    //tag {
    //  key = "key2"
    //  value = "value"
    //}
  }
}
`

func TestAccResourceTaikunProjectStandaloneAWSMinimal(t *testing.T) {
	cloudCredentialName := randomTestName()
	standaloneProfileName := randomTestName()
	projectName := shortRandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectStandaloneAWSMinimal,
					cloudCredentialName,
					standaloneProfileName,
					projectName,
					0,
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
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.#", "1"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.volume_size", "60"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.public_ip", "false"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.access_ip", ""),
					//resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.tag.#", "2"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.disk.#", "2"),
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

func TestAccResourceTaikunProjectStandaloneAWSMinimalUpdateFlavor(t *testing.T) {
	cloudCredentialName := randomTestName()
	standaloneProfileName := randomTestName()
	projectName := shortRandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAWS(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectStandaloneAWSMinimal,
					cloudCredentialName,
					standaloneProfileName,
					projectName,
					0,
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
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.#", "1"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.volume_size", "60"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.public_ip", "false"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.access_ip", ""),
					//resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.tag.#", "2"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.disk.#", "2"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectStandaloneAWSMinimal,
					cloudCredentialName,
					standaloneProfileName,
					projectName,
					1,
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
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.#", "1"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.volume_size", "60"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.public_ip", "false"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.access_ip", ""),
					//resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.tag.#", "2"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.disk.#", "2"),
				),
			},
		},
	})
}

const testAccResourceTaikunProjectStandaloneAzureMinimal = `
resource "taikun_cloud_credential_azure" "foo" {
  name = "%s"
  location = "%s"
}

data "taikun_flavors" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_azure.foo.id
  min_cpu = 2
  max_cpu = 2
  min_ram = 4
  max_ram = 8
}

data "taikun_images_azure" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_azure.foo.id
  publisher = "Canonical"
  offer = "0001-com-ubuntu-server-jammy"
  sku = "22_04-lts"
  latest = true
}

locals {
  images = [for image in data.taikun_images_azure.foo.images: image.id]
  flavors = [for flavor in data.taikun_flavors.foo.flavors: flavor.name]
}

resource "taikun_standalone_profile" "foo" {
	name = "%s"
    public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDHIHRbKu4shd9OardTAj0yV1HSjcg0o57s27IHT+TpL5CSVd8o5qZl/rI14pFbAG1lCbKly+BI9ql4kEj4RCOd6uS0EnLS3vxH0DPSZqKnV+r+iV8w0/0GgxaihAK2Z7RdIVoizOjDkjCRIDNd9fkQ2/C6uHdDmrRxiFh+e98w7Ebes/xcCX6r0iMhAUkYFfMx7C/H7BANA53YOJBdtxcd1BZbRo5VktoZ0i0ie5d+OioeD1uR+nEnU12q2tJqo4j2WHpJ++Rba2aNesVrYq1V9OoKg3+hl5CFXVDHzcgq2PykfNQ2PKo/C5i3jjLISMSVKvqCJDjZTsJJsoifv5KClkOYGA12Aqe/qJEpeq7uPadbQFRdYK8FT74K71Pz3Qg1Ipy02o6QaNRHZtJyXnaO5TZciD2tiM3YthuMoh0/vnARlqxc2YElOmrfUtaAEv3bB/SiIFreyGgkb1VNkEWA1hQmqYMxnTFhGF0ZbwSLo6xXQRTKuYo39ts+4eaqcJ0= non_non"
}

resource "taikun_project" "foo" {
  name = "%s"
  cloud_credential_id = resource.taikun_cloud_credential_azure.foo.id
  flavors = local.flavors
  images = local.images

  vm {
    name = "tf-acc-vm"
    flavor = local.flavors[%d]
    image_id = local.images[0]
    username = "foobar"
    standalone_profile_id =  resource.taikun_standalone_profile.foo.id
    volume_size = 60
    %s
    disk {
      name = "tf-acc-disk"
      size = 30
    }
    disk {
      name = "tf-acc-disk2"
      size = 30
      volume_type = "Premium_LRS"
    }
    //tag {
    //  key = "key"
    //  value = "value"
    //}
    //tag {
    //  key = "key2"
    //  value = "value"
    //}
  }
}
`

func TestAccResourceTaikunProjectStandaloneAzureMinimal(t *testing.T) {
	cloudCredentialName := randomTestName()
	standaloneProfileName := randomTestName()
	projectName := shortRandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAzure(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectStandaloneAzureMinimal,
					cloudCredentialName,
					os.Getenv("AZURE_LOCATION"),
					standaloneProfileName,
					projectName,
					0,
					"",
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
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.#", "1"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.volume_size", "60"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.public_ip", "false"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.access_ip", ""),
					//resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.tag.#", "2"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.disk.#", "2"),
				),
			},
			{
				ResourceName: "taikun_project.foo",
				ImportState:  true,
			},
		},
	})
}

func TestAccResourceTaikunProjectStandaloneAzureMinimalUpdateFlavor(t *testing.T) {
	cloudCredentialName := randomTestName()
	standaloneProfileName := randomTestName()
	projectName := shortRandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAzure(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectStandaloneAzureMinimal,
					cloudCredentialName,
					os.Getenv("AZURE_LOCATION"),
					standaloneProfileName,
					projectName,
					0,
					"",
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
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.#", "1"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.volume_size", "60"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.public_ip", "false"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.access_ip", ""),
					//resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.tag.#", "2"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.disk.#", "2"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectStandaloneAzureMinimal,
					cloudCredentialName,
					os.Getenv("AZURE_LOCATION"),
					standaloneProfileName,
					projectName,
					1,
					"",
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
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.#", "1"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.volume_size", "60"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.public_ip", "false"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.access_ip", ""),
					//resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.tag.#", "2"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.disk.#", "2"),
				),
			},
		},
	})
}

func TestAccResourceTaikunProjectStandaloneAzureMinimalWithVolumeType(t *testing.T) {
	cloudCredentialName := randomTestName()
	standaloneProfileName := randomTestName()
	projectName := shortRandomTestName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAzure(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunProjectStandaloneAzureMinimal,
					cloudCredentialName,
					os.Getenv("AZURE_LOCATION"),
					standaloneProfileName,
					projectName,
					0,
					"volume_type = \"Premium_LRS\"",
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
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.#", "1"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.volume_size", "60"),
					//resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.tag.#", "2"),
					resource.TestCheckResourceAttr("taikun_project.foo", "vm.0.disk.#", "2"),
				),
			},
		},
	})
}
