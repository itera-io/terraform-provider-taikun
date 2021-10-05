package taikun

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceTaikunAccessProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaikunAccessProfilesRead,
		Schema:      map[string]*schema.Schema{},
	}
}
