package organization

import (
	"context"
	"fmt"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTaikunOrganizationBillingRuleAttachmentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"billing_rule_id": {
			Description:      "ID of the billing rule.",
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: utils.StringIsInt,
		},
		"discount_rate": {
			Description:  "Discount rate in percents (0-100 %).",
			Type:         schema.TypeFloat,
			Optional:     true,
			ForceNew:     true,
			Default:      100,
			ValidateFunc: validation.FloatBetween(0, 100),
		},
		"organization_id": {
			Description:      "ID of the organisation.",
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: utils.StringIsInt,
		},
		"organization_name": {
			Description: "Name of the organisation.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func ResourceTaikunOrganizationBillingRuleAttachment() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Organization - Billing Rule Attachment",
		CreateContext: resourceTaikunOrganizationBillingRuleAttachmentCreate,
		ReadContext:   generateResourceTaikunOrganizationBillingRuleAttachmentReadWithoutRetries(),
		DeleteContext: resourceTaikunOrganizationBillingRuleAttachmentDelete,
		Schema:        resourceTaikunOrganizationBillingRuleAttachmentSchema(),
	}
}

func resourceTaikunOrganizationBillingRuleAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	billingRuleId, err := utils.Atoi32(d.Get("billing_rule_id").(string))
	if err != nil {
		return diag.Errorf("billing_rule_id isn't valid: %s", d.Get("billing_rule_id").(string))
	}

	organizationId, err := utils.Atoi32(d.Get("organization_id").(string))
	if err != nil {
		return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
	}

	discountRate := d.Get("discount_rate").(float64)

	body := []tkcore.AddPrometheusRulesToOrganizationDto{
		{
			Id:           &billingRuleId,
			DiscountRate: *tkcore.NewNullableFloat64(&discountRate),
		},
	}

	_, response, err := apiClient.Client.OrganizationsAPI.OrganizationsAddPrometheusrules(context.TODO(), organizationId).AddPrometheusRulesToOrganizationDto(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(response, err))
	}

	id := fmt.Sprintf("%d/%d", organizationId, billingRuleId)
	d.SetId(id)

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunOrganizationBillingRuleAttachmentReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunOrganizationBillingRuleAttachmentReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunOrganizationBillingRuleAttachmentRead(true)
}
func generateResourceTaikunOrganizationBillingRuleAttachmentReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunOrganizationBillingRuleAttachmentRead(false)
}
func generateResourceTaikunOrganizationBillingRuleAttachmentRead(withRetries bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)

		id := d.Id()
		d.SetId("")
		organizationId, billingRuleId, err := ParseOrganizationBillingRuleAttachmentId(id)
		if err != nil {
			return diag.Errorf("Error while reading taikun_organization_billing_rule_attachment : %s", err)
		}

		response, res, err := apiClient.Client.PrometheusRulesAPI.PrometheusrulesList(ctx).Id(billingRuleId).Execute()

		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		if len(response.Data) != 1 {
			if withRetries {
				d.SetId(id)
				return diag.Errorf(utils.NotFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawBillingRule := response.Data[0]

		for _, e := range rawBillingRule.GetBoundOrganizations() {
			if e.GetId() == organizationId {
				if err := d.Set("organization_id", utils.I32toa(e.GetId())); err != nil {
					return diag.FromErr(err)
				}
				if err := d.Set("organization_name", e.GetName()); err != nil {
					return diag.FromErr(err)
				}
				if err := d.Set("billing_rule_id", utils.I32toa(rawBillingRule.GetId())); err != nil {
					return diag.FromErr(err)
				}
				if err := d.Set("discount_rate", e.GetRuleDiscountRate()); err != nil {
					return diag.FromErr(err)
				}
				d.SetId(id)
				return nil
			}
		}

		if withRetries {
			d.SetId(id)
			return diag.Errorf(utils.NotFoundAfterCreateOrUpdateError)
		}
		return nil
	}
}

func resourceTaikunOrganizationBillingRuleAttachmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	organizationId, billingRuleId, err := ParseOrganizationBillingRuleAttachmentId(d.Id())
	if err != nil {
		return diag.Errorf("Error while deleting taikun_organization_billing_rule_attachment : %s", err)
	}

	organizationsListResponse, res, err := apiClient.Client.OrganizationsAPI.OrganizationsList(context.TODO()).Id(organizationId).Execute()

	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}
	if len(organizationsListResponse.Data) != 1 {
		d.SetId("")
		return nil
	}

	billingRulesListResponse, res, err := apiClient.Client.PrometheusRulesAPI.PrometheusrulesList(context.TODO()).Id(billingRuleId).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}
	if len(billingRulesListResponse.Data) != 1 {
		d.SetId("")
		return nil
	}

	body := []int32{billingRuleId}
	_, response, err := apiClient.Client.OrganizationsAPI.OrganizationsDeletePrometheusrules(context.TODO(), organizationId).RequestBody(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(response, err))
	}

	d.SetId("")
	return nil
}

func ParseOrganizationBillingRuleAttachmentId(id string) (int32, int32, error) {
	list := strings.Split(id, "/")
	if len(list) != 2 {
		return 0, 0, fmt.Errorf("unable to determine taikun_organization_billing_rule_attachment ID")
	}

	organizationId, err := utils.Atoi32(list[0])
	billingRuleId, err2 := utils.Atoi32(list[1])
	if err != nil || err2 != nil {
		return 0, 0, fmt.Errorf("unable to determine taikun_organization_billing_rule_attachment ID")
	}

	return organizationId, billingRuleId, nil
}
