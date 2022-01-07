package taikun

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccResourceTaikunProjectConfigWithImages = `
resource "taikun_cloud_credential_openstack" "foo" {
  name = "%s"
}
data "taikun_images" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_openstack.foo.id
}
locals {
  images = [for image in data.taikun_images.foo.images: image.id]
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
		resource.TestCheckResourceAttrPair("taikun_project.foo", "images.#", "data.taikun_images.foo", "images.#"),
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
  max_ram = 8
}
data "taikun_images" "foo" {
  cloud_credential_id = resource.taikun_cloud_credential_openstack.foo.id
}

locals {
  images = [for image in data.taikun_images.foo.images: image.id]
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

  vm {
    name = "my-vm"
    flavor = local.flavors[0]
    image = local.images[0]
    standalone_profile_id =  resource.taikun_standalone_profile.foo.id
    volume_size = 30
    disk {
      name = "mydisk"
      size = 30
    }
    disk {
      name = "mydisk2"
      size = 30
      volume_type = "ssd-2000iops"
      lun_id = 10
      device_name = "/dev/sdc"
    }
    tag {
      key = "key"
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
