package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/prometheus"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunBillingRules() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all billing rules.",
		ReadContext: dataSourceTaikunBillingRulesRead,
		Schema: map[string]*schema.Schema{
			"billing_rules": {
				Description: "List of retrieved billing rules.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunBillingRuleSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunBillingRulesRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	params := prometheus.NewPrometheusListOfRulesParams().WithV(ApiVersion)

	var billingRulesList []*models.PrometheusRuleListDto
	for {
		response, err := apiClient.client.Prometheus.PrometheusListOfRules(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		billingRulesList = append(billingRulesList, response.GetPayload().Data...)
		if len(billingRulesList) == int(response.GetPayload().TotalCount) {
			break
		}
		offset := int32(len(billingRulesList))
		params = params.WithOffset(&offset)
	}

	billingRules := make([]map[string]interface{}, len(billingRulesList))
	for i, rawBillingRule := range billingRulesList {
		billingRules[i] = flattenDataSourceTaikunBillingRuleItem(rawBillingRule)
	}
	if err := data.Set("billing_rules", billingRules); err != nil {
		return diag.FromErr(err)
	}

	data.SetId("all")

	return nil
}

func flattenDataSourceTaikunBillingRuleItem(rawBillingRule *models.PrometheusRuleListDto) map[string]interface{} {

	labels := make([]map[string]interface{}, len(rawBillingRule.Labels))
	for i, rawLabel := range rawBillingRule.Labels {
		labels[i] = map[string]interface{}{
			"key":   rawLabel.Label,
			"value": rawLabel.Value,
			"id":    i32toa(rawLabel.ID),
		}
	}

	return map[string]interface{}{
		"billing_credential_id": i32toa(rawBillingRule.OperationCredential.OperationCredentialID),
		"created_by":            rawBillingRule.CreatedBy,
		"id":                    i32toa(rawBillingRule.ID),
		"label":                 labels,
		"last_modified":         rawBillingRule.LastModified,
		"last_modified_by":      rawBillingRule.LastModifiedBy,
		"name":                  rawBillingRule.Name,
		"metric_name":           rawBillingRule.MetricName,
		"price":                 rawBillingRule.Price,
		"type":                  rawBillingRule.Type,
	}
}
