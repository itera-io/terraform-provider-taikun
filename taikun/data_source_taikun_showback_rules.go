package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/showback"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunShowbackRules() *schema.Resource {
	return &schema.Resource{
		Description: "Get the list of showback rules, optionally filtered by organization.",
		ReadContext: dataSourceTaikunShowbackRulesRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:  "Organization id filter.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: stringIsInt,
			},
			"showback_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunShowbackRuleSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunShowbackRulesRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	dataSourceID := "all"

	params := showback.NewShowbackRulesListParams().WithV(ApiVersion)

	organizationIDData, organizationIDProvided := data.GetOk("organization_id")
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
		response, err := apiClient.client.Showback.ShowbackRulesList(params, apiClient)
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

	showbackRules := make([]map[string]interface{}, len(showbackRulesList), len(showbackRulesList))
	for i, rawShowbackRule := range showbackRulesList {
		showbackRules[i] = flattenDatasourceTaikunShowbackRuleItem(rawShowbackRule)
	}
	if err := data.Set("showback_rules", showbackRules); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(dataSourceID)

	return nil
}

func flattenDatasourceTaikunShowbackRuleItem(rawShowbackRule *models.ShowbackRulesListDto) map[string]interface{} {

	labels := make([]map[string]interface{}, len(rawShowbackRule.Labels), len(rawShowbackRule.Labels))
	for i, rawLabel := range rawShowbackRule.Labels {
		labels[i] = map[string]interface{}{
			"key":   rawLabel.Label,
			"value": rawLabel.Value,
		}
	}

	result := map[string]interface{}{
		"created_by":          rawShowbackRule.CreatedBy,
		"global_alert_limit":  rawShowbackRule.GlobalAlertLimit,
		"id":                  i32toa(rawShowbackRule.ID),
		"kind":                rawShowbackRule.Kind,
		"label":               labels,
		"last_modified":       rawShowbackRule.LastModified,
		"last_modified_by":    rawShowbackRule.LastModifiedBy,
		"metric_name":         rawShowbackRule.MetricName,
		"name":                rawShowbackRule.Name,
		"organization_id":     i32toa(rawShowbackRule.OrganizationID),
		"organization_name":   rawShowbackRule.OrganizationName,
		"price":               rawShowbackRule.Price,
		"project_alert_limit": rawShowbackRule.ProjectAlertLimit,
		"type":                rawShowbackRule.Type,
	}

	if rawShowbackRule.ShowbackCredentialID != 0 {
		result["showback_credential_id"] = rawShowbackRule.ShowbackCredentialID
		result["showback_credential_name"] = rawShowbackRule.ShowbackCredentialName
	}
	return result
}
