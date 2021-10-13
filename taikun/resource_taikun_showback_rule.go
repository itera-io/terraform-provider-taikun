package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/showback"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunShowbackRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "The id of the showback rule.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description: "The name of the showback rule.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"metric_name": {
			Description: "The metric name.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"kind": {
			Description:  "Type of the showback rule. `General` (data source is taikun) or `External` (data source is external see `showback_credential_id`)",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"General", "External"}, false),
		},
		"type": {
			Description:  "Type of the showback rule. `Count` (calculate package as unit) or `Sum` (calculate per quantity)",
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
			Description:  "The id of the organization which owns the showback rule.",
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: stringIsInt,
			ForceNew:     true,
		},
		"organization_name": {
			Description: "The name of the organization which owns the showback rule.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"showback_credential_id": {
			Description:  "Id of the showback rule.",
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: stringIsInt,
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
		ReadContext:   resourceTaikunShowbackRuleRead,
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
			return diag.Errorf("organization_id isn't valid: %s", data.Get("organization_id").(string))
		}
		body.ShowbackCredentialID = showbackCredentialId
	}

	rawLabelsList := data.Get("label").([]interface{})
	LabelsList := make([]*models.ShowbackLabelCreateDto, len(rawLabelsList), len(rawLabelsList))
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

	return resourceTaikunShowbackRuleRead(ctx, data, meta)
}

func resourceTaikunShowbackRuleRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if response.Payload.TotalCount == 1 {
		rawShowbackCredential := response.GetPayload().Data[0]

		labels := make([]map[string]interface{}, len(rawShowbackCredential.Labels), len(rawShowbackCredential.Labels))
		for i, rawLabel := range rawShowbackCredential.Labels {
			labels[i] = map[string]interface{}{
				"key":   rawLabel.Label,
				"value": rawLabel.Value,
			}
		}

		if err := data.Set("created_by", rawShowbackCredential.CreatedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("global_alert_limit", rawShowbackCredential.GlobalAlertLimit); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("id", i32toa(rawShowbackCredential.ID)); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("kind", rawShowbackCredential.Kind); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("label", labels); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("last_modified", rawShowbackCredential.LastModified); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("last_modified_by", rawShowbackCredential.LastModifiedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("metric_name", rawShowbackCredential.MetricName); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("name", rawShowbackCredential.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_id", i32toa(rawShowbackCredential.OrganizationID)); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_name", rawShowbackCredential.OrganizationName); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("price", rawShowbackCredential.Price); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("project_alert_limit", rawShowbackCredential.ProjectAlertLimit); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("type", rawShowbackCredential.Type); err != nil {
			return diag.FromErr(err)
		}

		if rawShowbackCredential.ShowbackCredentialID != 0 {
			if err := data.Set("showback_credential_id", i32toa(rawShowbackCredential.ShowbackCredentialID)); err != nil {
				return diag.FromErr(err)
			}
			if err := data.Set("showback_credential_name", rawShowbackCredential.ShowbackCredentialName); err != nil {
				return diag.FromErr(err)
			}
		}

		data.SetId(i32toa(id))
	}

	return nil
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
	LabelsList := make([]*models.ShowbackLabelCreateDto, len(rawLabelsList), len(rawLabelsList))
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

	return resourceTaikunShowbackRuleRead(ctx, data, meta)
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
