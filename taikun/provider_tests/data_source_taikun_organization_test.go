package provider_tests

import (
	"fmt"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"math"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataSourceOrganizationConfig = `
data "taikun_organization" "foo" {
}
`

func TestAccDataSourceTaikunOrganization(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOrganizationConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "discount_rate"),
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "name"),
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "full_name"),
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "lock"),
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "id"),
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "is_read_only"),
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "projects"),
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "servers"),
					resource.TestCheckResourceAttrSet("data.taikun_organization.foo", "users"),
				),
			},
		},
	})
}

const testAccDataSourceOrganizationNewConfig = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = %f

  vat_number = "%s"
  email = "%s"
  billing_email = "%s"
  phone = "%s"
  address = "%s"
  city = "%s"
  country = "%s"
}

data "taikun_organization" "foo" {
  id = resource.taikun_organization.foo.id
}
`

func TestAccDataSourceTaikunOrganizationNew(t *testing.T) {
	name := utils.RandomTestName()
	fullName := utils.RandomString()
	discountRate := math.Round(rand.Float64()*10000) / 100
	vatNumber := utils.RandomString()
	email := utils.RandomEmail()
	billingEmail := utils.RandomEmail()
	phone := "+42424242424242"
	address := "10 Downing Street"
	city := "London"
	country := "United Kingdom of Great Britain and Northern Ireland"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceOrganizationNewConfig,
					name,
					fullName,
					discountRate,
					vatNumber,
					email,
					billingEmail,
					phone,
					address,
					city,
					country,
				),
				Check: utils_testing.CheckDataSourceStateMatchesResourceState(
					"data.taikun_organization.foo",
					"taikun_organization.foo",
				),
			},
		},
	})
}
