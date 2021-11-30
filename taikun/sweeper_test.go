package taikun

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/itera-io/taikungoclient/client"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func sharedConfig() (interface{}, error) {
	email := os.Getenv("TAIKUN_EMAIL")
	password := os.Getenv("TAIKUN_PASSWORD")

	if email == "" {
		return nil, fmt.Errorf("TAIKUN_EMAIL must be set in order to run sweepers")
	}
	if password == "" {
		return nil, fmt.Errorf("TAIKUN_PASSWORD must be set in order to run sweepers")
	}

	transportConfig := client.DefaultTransportConfig().WithHost("api.taikun.dev")

	return &apiClient{
		client:              client.NewHTTPClientWithConfig(nil, transportConfig),
		email:               email,
		password:            password,
		useKeycloakEndpoint: false,
	}, nil
}
