package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]*schema.Provider
var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {

	//  Probably leftover from testing. Default endpoint is production.
	// TF provider is now respecting endpoint in ENV, but provider configuration has precedence.
	// os.Setenv("TAIKUN_API_HOST", "api.taikun.dev")

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
	// What enviroment variables do we require to be set
	requiredEnvSlice := []string{
		"TAIKUN_EMAIL",
		"TAIKUN_PASSWORD",
	}

	// Check if all are set
	checkEnvVariables(requiredEnvSlice, t)
}

func testAccPreCheckPrometheus(t *testing.T) {
	// What enviroment variables do we require to be set
	requiredEnvSlice := []string{
		"PROMETHEUS_URL",
		"PROMETHEUS_USERNAME",
		"PROMETHEUS_PASSWORD",
	}

	// Check if all are set
	checkEnvVariables(requiredEnvSlice, t)
}

func testAccPreCheckOpenStack(t *testing.T) {
	// What enviroment variables do we require to be set
	requiredEnvSlice := []string{
		"OS_AUTH_URL",
		"OS_USERNAME",
		"OS_PASSWORD",
		"OS_USER_DOMAIN_NAME",
		"OS_PROJECT_NAME",
		"OS_INTERFACE",
		"OS_REGION_NAME",
	}

	// Check if all are set
	checkEnvVariables(requiredEnvSlice, t)
}

func testAccPreCheckAWS(t *testing.T) {
	// What enviroment variables do we require to be set
	requiredEnvSlice := []string{
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"AWS_DEFAULT_REGION",
		"AWS_AZ_COUNT",
	}

	// Check if all are set
	checkEnvVariables(requiredEnvSlice, t)
}

func testAccPreCheckGCP(t *testing.T) {
	// What enviroment variables do we require to be set
	requiredEnvSlice := []string{
		"GCP_CONFIG_FILE",
		"GCP_REGION",
	}

	// Check if all are set
	checkEnvVariables(requiredEnvSlice, t)
}

func testAccPreCheckS3(t *testing.T) {
	// What enviroment variables do we require to be set
	requiredEnvSlice := []string{
		"S3_ACCESS_KEY_ID",
		"S3_SECRET_ACCESS_KEY",
		"S3_ENDPOINT",
		"S3_REGION",
	}

	// Check if all are set
	checkEnvVariables(requiredEnvSlice, t)
}

func testAccPreCheckAzure(t *testing.T) {
	// What enviroment variables do we require to be set
	requiredEnvSlice := []string{
		"AZURE_SUBSCRIPTION",
		"AZURE_CLIENT_ID",
		"AZURE_TENANT",
		"AZURE_SECRET",
		"AZURE_AZ_COUNT",
		"AZURE_LOCATION",
	}

	// Check if all are set
	checkEnvVariables(requiredEnvSlice, t)
}

func testAccPreCheckProxmox(t *testing.T) {
	// What enviroment variables do we require to be set
	requiredEnvSlice := []string{
		"PROXMOX_API_HOST",
		"PROXMOX_CLIENT_ID",
		"PROXMOX_CLIENT_SECRET",
		"PROXMOX_STORAGE",
		"PROXMOX_VM_TEMPLATE_NAME",
		"PROXMOX_PUBLIC_NETWORK",
		"PROXMOX_PUBLIC_NETMASK",
		"PROXMOX_PUBLIC_GATEWAY",
		"PROXMOX_PUBLIC_BEGIN_RANGE",
		"PROXMOX_PUBLIC_END_RANGE",
		"PROXMOX_PUBLIC_BRIDGE",
		"PROXMOX_PRIVATE_NETWORK",
		"PROXMOX_PRIVATE_NETMASK",
		"PROXMOX_PRIVATE_GATEWAY",
		"PROXMOX_PRIVATE_BEGIN_RANGE",
		"PROXMOX_PRIVATE_END_RANGE",
		"PROXMOX_PRIVATE_BRIDGE",
		"PROXMOX_HYPERVISOR",
		"PROXMOX_HYPERVISOR2",
	}

	// Check if all are set
	checkEnvVariables(requiredEnvSlice, t)
}

func checkEnvVariables(requiredEnvSlice []string, t *testing.T) {
	// Iterate through the required enviroment variables and check if all are set.
	for _, requiredEnv := range requiredEnvSlice {
		if err := os.Getenv(requiredEnv); err == "" {
			fatalString := fmt.Sprintf("%s must be set for acceptance tests", requiredEnv)
			t.Fatal(fatalString)
		}
	}
}
