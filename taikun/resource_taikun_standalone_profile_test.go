package taikun

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/stand_alone_profile"
)

const testAccResourceTaikunStandaloneProfileConfig = `
resource "taikun_standalone_profile" "foo" {
	name = "%s"
    public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"
    lock = %t
    security_group {
        name = "http"
        from_port = 80
        to_port = 80
        ip_protocol = "TCP"
        cidr = "0.0.0.0/0"
    }
    security_group {
        name = "https"
        from_port = 443
        to_port = 443
        ip_protocol = "TCP"
        cidr = "0.0.0.0/0"
    }
    %s
}
`

const testAccResourceTaikunStandaloneProfileExtraSecurityGroup = `
security_group {
    name = "http2"
    from_port = 80
    to_port = 80
    ip_protocol = "UDP"
    cidr = "0.0.0.0/0"
}
security_group {
    name = "icmp"
    ip_protocol = "ICMP"
    cidr = "0.0.0.0/0"
}
`

func TestAccResourceTaikunStandaloneProfile(t *testing.T) {
	name := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunStandaloneProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunStandaloneProfileConfig, name, false, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunStandaloneProfileExists,
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "public_key"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_name"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.#", "2"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.name", "http"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.from_port", "80"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.to_port", "80"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.cidr", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.ip_protocol", "TCP"),
				),
			},
			{
				ResourceName:      "taikun_standalone_profile.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceTaikunStandaloneProfileLock(t *testing.T) {
	name := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunStandaloneProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunStandaloneProfileConfig, name, false, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunStandaloneProfileExists,
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "public_key"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_name"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.#", "2"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.name", "http"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.from_port", "80"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.to_port", "80"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.cidr", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.ip_protocol", "TCP"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunStandaloneProfileConfig, name, true, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunStandaloneProfileExists,
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "lock", "true"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "public_key"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_name"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.#", "2"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.name", "http"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.from_port", "80"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.to_port", "80"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.cidr", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.ip_protocol", "TCP"),
				),
			},
		},
	})
}

func TestAccResourceTaikunStandaloneProfileRename(t *testing.T) {
	name := randomTestName()
	newName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunStandaloneProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunStandaloneProfileConfig, name, false, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunStandaloneProfileExists,
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "public_key"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_name"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.#", "2"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.name", "http"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.from_port", "80"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.to_port", "80"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.cidr", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.ip_protocol", "TCP"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunStandaloneProfileConfig, newName, false, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunStandaloneProfileExists,
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "name", newName),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "public_key"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_name"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.#", "2"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.name", "http"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.from_port", "80"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.to_port", "80"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.cidr", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.ip_protocol", "TCP"),
				),
			},
		},
	})
}

func TestAccResourceTaikunStandaloneProfileAddGroups(t *testing.T) {
	name := randomTestName()
	newName := randomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckPrometheus(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckTaikunStandaloneProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccResourceTaikunStandaloneProfileConfig, name, false, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunStandaloneProfileExists,
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "name", name),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "public_key"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_name"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.#", "2"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.name", "http"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.from_port", "80"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.to_port", "80"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.cidr", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.ip_protocol", "TCP"),
				),
			},
			{
				Config: fmt.Sprintf(testAccResourceTaikunStandaloneProfileConfig, newName, false, testAccResourceTaikunStandaloneProfileExtraSecurityGroup),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTaikunStandaloneProfileExists,
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "name", newName),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "lock", "false"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "public_key"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_id"),
					resource.TestCheckResourceAttrSet("taikun_standalone_profile.foo", "organization_name"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.#", "4"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.0.name", "http"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.1.name", "https"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.2.name", "http2"),
					resource.TestCheckResourceAttr("taikun_standalone_profile.foo", "security_group.3.name", "icmp"),
				),
			},
		},
	})
}

func testAccCheckTaikunStandaloneProfileExists(state *terraform.State) error {
	client := testAccProvider.Meta().(*taikungoclient.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_standalone_profile" {
			continue
		}

		id, _ := atoi32(rs.Primary.ID)
		params := stand_alone_profile.NewStandAloneProfileListParams().WithV(ApiVersion).WithID(&id)

		response, err := client.Client.StandAloneProfile.StandAloneProfileList(params, client)
		if err != nil || response.Payload.TotalCount != 1 {
			return fmt.Errorf("standalone profile doesn't exist (id = %s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckTaikunStandaloneProfileDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*taikungoclient.Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "taikun_standalone_profile" {
			continue
		}

		retryErr := resource.RetryContext(context.Background(), getReadAfterOpTimeout(false), func() *resource.RetryError {
			id, _ := atoi32(rs.Primary.ID)
			params := stand_alone_profile.NewStandAloneProfileListParams().WithV(ApiVersion).WithID(&id)

			response, err := client.Client.StandAloneProfile.StandAloneProfileList(params, client)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if response.Payload.TotalCount != 0 {
				return resource.RetryableError(errors.New("standalone profile still exists"))
			}
			return nil
		})
		if timedOut(retryErr) {
			return errors.New("standalone profile still exists (timed out)")
		}
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}
