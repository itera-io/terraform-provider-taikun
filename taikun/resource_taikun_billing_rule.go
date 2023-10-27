package taikun

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTaikunBillingRuleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"billing_credential_id": {
			Description:      "ID of the billing credential.",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"created_by": {
			Description: "The creator of the billing rule.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"id": {
			Description: "The ID of the billing rule.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"label": {
			Description: "Labels linked to the billing rule.",
			Type:        schema.TypeSet,
			Required:    true,
			Set:         hashAttributes("key", "value"),
			MinItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Description: "ID of the label.",
						Type:        schema.TypeString,
						Computed:    true,
					},
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
			Description: "Time and date of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user to have modified the billing rule.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"metric_name": {
			Description:  "The name of the Prometheus metric (e.g. volumes, flavors, networks) to bill.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(3, 256),
		},
		"name": {
			Description:  "The name of the billing rule.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(3, 30),
		},
		"price": {
			Description:  "The price in CZK per selected unit.",
			Type:         schema.TypeFloat,
			Required:     true,
			ValidateFunc: validation.FloatAtLeast(0),
		},
		"type": {
			Description:  "The type of billing rule: `Count` (calculate package as unit) or `Sum` (calculate per quantity).",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"Count", "Sum"}, false),
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

func resourceTaikunBillingRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	billingCredentialId, err := atoi32(d.Get("billing_credential_id").(string))
	if err != nil {
		return diag.Errorf("billing_credential_id isn't valid: %s", d.Get("billing_credential_id").(string))
	}

	body := tkcore.RuleCreateCommand{}
	body.SetLabels(resourceTaikunBillingRuleLabelsToAdd(d))
	body.SetName(d.Get("name").(string))
	body.SetMetricName(d.Get("metric_name").(string))
	body.SetPrice(d.Get("price").(float64))
	body.SetOperationCredentialId(billingCredentialId)
	body.SetType(getPrometheusType(d.Get("type").(string)))

	createResult, res, err := apiClient.Client.PrometheusRulesAPI.PrometheusrulesCreate(ctx).RuleCreateCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId(createResult.GetId())

	return readAfterCreateWithRetries(generateResourceTaikunBillingRuleReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunBillingRuleReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunBillingRuleRead(true)
}
func generateResourceTaikunBillingRuleReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunBillingRuleRead(false)
}
func generateResourceTaikunBillingRuleRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)
		id, err := atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, res, err := apiClient.Client.PrometheusRulesAPI.PrometheusrulesList(context.TODO()).Id(id).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		if len(response.Data) != 1 {
			if withRetries {
				d.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawBillingRule := response.Data[0]

		err = setResourceDataFromMap(d, flattenTaikunBillingRule(&rawBillingRule))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunBillingRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	billingCredentialId, err := atoi32(d.Get("billing_credential_id").(string))
	if err != nil {
		return diag.Errorf("billing_credential_id isn't valid: %s", d.Get("billing_credential_id").(string))
	}

	body := tkcore.RuleForUpdateDto{}
	body.SetLabelsToAdd(resourceTaikunBillingRuleLabelsToAdd(d))
	body.SetLabelsToDelete(resourceTaikunBillingRuleLabelsToDelete(d))
	body.SetName(d.Get("name").(string))
	body.SetMetricName(d.Get("metric_name").(string))
	body.SetPrice(d.Get("price").(float64))
	body.SetOperationCredentialId(billingCredentialId)
	body.SetType(getPrometheusType(d.Get("type").(string)))

	res, err := apiClient.Client.PrometheusRulesAPI.PrometheusrulesUpdate(context.TODO(), id).RuleForUpdateDto(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	return readAfterUpdateWithRetries(generateResourceTaikunBillingRuleReadWithRetries(), ctx, d, meta)
}

func resourceTaikunBillingRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := apiClient.Client.PrometheusRulesAPI.PrometheusrulesDelete(context.TODO(), id).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId("")
	return nil
}

func resourceTaikunBillingRuleLabelsToAdd(d *schema.ResourceData) []tkcore.PrometheusLabelListDto {
	oldLabelsData, newLabelsData := d.GetChange("label")
	oldLabels := oldLabelsData.(*schema.Set)
	newLabels := newLabelsData.(*schema.Set)
	toAdd := newLabels.Difference(oldLabels).List()
	labelsToAdd := make([]tkcore.PrometheusLabelListDto, len(toAdd))
	for i, labelData := range toAdd {
		label := labelData.(map[string]interface{})
		labelsToAdd[i] = tkcore.PrometheusLabelListDto{}
		labelsToAdd[i].SetLabel(label["key"].(string))
		labelsToAdd[i].SetValue(label["value"].(string))
	}
	return labelsToAdd
}

func resourceTaikunBillingRuleLabelsToDelete(d *schema.ResourceData) []tkcore.PrometheusLabelDeleteDto {
	oldLabelsData, newLabelsData := d.GetChange("label")
	oldLabels := oldLabelsData.(*schema.Set)
	newLabels := newLabelsData.(*schema.Set)
	toDelete := oldLabels.Difference(newLabels).List()
	labelsToDelete := make([]tkcore.PrometheusLabelDeleteDto, len(toDelete))
	for i, oldLabelData := range toDelete {
		oldLabel := oldLabelData.(map[string]interface{})
		oldLabelID, _ := atoi32(oldLabel["id"].(string))
		labelsToDelete[i] = tkcore.PrometheusLabelDeleteDto{}
		labelsToDelete[i].SetId(oldLabelID)
	}
	return labelsToDelete
}

func flattenTaikunBillingRule(rawBillingRule *tkcore.PrometheusRuleListDto) map[string]interface{} {

	labels := make([]map[string]interface{}, len(rawBillingRule.GetLabels()))
	for i, rawLabel := range rawBillingRule.GetLabels() {
		labels[i] = map[string]interface{}{
			"key":   rawLabel.GetLabel(),
			"value": rawLabel.GetValue(),
			"id":    i32toa(rawLabel.GetId()),
		}
	}

	return map[string]interface{}{
		"billing_credential_id": i32toa(rawBillingRule.OperationCredential.GetOperationCredentialId()),
		"created_by":            rawBillingRule.GetCreatedBy(),
		"id":                    i32toa(rawBillingRule.GetId()),
		"label":                 labels,
		"last_modified":         rawBillingRule.GetLastModified(),
		"last_modified_by":      rawBillingRule.GetLastModifiedBy(),
		"name":                  rawBillingRule.GetName(),
		"metric_name":           rawBillingRule.GetMetricName(),
		"price":                 rawBillingRule.GetPrice(),
		"type":                  rawBillingRule.GetType(),
	}
}
