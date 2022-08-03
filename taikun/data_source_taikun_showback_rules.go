package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/models"
	"github.com/itera-io/taikungoclient/showbackclient/showback_rules"
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
	apiClient := meta.(*apiClient)
	dataSourceID := "all"

	params := showback_rules.NewShowbackRulesListParams().WithV(ApiVersion)

	organizationIDData, organizationIDProvided := d.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.WithOrganizationID(&organizationID)
	}

	var showbackRulesList []*models.ShowbackRulesListDto
	for {
		response, err := apiClient.ShowbackClient.Showback.ShowbackRulesList(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		showbackRulesList = append(showbackRulesList, response.GetPayload().Data...)
		if len(showbackRulesList) == int(response.GetPayload().TotalCount) {
			break
		}
		offset := int32(len(showbackRulesList))
		params = params.WithOffset(&offset)
	}

	showbackRules := make([]map[string]interface{}, len(showbackRulesList))
	for i, rawShowbackRule := range showbackRulesList {
		showbackRules[i] = flattenTaikunShowbackRule(rawShowbackRule)
	}
	if err := d.Set("showback_rules", showbackRules); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
