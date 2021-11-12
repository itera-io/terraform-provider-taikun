package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/prometheus"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunBillingRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "The ID of the billing rule.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description:  "The name of the billing rule.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(3, 30),
		},
		"metric_name": {
			Description:  "The name of the Prometheus metric (e.g. volumes, flavors, networks) to bill.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(3, 256),
		},
		"label": {
			Description: "Labels linked to the billing rule.",
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    1,
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
						Description: "ID of the label.",
						Type:        schema.TypeString,
						Computed:    true,
					},
				},
			},
		},
		"type": {
			Description:  "The type of billing rule: `Count` (calculate package as unit) or `Sum` (calculate per quantity).",
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
			Description: "The creator of the billing rule.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"billing_credential_id": {
			Description:      "ID of the billing credential.",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"last_modified": {
			Description: "Time and date of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user to have modified the billing rule.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func resourceTaikunBillingRule() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Billing Rule",
		CreateContext: resourceTaikunBillingRuleCreate,
		ReadContext:   generateResourceTaikunBillingRuleReadWithoutRetries(),
		UpdateContext: resourceTaikunBillingRuleUpdate,
		DeleteContext: resourceTaikunBillingRuleDelete,
		Schema:        resourceTaikunBillingRuleSchema(),
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
		Labels:                resourceTaikunBillingRuleLabelsToAdd(data),
		Name:                  data.Get("name").(string),
		MetricName:            data.Get("metric_name").(string),
		Price:                 data.Get("price").(float64),
		OperationCredentialID: billingCredentialId,
		Type:                  getPrometheusType(data.Get("type").(string)),
	}

	params := prometheus.NewPrometheusCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.client.Prometheus.PrometheusCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(createResult.Payload.ID)

	return readAfterCreateWithRetries(generateResourceTaikunBillingRuleReadWithRetries(), ctx, data, meta)
}
func generateResourceTaikunBillingRuleReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunBillingRuleRead(true)
}
func generateResourceTaikunBillingRuleReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunBillingRuleRead(false)
}
func generateResourceTaikunBillingRuleRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		if len(response.Payload.Data) != 1 {
			if withRetries {
				data.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawBillingRule := response.GetPayload().Data[0]

		err = setResourceDataFromMap(data, flattenTaikunBillingRule(rawBillingRule))
		if err != nil {
			return diag.FromErr(err)
		}

		data.SetId(i32toa(id))

		return nil
	}
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
		LabelsToAdd:           resourceTaikunBillingRuleLabelsToAdd(data),
		LabelsToDelete:        resourceTaikunBillingRuleLabelsToDelete(data),
		Name:                  data.Get("name").(string),
		MetricName:            data.Get("metric_name").(string),
		Price:                 data.Get("price").(float64),
		OperationCredentialID: billingCredentialId,
		Type:                  getPrometheusType(data.Get("type").(string)),
	}

	params := prometheus.NewPrometheusUpdateParams().WithV(ApiVersion).WithID(id).WithBody(body)
	_, err = apiClient.client.Prometheus.PrometheusUpdate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	return readAfterUpdateWithRetries(generateResourceTaikunBillingRuleReadWithRetries(), ctx, data, meta)
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

func resourceTaikunBillingRuleLabelsToAdd(data *schema.ResourceData) []*models.PrometheusLabelListDto {
	labels := data.Get("label").([]interface{})
	labelsToAdd := make([]*models.PrometheusLabelListDto, len(labels))
	for i, labelData := range labels {
		label := labelData.(map[string]interface{})
		labelsToAdd[i] = &models.PrometheusLabelListDto{
			Label: label["key"].(string),
			Value: label["value"].(string),
		}
	}
	return labelsToAdd
}

func resourceTaikunBillingRuleLabelsToDelete(data *schema.ResourceData) []*models.PrometheusLabelDeleteDto {
	oldLabelsData, _ := data.GetChange("label")
	oldLabels := oldLabelsData.([]interface{})
	labelsToDelete := make([]*models.PrometheusLabelDeleteDto, len(oldLabels))
	for i, oldLabelData := range oldLabels {
		oldLabel := oldLabelData.(map[string]interface{})
		oldLabelID, _ := atoi32(oldLabel["id"].(string))
		labelsToDelete[i] = &models.PrometheusLabelDeleteDto{
			ID: oldLabelID,
		}
	}
	return labelsToDelete
}

func flattenTaikunBillingRule(rawBillingRule *models.PrometheusRuleListDto) map[string]interface{} {

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
