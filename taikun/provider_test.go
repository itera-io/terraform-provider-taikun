package taikun

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"testing"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]*schema.Provider
var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"taikun": testAccProvider,
	}
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"taikun": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if err := os.Getenv("TAIKUN_EMAIL"); err == "" {
		t.Fatal("TAIKUN_EMAIL must be set for acceptance tests")
	}
	if err := os.Getenv("TAIKUN_PASSWORD"); err == "" {
		t.Fatal("TAIKUN_PASSWORD must be set for acceptance tests")
	}
}

func testAccPreCheckPrometheus(t *testing.T) {
	if err := os.Getenv("PROMETHEUS_URL"); err == "" {
		t.Fatal("PROMETHEUS_URL must be set for acceptance tests")
	}
	if err := os.Getenv("PROMETHEUS_USERNAME"); err == "" {
		t.Fatal("PROMETHEUS_USERNAME must be set for acceptance tests")
	}
	if err := os.Getenv("PROMETHEUS_PASSWORD"); err == "" {
		t.Fatal("PROMETHEUS_PASSWORD must be set for acceptance tests")
	}
}

func testAccPreCheckOpenStack(t *testing.T) {
	if err := os.Getenv("OPENSTACK_URL"); err == "" {
		t.Fatal("OPENSTACK_URL must be set for acceptance tests")
	}
	if err := os.Getenv("OPENSTACK_USERNAME"); err == "" {
		t.Fatal("OPENSTACK_USERNAME must be set for acceptance tests")
	}
	if err := os.Getenv("OPENSTACK_PASSWORD"); err == "" {
		t.Fatal("OPENSTACK_PASSWORD must be set for acceptance tests")
	}
	if err := os.Getenv("OPENSTACK_DOMAIN"); err == "" {
		t.Fatal("OPENSTACK_DOMAIN must be set for acceptance tests")
	}
	if err := os.Getenv("OPENSTACK_PROJECT"); err == "" {
		t.Fatal("OPENSTACK_PROJECT must be set for acceptance tests")
	}
	if err := os.Getenv("OPENSTACK_PUBLIC_NETWORK"); err == "" {
		t.Fatal("OPENSTACK_PUBLIC_NETWORK must be set for acceptance tests")
	}
	if err := os.Getenv("OPENSTACK_REGION"); err == "" {
		t.Fatal("OPENSTACK_REGION must be set for acceptance tests")
	}
}
