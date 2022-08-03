package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient"
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

func dataSourceTaikunBillingRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)

	params := prometheus.NewPrometheusListOfRulesParams().WithV(ApiVersion)

	var billingRulesList []*models.PrometheusRuleListDto
	for {
		response, err := apiClient.Client.Prometheus.PrometheusListOfRules(params, apiClient)
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
		billingRules[i] = flattenTaikunBillingRule(rawBillingRule)
	}
	if err := d.Set("billing_rules", billingRules); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("all")

	return nil
}
