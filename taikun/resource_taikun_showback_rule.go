package taikun

import (
	"context"

	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/showback"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunShowbackRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "The ID of the showback rule.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description: "The name of the showback rule.",
			Type:        schema.TypeString,
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 30),
				validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z0-9-_.]+$"),
					"expected only alpha numeric characters or non alpha numeric (_-.)",
				),
			),
		},
		"metric_name": {
			Description:  "The metric name.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(3, 256),
		},
		"kind": {
			Description:  "The kind of showback rule: `General` (data source is Taikun) or `External` (data source is external, see `showback_credential_id`).",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"General", "External"}, false),
		},
		"type": {
			Description:  "The type of showback rule: `Count` (calculate package as unit) or `Sum` (calculate per quantity).",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"Count", "Sum"}, false),
		},
		"price": {
			Description:  "Billing in CZK per selected unit.",
			Type:         schema.TypeFloat,
			Required:     true,
			ValidateFunc: validation.FloatAtLeast(0),
		},
		"project_alert_limit": {
			Description:  "Set limit of alerts for one project.",
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntAtLeast(0),
		},
		"global_alert_limit": {
			Description:  "Set limit of alerts for all projects.",
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntAtLeast(0),
		},
		"organization_id": {
			Description:      "The ID of the organization which owns the showback rule.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: stringIsInt,
			ForceNew:         true,
		},
		"organization_name": {
			Description: "The name of the organization which owns the showback rule.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"showback_credential_id": {
			Description:      "Id of the showback rule.",
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"showback_credential_name": {
			Description: "Name of the showback rule.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"created_by": {
			Description: "The creator of the showback rule.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified": {
			Description: "Time of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user who modified the showback rule.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"label": {
			Description: "Labels linked to this showback rule.",
			Type:        schema.TypeList,
			Optional:    true,
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
				},
			},
		},
	}
}

func resourceTaikunShowbackRule() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Showback Rule",
		CreateContext: resourceTaikunShowbackRuleCreate,
		ReadContext:   generateResourceTaikunShowbackRuleReadWithoutRetries(),
		UpdateContext: resourceTaikunShowbackRuleUpdate,
		DeleteContext: resourceTaikunShowbackRuleDelete,
		Schema:        resourceTaikunShowbackRuleSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunShowbackRuleCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := &models.CreateShowbackRuleCommand{
		Name:              data.Get("name").(string),
		MetricName:        data.Get("metric_name").(string),
		Type:              getPrometheusType(data.Get("type").(string)),
		Kind:              getShowbackType(data.Get("kind").(string)),
		Price:             data.Get("price").(float64),
		ProjectAlertLimit: int32(data.Get("project_alert_limit").(int)),
		GlobalAlertLimit:  int32(data.Get("global_alert_limit").(int)),
	}

	organizationIDData, organizationIDIsSet := data.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", data.Get("organization_id").(string))
		}
		body.OrganizationID = organizationId
	}

	showbackCredentialIDData, showbackCredentialIDIsSet := data.GetOk("showback_credential_id")
	if showbackCredentialIDIsSet {
		showbackCredentialId, err := atoi32(showbackCredentialIDData.(string))
		if err != nil {
			return diag.Errorf("showback_credential_id isn't valid: %s", data.Get("showback_credential_id").(string))
		}
		body.ShowbackCredentialID = &showbackCredentialId
	}

	rawLabelsList := data.Get("label").([]interface{})
	LabelsList := make([]*models.ShowbackLabelCreateDto, len(rawLabelsList))
	for i, e := range rawLabelsList {
		rawLabel := e.(map[string]interface{})
		LabelsList[i] = &models.ShowbackLabelCreateDto{
			Label: rawLabel["key"].(string),
			Value: rawLabel["value"].(string),
		}
	}
	body.Labels = LabelsList

	params := showback.NewShowbackCreateRuleParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.client.Showback.ShowbackCreateRule(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(createResult.Payload.ID)

	return readAfterCreateWithRetries(generateResourceTaikunShowbackRuleReadWithRetries(), ctx, data, meta)
}
func generateResourceTaikunShowbackRuleReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunShowbackRuleRead(true)
}
func generateResourceTaikunShowbackRuleReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunShowbackRuleRead(false)
}
func generateResourceTaikunShowbackRuleRead(withRetries bool) schema.ReadContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*apiClient)
		id, err := atoi32(data.Id())
		data.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, err := apiClient.client.Showback.ShowbackRulesList(showback.NewShowbackRulesListParams().WithV(ApiVersion).WithID(&id), apiClient)
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

		rawShowbackRule := response.GetPayload().Data[0]

		err = setResourceDataFromMap(data, flattenTaikunShowbackRule(rawShowbackRule))
		if err != nil {
			return diag.FromErr(err)
		}

		data.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunShowbackRuleUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	body := &models.UpdateShowbackRuleCommand{
		ID:                id,
		Name:              data.Get("name").(string),
		MetricName:        data.Get("metric_name").(string),
		Type:              getPrometheusType(data.Get("type").(string)),
		Kind:              getShowbackType(data.Get("kind").(string)),
		Price:             data.Get("price").(float64),
		ProjectAlertLimit: int32(data.Get("project_alert_limit").(int)),
		GlobalAlertLimit:  int32(data.Get("global_alert_limit").(int)),
	}

	rawLabelsList := data.Get("label").([]interface{})
	LabelsList := make([]*models.ShowbackLabelCreateDto, len(rawLabelsList))
	for i, e := range rawLabelsList {
		rawLabel := e.(map[string]interface{})
		LabelsList[i] = &models.ShowbackLabelCreateDto{
			Label: rawLabel["key"].(string),
			Value: rawLabel["value"].(string),
		}
	}
	body.Labels = LabelsList

	params := showback.NewShowbackUpdateRuleParams().WithV(ApiVersion).WithBody(body)
	_, err = apiClient.client.Showback.ShowbackUpdateRule(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	return readAfterUpdateWithRetries(generateResourceTaikunShowbackRuleReadWithRetries(), ctx, data, meta)
}

func resourceTaikunShowbackRuleDelete(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	params := showback.NewShowbackDeleteRuleParams().WithV(ApiVersion).WithBody(&models.DeleteShowbackRuleCommand{ID: id})
	_, err = apiClient.client.Showback.ShowbackDeleteRule(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

func flattenTaikunShowbackRule(rawShowbackRule *models.ShowbackRulesListDto) map[string]interface{} {

	labels := make([]map[string]interface{}, len(rawShowbackRule.Labels))
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

	if rawShowbackRule.ShowbackCredentialID != nil {
		result["showback_credential_id"] = i32toa(*rawShowbackRule.ShowbackCredentialID)
		result["showback_credential_name"] = rawShowbackRule.ShowbackCredentialName
	}
	return result
}
