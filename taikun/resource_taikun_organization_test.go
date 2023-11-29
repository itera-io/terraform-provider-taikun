package taikun

import (
	"context"
	"errors"
	"fmt"
	tk "github.com/itera-io/taikungoclient"
	"math"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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

  lock = %t

  managers_can_change_subscription = %t
}
`

func TestAccResourceTaikunOrganization(t *testing.T) {
	name := randomTestName()
	fullName := randomString()
	discountRate := math.Round(rand.Float64()*10000) / 100
	vatNumber := randomString()
	email := randomEmail()
	billingEmail := randomEmail()
	phone := "+42424242424242"
	address := "10 Downing Street"
	city := "London"
	country := "United Kingdom of Great Britain and Northern Ireland"
	isLocked := rand.Int()%2 == 0
	letManagersChangeSubscription := rand.Int()%2 == 0

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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
					isLocked,
					letManagersChangeSubscription),
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
					resource.TestCheckResourceAttr("taikun_organization.foo", "lock", fmt.Sprint(isLocked)),
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
	name := randomTestName()
	newName := randomTestName()
	fullName := randomString()
	newFullName := randomString()
	discountRate := math.Round(rand.Float64()*10000) / 100
	newDiscountRate := math.Round(rand.Float64()*10000) / 100
	vatNumber := randomString()
	newVatNumber := randomString()
	email := randomEmail()
	newEmail := randomEmail()
	billingEmail := randomEmail()
	newBillingEmail := randomEmail()
	phone := "+42424242424242"
	newPhone := "+43434343434343"
	address := "10 Downing Street"
	newAddress := "1600 Pennsylvania Avenue NW"
	city := "London"
	newCity := "Washington, D.C"
	country := "United Kingdom of Great Britain and Northern Ireland"
	newCountry := "United States of America"
	isLocked := rand.Int()%2 == 0
	newIsLocked := !isLocked
	letManagersChangeSubscription := rand.Int()%2 == 0
	newLetManagersChangeSubscription := !letManagersChangeSubscription

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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
					isLocked,
					letManagersChangeSubscription),
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
					resource.TestCheckResourceAttr("taikun_organization.foo", "lock", fmt.Sprint(isLocked)),
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
					newIsLocked,
					newLetManagersChangeSubscription),
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
					resource.TestCheckResourceAttr("taikun_organization.foo", "lock", fmt.Sprint(newIsLocked)),
					resource.TestCheckResourceAttr("taikun_organization.foo", "managers_can_change_subscription", fmt.Sprint(newLetManagersChangeSubscription)),
				),
			},
		},
	})
}

func testAccCheckTaikunOrganizationExists(state *terraform.State) error {
	apiClient := testAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_organization" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		response, _, err := apiClient.Client.OrganizationsAPI.OrganizationsList(context.TODO()).Id(id).Execute()
		if err != nil || len(response.GetData()) != 1 {
			return fmt.Errorf("organization doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunOrganizationDestroy(state *terraform.State) error {
	apiClient := testAccProvider.Meta().(*tk.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_organization" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			id, _ := atoi32(rs.Primary.ID)
			response, _, err := apiClient.Client.OrganizationsAPI.OrganizationsList(context.TODO()).Id(id).Execute()

			if err != nil {
				return resource.NonRetryableError(err)
			}
			if response.GetTotalCount() != 0 {
				return resource.RetryableError(errors.New("organization still exists"))
			}
			return nil
		})
		if timedOut(retryErr) {
			return errors.New("organization still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
