package utils_testing

import (
	"fmt"
	"os"
	"testing"
)

func TestAccPreCheck(t *testing.T) {
	// What enviroment variables do we require to be set
	requiredEnvSlice := []string{
		"TAIKUN_EMAIL",
		"TAIKUN_PASSWORD",
	}

	// Check if all are set
	checkEnvVariables(requiredEnvSlice, t)
}

func TestAccPreCheckPrometheus(t *testing.T) {
	// What enviroment variables do we require to be set
	requiredEnvSlice := []string{
		"PROMETHEUS_URL",
		"PROMETHEUS_USERNAME",
		"PROMETHEUS_PASSWORD",
	}

	// Check if all are set
	checkEnvVariables(requiredEnvSlice, t)
}

func TestAccPreCheckOpenStack(t *testing.T) {
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

func TestAccPreCheckAWS(t *testing.T) {
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

func TestAccPreCheckZadara(t *testing.T) {
	// What enviroment variables do we require to be set
	requiredEnvSlice := []string{
		"ZADARA_ACCESS_KEY_ID",
		"ZADARA_SECRET_ACCESS_KEY",
		"ZADARA_DEFAULT_REGION",
		"ZADARA_AZ_COUNT",
	}

	// Check if all are set
	checkEnvVariables(requiredEnvSlice, t)
}

func TestAccPreCheckGCP(t *testing.T) {
	// What enviroment variables do we require to be set
	requiredEnvSlice := []string{
		"GCP_CONFIG_FILE",
		"GCP_REGION",
	}

	// Check if all are set
	checkEnvVariables(requiredEnvSlice, t)
}

func TestAccPreCheckS3(t *testing.T) {
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

func TestAccPreCheckAzure(t *testing.T) {
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

func TestAccPreCheckProxmox(t *testing.T) {
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

func TestAccPreCheckVsphere(t *testing.T) {
	// What enviroment variables do we require to be set
	requiredEnvSlice := []string{
		"VSPHERE_USERNAME",
		"VSPHERE_PASSWORD",
		"VSPHERE_API_URL",
		"VSPHERE_DATACENTER",
		"VSPHERE_RESOURCE_POOL",
		"VSPHERE_DATA_STORE",
		"VSPHERE_DRS_ENABLED",
		"VSPHERE_HYPERVISOR",
		"VSPHERE_HYPERVISOR2",
		"VSPHERE_VM_TEMPLATE",
		"VSPHERE_CONTINENT",

		"VSPHERE_PUBLIC_NETWORK_NAME",
		"VSPHERE_PUBLIC_NETWORK_ADDRESS",
		"VSPHERE_PUBLIC_NETMASK",
		"VSPHERE_PUBLIC_GATEWAY",
		"VSPHERE_PUBLIC_BEGIN_RANGE",
		"VSPHERE_PUBLIC_END_RANGE",

		"VSPHERE_PRIVATE_NETWORK_NAME",
		"VSPHERE_PRIVATE_NETWORK_ADDRESS",
		"VSPHERE_PRIVATE_NETMASK",
		"VSPHERE_PRIVATE_GATEWAY",
		"VSPHERE_PRIVATE_BEGIN_RANGE",
		"VSPHERE_PRIVATE_END_RANGE",
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
