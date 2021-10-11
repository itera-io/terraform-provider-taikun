package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/prometheus"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunBillingRule() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Billing Rule",
		CreateContext: resourceTaikunBillingRuleCreate,
		ReadContext:   resourceTaikunBillingRuleRead,
		UpdateContext: resourceTaikunBillingRuleUpdate,
		DeleteContext: resourceTaikunBillingRuleDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The id of the billing rule.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the billing rule.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"metric_name": {
				Description: "The name of the metric from Prometheus you want to bill.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"label": {
				Description: "Labels linked to this billing rule.",
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Description: "Key of the label.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"value": {
							Description: "Value of the label.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"id": {
							Description: "Id of the label.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"type": {
				Description:  "Type of the billing rule. `Count` (calculate package as unit) or `Sum` (calculate per quantity)",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Count", "Sum"}, false),
			},
			"price": {
				Description:  "The price in CZK per selected unit.",
				Type:         schema.TypeFloat,
				Required:     true,
				ValidateFunc: validation.FloatAtLeast(0),
			},
			"created_by": {
				Description: "The creator of the billing credential.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"billing_credential_id": {
				Description:  "Id of the billing credential.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: stringIsInt,
			},
			"last_modified": {
				Description: "Time of last modification.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_modified_by": {
				Description: "The last user who modified the billing credential.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunBillingRuleCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	billingCredentialId, err := atoi32(data.Get("billing_credential_id").(string))
	if err != nil {
		return diag.Errorf("billing_credential_id isn't valid: %s", data.Get("billing_credential_id").(string))
	}

	body := &models.RuleCreateCommand{
		Name:                  data.Get("name").(string),
		MetricName:            data.Get("metric_name").(string),
		Price:                 data.Get("price").(float64),
		OperationCredentialID: billingCredentialId,
		Type:                  models.PrometheusType(getPrometheusType(data.Get("type").(string))),
	}

	rawLabelsList := data.Get("label").([]interface{})
	LabelsList := make([]*models.PrometheusLabelListDto, len(rawLabelsList), len(rawLabelsList))
	for i, e := range rawLabelsList {
		rawLabel := e.(map[string]interface{})
		LabelsList[i] = &models.PrometheusLabelListDto{
			Label: rawLabel["key"].(string),
			Value: rawLabel["value"].(string),
		}
	}
	body.Labels = LabelsList

	params := prometheus.NewPrometheusCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.client.Prometheus.PrometheusCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(createResult.Payload.ID)

	return resourceTaikunBillingRuleRead(ctx, data, meta)
}

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
				"key":   rawLabel.Label,
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
		if err := data.Set("label", labels); err != nil {
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

func resourceTaikunBillingRuleUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	billingCredentialId, err := atoi32(data.Get("billing_credential_id").(string))
	if err != nil {
		return diag.Errorf("billing_credential_id isn't valid: %s", data.Get("billing_credential_id").(string))
	}

	body := &models.RuleForUpdateDto{
		Name:                  data.Get("name").(string),
		MetricName:            data.Get("metric_name").(string),
		Price:                 data.Get("price").(float64),
		OperationCredentialID: billingCredentialId,
		Type:                  models.PrometheusType(getPrometheusType(data.Get("type").(string))),
	}

	params := prometheus.NewPrometheusUpdateParams().WithV(ApiVersion).WithID(id).WithBody(body)
	_, err = apiClient.client.Prometheus.PrometheusUpdate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceTaikunBillingRuleRead(ctx, data, meta)
}

func resourceTaikunBillingRuleDelete(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	params := prometheus.NewPrometheusDeleteParams().WithV(ApiVersion).WithID(id)
	_, err = apiClient.client.Prometheus.PrometheusDelete(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}
