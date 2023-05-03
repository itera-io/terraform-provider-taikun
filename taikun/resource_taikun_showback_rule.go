package taikun

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/models"
	"github.com/itera-io/taikungoclient/showbackclient/showback_rules"
)

func resourceTaikunShowbackRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"created_by": {
			Description: "The creator of the showback rule.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"global_alert_limit": {
			Description:  "Set limit of alerts for all projects.",
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntAtLeast(0),
		},
		"id": {
			Description: "The ID of the showback rule.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"kind": {
			Description:  "The kind of showback rule: `General` (data source is Taikun) or `External` (data source is external, see `showback_credential_id`).",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"General", "External"}, false),
		},
		"label": {
			Description: "Labels linked to this showback rule.",
			Type:        schema.TypeSet,
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
		"metric_name": {
			Description:  "The metric name.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(3, 256),
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
		"showback_credential_id": {
			Description:      "ID of the showback credential.",
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"showback_credential_name": {
			Description: "Name of the showback credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"type": {
			Description:  "The type of showback rule: `Count` (calculate package as unit) or `Sum` (calculate per quantity).",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"Count", "Sum"}, false),
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

func resourceTaikunShowbackRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)
	// temporary hack to fix
	apiClient.Refresh()

	body := &models.CreateShowbackRuleCommand{
		Name:              d.Get("name").(string),
		MetricName:        d.Get("metric_name").(string),
		Type:              getEPrometheusType(d.Get("type").(string)),
		Kind:              getShowbackType(d.Get("kind").(string)),
		Price:             d.Get("price").(float64),
		ProjectAlertLimit: int32(d.Get("project_alert_limit").(int)),
		GlobalAlertLimit:  int32(d.Get("global_alert_limit").(int)),
	}

	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
		body.OrganizationID = organizationId
	}

	showbackCredentialIDData, showbackCredentialIDIsSet := d.GetOk("showback_credential_id")
	if showbackCredentialIDIsSet {
		showbackCredentialId, err := atoi32(showbackCredentialIDData.(string))
		if err != nil {
			return diag.Errorf("showback_credential_id isn't valid: %s", d.Get("showback_credential_id").(string))
		}
		body.ShowbackCredentialID = &showbackCredentialId
	}

	rawLabelsList := d.Get("label").(*schema.Set).List()
	LabelsList := make([]*models.ShowbackLabelCreateDto, len(rawLabelsList))
	for i, e := range rawLabelsList {
		rawLabel := e.(map[string]interface{})
		LabelsList[i] = &models.ShowbackLabelCreateDto{
			Label: rawLabel["key"].(string),
			Value: rawLabel["value"].(string),
		}
	}
	body.Labels = LabelsList

	params := showback_rules.NewShowbackRulesCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.ShowbackClient.ShowbackRules.ShowbackRulesCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.Payload.ID)

	return readAfterCreateWithRetries(generateResourceTaikunShowbackRuleReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunShowbackRuleReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunShowbackRuleRead(true)
}
func generateResourceTaikunShowbackRuleReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunShowbackRuleRead(false)
}
func generateResourceTaikunShowbackRuleRead(withRetries bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*taikungoclient.Client)
		id, err := atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, err := apiClient.ShowbackClient.ShowbackRules.ShowbackRulesList(showback_rules.NewShowbackRulesListParams().WithV(ApiVersion).WithID(&id), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(response.Payload.Data) != 1 {
			if withRetries {
				d.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawShowbackRule := response.GetPayload().Data[0]

		err = setResourceDataFromMap(d, flattenTaikunShowbackRule(rawShowbackRule))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunShowbackRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	body := &models.UpdateShowbackRuleCommand{
		ID:                id,
		Name:              d.Get("name").(string),
		MetricName:        d.Get("metric_name").(string),
		Type:              getEPrometheusType(d.Get("type").(string)),
		Kind:              getShowbackType(d.Get("kind").(string)),
		Price:             d.Get("price").(float64),
		ProjectAlertLimit: int32(d.Get("project_alert_limit").(int)),
		GlobalAlertLimit:  int32(d.Get("global_alert_limit").(int)),
	}

	rawLabelsList := d.Get("label").(*schema.Set).List()
	LabelsList := make([]*models.ShowbackLabelCreateDto, len(rawLabelsList))
	for i, e := range rawLabelsList {
		rawLabel := e.(map[string]interface{})
		LabelsList[i] = &models.ShowbackLabelCreateDto{
			Label: rawLabel["key"].(string),
			Value: rawLabel["value"].(string),
		}
	}
	body.Labels = LabelsList

	params := showback_rules.NewShowbackRulesUpdateParams().WithV(ApiVersion).WithBody(body)
	_, err = apiClient.ShowbackClient.ShowbackRules.ShowbackRulesUpdate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	return readAfterUpdateWithRetries(generateResourceTaikunShowbackRuleReadWithRetries(), ctx, d, meta)
}

func resourceTaikunShowbackRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	params := showback_rules.NewShowbackRulesDeleteParams().WithV(ApiVersion).WithID(id)
	_, _, err = apiClient.ShowbackClient.ShowbackRules.ShowbackRulesDelete(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
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
