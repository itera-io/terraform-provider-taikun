package taikun

import (
	"context"
	tk "github.com/chnyda/taikungoclient"
	tkshowback "github.com/chnyda/taikungoclient/showbackclient"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunShowbackRules() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all showback rules.",
		ReadContext: dataSourceTaikunShowbackRulesRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsInt,
			},
			"showback_rules": {
				Description: "List of retrieved showback rules.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunShowbackRuleSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunShowbackRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"

	params := apiClient.ShowbackClient.ShowbackRulesApi.ShowbackrulesList(context.TODO())
	var offset int32 = 0

	organizationIDData, organizationIDProvided := d.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.OrganizationId(organizationID)
	}

	var showbackRulesList []tkshowback.ShowbackRulesListDto
	for {
		response, resp, err := params.Offset(offset).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(resp, err))
		}
		showbackRulesList = append(showbackRulesList, response.GetData()...)
		if len(showbackRulesList) == int(response.GetTotalCount()) {
			break
		}
		offset = int32(len(showbackRulesList))
	}

	showbackRules := make([]map[string]interface{}, len(showbackRulesList))
	for i, rawShowbackRule := range showbackRulesList {
		showbackRules[i] = flattenTaikunShowbackRule(&rawShowbackRule)
	}
	if err := d.Set("showback_rules", showbackRules); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
