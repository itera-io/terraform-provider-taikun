package testing

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	tk "github.com/itera-io/taikungoclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils_testing"
	"math"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccResourceTaikunOrganizationConfig = `
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

  managers_can_change_subscription = %t

  %s
}
`

func TestAccResourceTaikunOrganization(t *testing.T) {
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
	letManagersChangeSubscription := rand.Int()%2 == 0

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunOrganizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunOrganizationConfig,
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
					letManagersChangeSubscription,
					""),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunOrganizationExists,
					resource.TestCheckResourceAttr("taikun_organization.foo", "name", fmt.Sprint(name)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "full_name", fmt.Sprint(fullName)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "discount_rate", fmt.Sprint(discountRate)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "vat_number", fmt.Sprint(vatNumber)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "email", fmt.Sprint(email)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "billing_email", fmt.Sprint(billingEmail)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "phone", fmt.Sprint(phone)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "address", fmt.Sprint(address)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "city", fmt.Sprint(city)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "country", fmt.Sprint(country)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "lock", fmt.Sprint(false)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "managers_can_change_subscription", fmt.Sprint(letManagersChangeSubscription)),
				),
			},
			{
				ResourceName:      "taikun_organization.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceTaikunOrganizationUpdate(t *testing.T) {
	name := utils.RandomTestName()
	newName := utils.RandomTestName()
	fullName := utils.RandomString()
	newFullName := utils.RandomString()
	discountRate := math.Round(rand.Float64()*10000) / 100
	newDiscountRate := math.Round(rand.Float64()*10000) / 100
	vatNumber := utils.RandomString()
	newVatNumber := utils.RandomString()
	email := utils.RandomEmail()
	newEmail := utils.RandomEmail()
	billingEmail := utils.RandomEmail()
	newBillingEmail := utils.RandomEmail()
	phone := "+42424242424242"
	newPhone := "+43434343434343"
	address := "10 Downing Street"
	newAddress := "1600 Pennsylvania Avenue NW"
	city := "London"
	newCity := "Washington, D.C"
	country := "United Kingdom of Great Britain and Northern Ireland"
	newCountry := "United States of America"
	letManagersChangeSubscription := rand.Int()%2 == 0
	newLetManagersChangeSubscription := !letManagersChangeSubscription

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { utils_testing.TestAccPreCheck(t) },
		ProviderFactories: utils_testing.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunOrganizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunOrganizationConfig,
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
					letManagersChangeSubscription,
					""),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunOrganizationExists,
					resource.TestCheckResourceAttr("taikun_organization.foo", "name", fmt.Sprint(name)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "full_name", fmt.Sprint(fullName)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "discount_rate", fmt.Sprint(discountRate)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "vat_number", fmt.Sprint(vatNumber)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "email", fmt.Sprint(email)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "billing_email", fmt.Sprint(billingEmail)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "phone", fmt.Sprint(phone)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "address", fmt.Sprint(address)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "city", fmt.Sprint(city)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "country", fmt.Sprint(country)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "lock", fmt.Sprint(false)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "managers_can_change_subscription", fmt.Sprint(letManagersChangeSubscription)),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunOrganizationConfig,
					newName,
					newFullName,
					newDiscountRate,
					newVatNumber,
					newEmail,
					newBillingEmail,
					newPhone,
					newAddress,
					newCity,
					newCountry,
					newLetManagersChangeSubscription,
					"lock = true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunOrganizationExists,
					resource.TestCheckResourceAttr("taikun_organization.foo", "name", fmt.Sprint(newName)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "full_name", fmt.Sprint(newFullName)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "discount_rate", fmt.Sprint(newDiscountRate)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "vat_number", fmt.Sprint(newVatNumber)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "email", fmt.Sprint(newEmail)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "billing_email", fmt.Sprint(newBillingEmail)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "phone", fmt.Sprint(newPhone)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "address", fmt.Sprint(newAddress)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "city", fmt.Sprint(newCity)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "country", fmt.Sprint(newCountry)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "lock", fmt.Sprint(true)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "managers_can_change_subscription", fmt.Sprint(newLetManagersChangeSubscription)),
				),
			},
		},
	})
}

func testAccCheckTaikunOrganizationExists(state *terraform.State) error {
	apiClient := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_organization" {
			continue
		}

		id, _ := utils.Atoi32(rs.Primary.ID)
		response, _, err := apiClient.Client.OrganizationsAPI.OrganizationsList(context.TODO()).Id(id).Execute()
		if err != nil || len(response.GetData()) != 1 {
			return fmt.Errorf("organization doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunOrganizationDestroy(state *terraform.State) error {
	apiClient := utils_testing.TestAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_organization" {
			continue
		}

		retryErr := retry.RetryContext(context.Background(), utils.GetReadAfterOpTimeout(false), func() *retry.RetryError {
			id, _ := utils.Atoi32(rs.Primary.ID)
			response, _, err := apiClient.Client.OrganizationsAPI.OrganizationsList(context.TODO()).Id(id).Execute()

			if err != nil {
				return retry.NonRetryableError(err)
			}
			if response.GetTotalCount() != 0 {
				return retry.RetryableError(errors.New("organization still exists"))
			}
			return nil
		})
		if utils.TimedOut(retryErr) {
			return errors.New("organization still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
