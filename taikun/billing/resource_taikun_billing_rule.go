package billing

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

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
			ValidateDiagFunc: utils.StringIsInt,
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
			Set:         utils.HashAttributes("key", "value"),
			MinItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					//"id": {
					//	Description: "ID of the label.",
					//	Type:        schema.TypeString,
					//	Computed:    true,
					//},
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

func ResourceTaikunBillingRule() *schema.Resource {
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

	billingCredentialId, err := utils.Atoi32(d.Get("billing_credential_id").(string))
	if err != nil {
		return diag.Errorf("billing_credential_id isn't valid: %s", d.Get("billing_credential_id").(string))
	}

	body := tkcore.RuleCreateCommand{}
	body.SetLabels(resourceTaikunBillingRuleLabelsAll(d))
	body.SetName(d.Get("name").(string))
	body.SetMetricName(d.Get("metric_name").(string))
	body.SetPrice(d.Get("price").(float64))
	body.SetOperationCredentialId(billingCredentialId)
	body.SetType(utils.GetPrometheusType(d.Get("type").(string)))

	createResult, res, err := apiClient.Client.PrometheusRulesAPI.PrometheusrulesCreate(ctx).RuleCreateCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId(createResult.GetId())

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunBillingRuleReadWithRetries(), ctx, d, meta)
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
		id, err := utils.Atoi32(d.Id())
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
				d.SetId(utils.I32toa(id))
				return diag.Errorf(utils.NotFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawBillingRule := response.Data[0]

		err = utils.SetResourceDataFromMap(d, flattenTaikunBillingRule(&rawBillingRule))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(utils.I32toa(id))

		return nil
	}
}

func resourceTaikunBillingRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	id, err := utils.Atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	billingCredentialId, err := utils.Atoi32(d.Get("billing_credential_id").(string))
	if err != nil {
		return diag.Errorf("billing_credential_id isn't valid: %s", d.Get("billing_credential_id").(string))
	}

	body := tkcore.RuleForUpdateDto{}
	//body.SetLabelsToAdd(resourceTaikunBillingRuleLabelsToAdd(d))
	//body.SetLabelsToDelete(resourceTaikunBillingRuleLabelsToDelete(d))
	//_, newLabelsData := d.GetChange("label")
	body.SetLabels(resourceTaikunBillingRuleLabelsAll(d))
	body.SetName(d.Get("name").(string))
	body.SetMetricName(d.Get("metric_name").(string))
	body.SetPrice(d.Get("price").(float64))
	body.SetOperationCredentialId(billingCredentialId)
	body.SetType(utils.GetPrometheusType(d.Get("type").(string)))

	_, res, err := apiClient.Client.PrometheusRulesAPI.PrometheusrulesUpdate(context.TODO(), id).RuleForUpdateDto(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	return utils.ReadAfterUpdateWithRetries(generateResourceTaikunBillingRuleReadWithRetries(), ctx, d, meta)
}

func resourceTaikunBillingRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	id, err := utils.Atoi32(d.Id())
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

func resourceTaikunBillingRuleLabelsAll(d *schema.ResourceData) []tkcore.PrometheusLabelListDto {
	_, newLabelsData := d.GetChange("label")
	newLabels := newLabelsData.(*schema.Set).List()
	labelsToAdd := make([]tkcore.PrometheusLabelListDto, len(newLabels))
	for i, labelData := range newLabels {
		label := labelData.(map[string]interface{})
		labelsToAdd[i] = tkcore.PrometheusLabelListDto{}
		labelsToAdd[i].SetLabel(label["key"].(string))
		labelsToAdd[i].SetValue(label["value"].(string))
	}
	return labelsToAdd
}

//func resourceTaikunBillingRuleLabelsToDelete(d *schema.ResourceData) []tkcore.PrometheusLabelDeleteDto {
//	oldLabelsData, newLabelsData := d.GetChange("label")
//	oldLabels := oldLabelsData.(*schema.Set)
//	newLabels := newLabelsData.(*schema.Set)
//	toDelete := oldLabels.Difference(newLabels).List()
//	labelsToDelete := make([]tkcore.PrometheusLabelDeleteDto, len(toDelete))
//	for i, oldLabelData := range toDelete {
//		oldLabel := oldLabelData.(map[string]interface{})
//		oldLabelID, _ := utils.Atoi32(oldLabel["id"].(string))
//		labelsToDelete[i] = tkcore.PrometheusLabelDeleteDto{}
//		labelsToDelete[i].SetId(oldLabelID)
//	}
//	return labelsToDelete
//}

func flattenTaikunBillingRule(rawBillingRule *tkcore.PrometheusRuleListDto) map[string]interface{} {

	labels := make([]map[string]interface{}, len(rawBillingRule.GetLabels()))
	for i, rawLabel := range rawBillingRule.GetLabels() {
		labels[i] = map[string]interface{}{
			"key":   rawLabel.GetLabel(),
			"value": rawLabel.GetValue(),
			//"id":    utils.I32toa(rawLabel.GetId()),
			//"id": utils.I32toa(hashStrings(rawLabel.GetLabel(), rawLabel.GetValue())),
		}
	}

	return map[string]interface{}{
		"billing_credential_id": utils.I32toa(rawBillingRule.OperationCredential.GetOperationCredentialId()),
		"created_by":            rawBillingRule.GetCreatedBy(),
		"id":                    utils.I32toa(rawBillingRule.GetId()),
		"label":                 labels,
		"last_modified":         rawBillingRule.GetLastModified(),
		"last_modified_by":      rawBillingRule.GetLastModifiedBy(),
		"name":                  rawBillingRule.GetName(),
		"metric_name":           rawBillingRule.GetMetricName(),
		"price":                 rawBillingRule.GetPrice(),
		"type":                  rawBillingRule.GetType(),
	}
}

//func hashStrings(string1 string, string2 string) int32 {
//	// Concatenate the strings
//	combinedString := string1 + string2
//	// Hash the combined string using SHA-256
//	hashedString := sha256.Sum256([]byte(combinedString))
//	// Convert first 4 bytes of the hash to a unique int32
//	return int32(binaryToInt32(hashedString[:4]))
//}
//
//func binaryToInt32(bytes []byte) int32 {
//	result := int32(0)
//	for i, b := range bytes {
//		result |= int32(b) << (uint(i) * 8)
//	}
//	return result
//}
