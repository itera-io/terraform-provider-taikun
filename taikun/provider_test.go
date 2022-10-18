package taikun

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]*schema.Provider
var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {

	os.Setenv("TAIKUN_API_HOST", "api.taikun.dev")

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
	if err := os.Getenv("OS_AUTH_URL"); err == "" {
		t.Fatal("OS_AUTH_URL must be set for acceptance tests")
	}
	if err := os.Getenv("OS_USERNAME"); err == "" {
		t.Fatal("OS_USERNAME must be set for acceptance tests")
	}
	if err := os.Getenv("OS_PASSWORD"); err == "" {
		t.Fatal("OS_PASSWORD must be set for acceptance tests")
	}
	if err := os.Getenv("OS_USER_DOMAIN_NAME"); err == "" {
		t.Fatal("OS_USER_DOMAIN_NAME must be set for acceptance tests")
	}
	if err := os.Getenv("OS_PROJECT_NAME"); err == "" {
		t.Fatal("OS_PROJECT_NAME must be set for acceptance tests")
	}
	if err := os.Getenv("OS_INTERFACE"); err == "" {
		t.Fatal("OS_INTERFACE must be set for acceptance tests")
	}
	if err := os.Getenv("OS_REGION_NAME"); err == "" {
		t.Fatal("OS_REGION_NAME must be set for acceptance tests")
	}
}

func testAccPreCheckAWS(t *testing.T) {
	if err := os.Getenv("AWS_ACCESS_KEY_ID"); err == "" {
		t.Fatal("AWS_ACCESS_KEY_ID must be set for acceptance tests")
	}
	if err := os.Getenv("AWS_SECRET_ACCESS_KEY"); err == "" {
		t.Fatal("AWS_SECRET_ACCESS_KEY must be set for acceptance tests")
	}
	if err := os.Getenv("AWS_DEFAULT_REGION"); err == "" {
		t.Fatal("AWS_DEFAULT_REGION must be set for acceptance tests")
	}
	if err := os.Getenv("AWS_AZ_COUNT"); err == "" {
		t.Fatal("AWS_AZ_COUNT must be set for acceptance tests")
	}
}

func testAccPreCheckGCP(t *testing.T) {
	if err := os.Getenv("GCP_FOLDER_ID"); err == "" {
		t.Fatal("GCP_FOLDER_ID must be set for acceptance tests")
	}
	if err := os.Getenv("GCP_BILLING_ACCOUNT"); err == "" {
		t.Fatal("GCP_BILLING_ACCOUNT must be set for acceptance tests")
	}
	if err := os.Getenv("GCP_REGION"); err == "" {
		t.Fatal("GCP_REGION must be set for acceptance tests")
	}
	if err := os.Getenv("GCP_AZ_COUNT"); err == "" {
		t.Fatal("GCP_ZONE must be set for acceptance tests")
	}
}

func testAccPreCheckS3(t *testing.T) {
	if err := os.Getenv("AWS_ACCESS_KEY_ID"); err == "" {
		t.Fatal("AWS_ACCESS_KEY_ID must be set for acceptance tests")
	}
	if err := os.Getenv("AWS_SECRET_ACCESS_KEY"); err == "" {
		t.Fatal("AWS_SECRET_ACCESS_KEY must be set for acceptance tests")
	}
	if err := os.Getenv("S3_ENDPOINT"); err == "" {
		t.Fatal("S3_ENDPOINT must be set for acceptance tests")
	}
	if err := os.Getenv("S3_REGION"); err == "" {
		t.Fatal("S3_REGION must be set for acceptance tests")
	}
}

func testAccPreCheckAzure(t *testing.T) {
	if err := os.Getenv("ARM_SUBSCRIPTION_ID"); err == "" {
		t.Fatal("ARM_SUBSCRIPTION_ID must be set for acceptance tests")
	}
	if err := os.Getenv("ARM_CLIENT_ID"); err == "" {
		t.Fatal("ARM_CLIENT_ID must be set for acceptance tests")
	}
	if err := os.Getenv("ARM_TENANT_ID"); err == "" {
		t.Fatal("ARM_TENANT_ID must be set for acceptance tests")
	}
	if err := os.Getenv("ARM_CLIENT_SECRET"); err == "" {
		t.Fatal("ARM_CLIENT_SECRET must be set for acceptance tests")
	}
	if err := os.Getenv("ARM_AZ_COUNT"); err == "" {
		t.Fatal("ARM_AZ_COUNT must be set for acceptance tests")
	}
	if err := os.Getenv("ARM_LOCATION"); err == "" {
		t.Fatal("ARM_LOCATION must be set for acceptance tests")
	}
}
