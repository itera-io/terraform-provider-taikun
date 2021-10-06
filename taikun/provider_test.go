package taikun

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"testing"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"taikun": func() (*schema.Provider, error) {
		return Provider(), nil
	},
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
