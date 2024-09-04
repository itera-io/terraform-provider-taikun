package billing

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceTaikunBillingRules() *schema.Resource {
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
	apiClient := meta.(*tk.Client)
	var offset int32 = 0

	params := apiClient.Client.PrometheusRulesAPI.PrometheusrulesList(context.TODO())

	var billingRulesList []tkcore.PrometheusRuleListDto
	for {
		response, res, err := params.Offset(offset).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		billingRulesList = append(billingRulesList, response.Data...)
		if len(billingRulesList) == int(response.GetTotalCount()) {
			break
		}
		offset = int32(len(billingRulesList))
	}

	billingRules := make([]map[string]interface{}, len(billingRulesList))
	for i, rawBillingRule := range billingRulesList {
		billingRules[i] = flattenTaikunBillingRule(&rawBillingRule)
	}
	if err := d.Set("billing_rules", billingRules); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("all")

	return nil
}
