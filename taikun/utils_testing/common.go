package utils_testing

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/terraform-provider-taikun/taikun/provider"
)

var TestAccProvider *schema.Provider
var TestAccProviders map[string]*schema.Provider
var TestAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	TestAccProvider = provider.Provider()
	if err := TestAccProvider.InternalValidate(); err != nil {
		panic(fmt.Errorf("err: %s", err))
	}

	TestAccProviders = map[string]*schema.Provider{
		"taikun": TestAccProvider,
	}
	TestAccProviderFactories = map[string]func() (*schema.Provider, error){
		"taikun": func() (*schema.Provider, error) {
			return TestAccProvider, nil
		},
	}
}
