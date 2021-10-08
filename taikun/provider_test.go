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
