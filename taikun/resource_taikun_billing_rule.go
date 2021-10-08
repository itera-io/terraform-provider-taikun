package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/prometheus"
)

func resourceTaikunBillingRuleRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	data.SetId("")
	if err != nil {
		return diag.FromErr(err)
	}

	params := prometheus.NewPrometheusListOfRulesParams().WithV(ApiVersion).WithID(&id)
	response, err := apiClient.client.Prometheus.PrometheusListOfRules(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	if response.Payload.TotalCount == 1 {
		rawBillingRule := response.GetPayload().Data[0]

		labels := make([]map[string]interface{}, len(rawBillingRule.Labels), len(rawBillingRule.Labels))
		for i, rawLabel := range rawBillingRule.Labels {
			labels[i] = map[string]interface{}{
				"label": rawLabel.Label,
				"value": rawLabel.Value,
				"id":    i32toa(rawLabel.ID),
			}
		}

		if err := data.Set("billing_credential_id", i32toa(rawBillingRule.OperationCredential.OperationCredentialID)); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("created_by", rawBillingRule.CreatedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("id", i32toa(rawBillingRule.ID)); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("labels", labels); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("last_modified", rawBillingRule.LastModified); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("last_modified_by", rawBillingRule.LastModifiedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("metric_name", rawBillingRule.MetricName); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("name", rawBillingRule.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("price", rawBillingRule.Price); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("type", rawBillingRule.Type); err != nil {
			return diag.FromErr(err)
		}

		data.SetId(i32toa(id))
	}

	return nil
}
