package showback

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkshowback "github.com/itera-io/taikungoclient/showbackclient"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
			ValidateDiagFunc: utils.StringIsInt,
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
			ValidateDiagFunc: utils.StringIsInt,
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

func ResourceTaikunShowbackRule() *schema.Resource {
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
	apiClient := meta.(*tk.Client)

	body := tkshowback.CreateShowbackRuleCommand{}
	body.SetName(d.Get("name").(string))
	body.SetMetricName(d.Get("metric_name").(string))
	body.SetType(utils.GetEPrometheusType(d.Get("type").(string)))
	body.SetKind(utils.GetShowbackType(d.Get("kind").(string)))
	body.SetPrice(d.Get("price").(float64))
	body.SetProjectAlertLimit(int32(d.Get("project_alert_limit").(int)))
	body.SetGlobalAlertLimit(int32(d.Get("global_alert_limit").(int)))

	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := utils.Atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
		body.SetOrganizationId(organizationId)
	}

	showbackCredentialIDData, showbackCredentialIDIsSet := d.GetOk("showback_credential_id")
	if showbackCredentialIDIsSet {
		showbackCredentialId, err := utils.Atoi32(showbackCredentialIDData.(string))
		if err != nil {
			return diag.Errorf("showback_credential_id isn't valid: %s", d.Get("showback_credential_id").(string))
		}
		body.SetShowbackCredentialId(showbackCredentialId)
	}

	rawLabelsList := d.Get("label").(*schema.Set).List()
	LabelsList := make([]tkshowback.ShowbackLabelCreateDto, len(rawLabelsList))
	for i, e := range rawLabelsList {
		rawLabel := e.(map[string]interface{})
		LabelsList[i] = tkshowback.ShowbackLabelCreateDto{}
		LabelsList[i].SetLabel(rawLabel["key"].(string))
		LabelsList[i].SetValue(rawLabel["value"].(string))
	}
	body.Labels = LabelsList

	createResult, resp, err := apiClient.ShowbackClient.ShowbackRulesAPI.ShowbackrulesCreate(context.TODO()).CreateShowbackRuleCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(resp, err))
	}

	d.SetId(createResult.GetId())

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunShowbackRuleReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunShowbackRuleReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunShowbackRuleRead(true)
}
func generateResourceTaikunShowbackRuleReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunShowbackRuleRead(false)
}
func generateResourceTaikunShowbackRuleRead(withRetries bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)
		id, err := utils.Atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, resp, err := apiClient.ShowbackClient.ShowbackRulesAPI.ShowbackrulesList(context.TODO()).Id(id).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(resp, err))
		}
		if len(response.GetData()) != 1 {
			if withRetries {
				d.SetId(utils.I32toa(id))
				return diag.Errorf(utils.NotFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawShowbackRule := response.GetData()[0]

		err = utils.SetResourceDataFromMap(d, flattenTaikunShowbackRule(&rawShowbackRule))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.I32toa(id))

		return nil
	}
}

func resourceTaikunShowbackRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	id, err := utils.Atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	body := tkshowback.UpdateShowbackRuleCommand{}
	body.SetId(id)
	body.SetName(d.Get("name").(string))
	body.SetMetricName(d.Get("metric_name").(string))
	body.SetType(utils.GetEPrometheusType(d.Get("type").(string)))
	body.SetKind(utils.GetShowbackType(d.Get("kind").(string)))
	body.SetPrice(d.Get("price").(float64))
	body.SetProjectAlertLimit(int32(d.Get("project_alert_limit").(int)))
	body.SetGlobalAlertLimit(int32(d.Get("global_alert_limit").(int)))

	rawLabelsList := d.Get("label").(*schema.Set).List()
	LabelsList := make([]tkshowback.ShowbackLabelCreateDto, len(rawLabelsList))
	for i, e := range rawLabelsList {
		rawLabel := e.(map[string]interface{})
		LabelsList[i] = tkshowback.ShowbackLabelCreateDto{}
		LabelsList[i].SetLabel(rawLabel["key"].(string))
		LabelsList[i].SetValue(rawLabel["value"].(string))
	}
	body.Labels = LabelsList

	resp, err := apiClient.ShowbackClient.ShowbackRulesAPI.ShowbackrulesUpdate(context.TODO()).UpdateShowbackRuleCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(resp, err))
	}

	return utils.ReadAfterUpdateWithRetries(generateResourceTaikunShowbackRuleReadWithRetries(), ctx, d, meta)
}

func resourceTaikunShowbackRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	id, err := utils.Atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := apiClient.ShowbackClient.ShowbackRulesAPI.ShowbackrulesDelete(context.TODO(), id).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(resp, err))
	}

	d.SetId("")
	return nil
}

func flattenTaikunShowbackRule(rawShowbackRule *tkshowback.ShowbackRulesListDto) map[string]interface{} {

	labels := make([]map[string]interface{}, len(rawShowbackRule.GetLabels()))
	for i, rawLabel := range rawShowbackRule.GetLabels() {
		labels[i] = map[string]interface{}{
			"key":   rawLabel.GetLabel(),
			"value": rawLabel.GetValue(),
		}
	}

	result := map[string]interface{}{
		"created_by":          rawShowbackRule.GetCreatedBy(),
		"global_alert_limit":  rawShowbackRule.GetGlobalAlertLimit(),
		"id":                  utils.I32toa(rawShowbackRule.GetId()),
		"kind":                rawShowbackRule.GetKind(),
		"label":               labels,
		"last_modified":       rawShowbackRule.GetLastModified(),
		"last_modified_by":    rawShowbackRule.GetLastModifiedBy(),
		"metric_name":         rawShowbackRule.GetMetricName(),
		"name":                rawShowbackRule.GetName(),
		"organization_id":     utils.I32toa(rawShowbackRule.GetOrganizationId()),
		"organization_name":   rawShowbackRule.GetOrganizationName(),
		"price":               rawShowbackRule.GetPrice(),
		"project_alert_limit": rawShowbackRule.GetProjectAlertLimit(),
		"type":                rawShowbackRule.GetType(),
	}

	if _, ok := rawShowbackRule.GetShowbackCredentialIdOk(); ok {
		// It seems there was a slight change in the API. Now it returns id 0 and name "" for empty credential.
		// Set this only if it is actually set.
		if rawShowbackRule.GetShowbackCredentialName() != "" {
			result["showback_credential_id"] = utils.I32toa(rawShowbackRule.GetShowbackCredentialId())
			result["showback_credential_name"] = rawShowbackRule.GetShowbackCredentialName()
		}
	}
	return result
}
