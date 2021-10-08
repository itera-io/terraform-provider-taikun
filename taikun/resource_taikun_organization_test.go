package taikun

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient/client/organizations"

	"github.com/itera-io/taikungoclient/models"
)

func init() {
	resource.AddTestSweepers("taikun_organization", &resource.Sweeper{
		Name: "taikun_organization",
		F: func(r string) error {

			meta, err := sharedConfig()
			if err != nil {
				return err
			}
			apiClient := meta.(*apiClient)

			params := organizations.NewOrganizationsListParams().WithV(ApiVersion)

			var organizationsList []*models.OrganizationDetailsDto

			for {
				response, err := apiClient.client.Organizations.OrganizationsList(params, apiClient)
				if err != nil {
					return err
				}
				organizationsList = append(organizationsList, response.GetPayload().Data...)
				if len(organizationsList) == int(response.GetPayload().TotalCount) {
					break
				}
				offset := int32(len(organizationsList))
				params = params.WithOffset(&offset)
			}

			for _, e := range organizationsList {
				if strings.HasPrefix(e.Name, testNamePrefix) {
					params := organizations.NewOrganizationsDeleteParams().WithV(ApiVersion).WithOrganizationID(e.ID)
					_, _, err = apiClient.client.Organizations.OrganizationsDelete(params, apiClient)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	})
}

const testAccResourceTaikunOrganization = `
resource "taikun_organization" "foo" {
  name = "%s"
  full_name = "%s"
  discount_rate = %f

  vat_number = %s
  email = "%s"
  billing_email = "%s"
  phone = "%s"
  address = "%s"
  city = "%s"
  country = "%s"

  is_locked = %t

  let_managers_change_subscription = %t
}
`

func TestAccResourceTaikunOrganization(t *testing.T) {
	name := randomTestName()
	fullName := randomString()
	discountRate := rand.Float64() * 100
	vatNumber := randomString()
	email := "manager@example.org"
	billingEmail := "billing@example.org"
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
				Config: fmt.Sprintf(testAccResourceTaikunOrganization,
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunOrganizationExists,
					resource.TestCheckResourceAttr("taikun_organization.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_organization.foo", "full_name", fullName),
					resource.TestCheckResourceAttr("taikun_organization.foo", "discount_rate", discountRate),
					resource.TestCheckResourceAttr("taikun_organization.foo", "vat_number", vatNumber),
					resource.TestCheckResourceAttr("taikun_organization.foo", "email", email),
					resource.TestCheckResourceAttr("taikun_organization.foo", "billing_email", billingEmail),
					resource.TestCheckResourceAttr("taikun_organization.foo", "phone", phone),
					resource.TestCheckResourceAttr("taikun_organization.foo", "address", address),
					resource.TestCheckResourceAttr("taikun_organization.foo", "city", city),
					resource.TestCheckResourceAttr("taikun_organization.foo", "country", country),
					resource.TestCheckResourceAttr("taikun_organization.foo", "is_locked", isLocked),
					resource.TestCheckResourceAttr("taikun_organization.foo", "let_managers_change_subscription", letManagersChangeSubscription),
				),
			},
		},
	})
}

func TestAccResourceTaikunOrganizationUpdate(t *testing.T) {
	name := randomTestName()
	newName := randomTestName()
	fullName := randomString()
	newFullName := randomString()
	discountRate := rand.Float64() * 100
	newDiscountRate := rand.Float64() * 100
	vatNumber := randomString()
	newVatNumber := randomString()
	email := "manager@example.org"
	newEmail := "manager@example.com"
	billingEmail := "billing@example.org"
	newBillingEmail := "billing@example.com"
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
				Config: fmt.Sprintf(testAccResourceTaikunOrganization,
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunOrganizationExists,
					resource.TestCheckResourceAttr("taikun_organization.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_organization.foo", "full_name", fullName),
					resource.TestCheckResourceAttr("taikun_organization.foo", "discount_rate", discountRate),
					resource.TestCheckResourceAttr("taikun_organization.foo", "vat_number", vatNumber),
					resource.TestCheckResourceAttr("taikun_organization.foo", "email", email),
					resource.TestCheckResourceAttr("taikun_organization.foo", "billing_email", billingEmail),
					resource.TestCheckResourceAttr("taikun_organization.foo", "phone", phone),
					resource.TestCheckResourceAttr("taikun_organization.foo", "address", address),
					resource.TestCheckResourceAttr("taikun_organization.foo", "city", city),
					resource.TestCheckResourceAttr("taikun_organization.foo", "country", country),
					resource.TestCheckResourceAttr("taikun_organization.foo", "is_locked", isLocked),
					resource.TestCheckResourceAttr("taikun_organization.foo", "let_managers_change_subscription", letManagersChangeSubscription),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunOrganization,
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaikunOrganizationExists,
					resource.TestCheckResourceAttr("taikun_organization.foo", "name", newName),
					resource.TestCheckResourceAttr("taikun_organization.foo", "full_name", newFullName),
					resource.TestCheckResourceAttr("taikun_organization.foo", "discount_rate", newDiscountRate),
					resource.TestCheckResourceAttr("taikun_organization.foo", "vat_number", newVatNumber),
					resource.TestCheckResourceAttr("taikun_organization.foo", "email", newEmail),
					resource.TestCheckResourceAttr("taikun_organization.foo", "billing_email", newBillingEmail),
					resource.TestCheckResourceAttr("taikun_organization.foo", "phone", newPhone),
					resource.TestCheckResourceAttr("taikun_organization.foo", "address", newAddress),
					resource.TestCheckResourceAttr("taikun_organization.foo", "city", newCity),
					resource.TestCheckResourceAttr("taikun_organization.foo", "country", newCountry),
					resource.TestCheckResourceAttr("taikun_organization.foo", "is_locked", newIsLocked),
					resource.TestCheckResourceAttr("taikun_organization.foo", "let_managers_change_subscription", newLetManagersChangeSubscription),
				),
			},
		},
	})
}

func testAccCheckTaikunOrganizationExists(state *terraform.State) error {
	apiClient := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_organization" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		var limit int32 = 1
		params := organizations.NewOrganizationsListParams().WithV(ApiVersion).WithSearchID(&id).WithLimit(&limit)

		response, err := apiClient.client.Organizations.OrganizationsList(params, apiClient)
		if err != nil || len(response.Payload.Data) != 1 {
			return fmt.Errorf("organization doesn't exist")
		}
	}

	return nil
}

func testAccCheckTaikunOrganizationDestroy(state *terraform.State) error {
	apiClient := testAccProvider.Meta().(*apiClient)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_organization" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		var limit int32 = 1
		params := organizations.NewOrganizationsListParams().WithV(ApiVersion).WithID(&id).WithLimit(&limit)

		response, err := apiClient.client.Organizations.OrganizationsList(params, apiClient)
		if err == nil && len(response.Payload.Data) != 0 {
			return fmt.Errorf("organization still exists")
		}
	}

	return nil
}
