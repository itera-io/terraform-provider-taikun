package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/itera-io/terraform-provider-taikun/taikun"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

// Proofread the generated documentation
//go:generate ./scripts/docs_cleanup.sh

func main() {
	plugin.Serve(
		&plugin.ServeOpts{
			ProviderFunc: func() *schema.Provider {
				return taikun.Provider()
			},
		})
}
