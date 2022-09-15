package taikun

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/organizations"
	"github.com/itera-io/taikungoclient/client/prometheus"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunOrganizationBillingRuleAttachmentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"billing_rule_id": {
			Description:      "ID of the billing rule.",
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: stringIsInt,
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
			ValidateDiagFunc: stringIsInt,
		},
		"organization_name": {
			Description: "Name of the organisation.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func resourceTaikunOrganizationBillingRuleAttachment() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Organization - Billing Rule Attachment",
		CreateContext: resourceTaikunOrganizationBillingRuleAttachmentCreate,
		ReadContext:   generateResourceTaikunOrganizationBillingRuleAttachmentReadWithoutRetries(),
		DeleteContext: resourceTaikunOrganizationBillingRuleAttachmentDelete,
		Schema:        resourceTaikunOrganizationBillingRuleAttachmentSchema(),
	}
}

func resourceTaikunOrganizationBillingRuleAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*taikungoclient.Client)

	billingRuleId, err := atoi32(d.Get("billing_rule_id").(string))
	if err != nil {
		return diag.Errorf("billing_rule_id isn't valid: %s", d.Get("billing_rule_id").(string))
	}

	organizationId, err := atoi32(d.Get("organization_id").(string))
	if err != nil {
		return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
	}

	body := &models.BindPrometheusOrganizationsCommand{
		Organizations: []*models.BindOrganizationsToRuleDto{
			{
				IsBound:          true,
				OrganizationID:   organizationId,
				RuleDiscountRate: d.Get("discount_rate").(float64),
			},
		},
		PrometheusRuleID: billingRuleId,
	}
	params := prometheus.NewPrometheusBindOrganizationsParams().WithV(ApiVersion).WithBody(body)
	_, err = client.client.Prometheus.PrometheusBindOrganizations(params, client)
	if err != nil {
		return diag.FromErr(err)
	}

	id := fmt.Sprintf("%d/%d", organizationId, billingRuleId)
	d.SetId(id)

	return readAfterCreateWithRetries(generateResourceTaikunOrganizationBillingRuleAttachmentReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunOrganizationBillingRuleAttachmentReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunOrganizationBillingRuleAttachmentRead(true)
}
func generateResourceTaikunOrganizationBillingRuleAttachmentReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunOrganizationBillingRuleAttachmentRead(false)
}
func generateResourceTaikunOrganizationBillingRuleAttachmentRead(withRetries bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*taikungoclient.Client)

		id := d.Id()
		d.SetId("")
		organizationId, billingRuleId, err := parseOrganizationBillingRuleAttachmentId(id)
		if err != nil {
			return diag.Errorf("Error while reading taikun_organization_billing_rule_attachment : %s", err)
		}

		params := prometheus.NewPrometheusListOfRulesParams().WithV(ApiVersion).WithID(&billingRuleId)
		response, err := apiClient.Client.Prometheus.PrometheusListOfRules(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(response.Payload.Data) != 1 {
			if withRetries {
				d.SetId(id)
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawBillingRule := response.GetPayload().Data[0]

		for _, e := range rawBillingRule.BoundOrganizations {
			if e.OrganizationID == organizationId {
				if err := d.Set("organization_id", i32toa(e.OrganizationID)); err != nil {
					return diag.FromErr(err)
				}
				if err := d.Set("organization_name", e.OrganizationName); err != nil {
					return diag.FromErr(err)
				}
				if err := d.Set("billing_rule_id", i32toa(rawBillingRule.ID)); err != nil {
					return diag.FromErr(err)
				}
				if err := d.Set("discount_rate", e.RuleDiscountRate); err != nil {
					return diag.FromErr(err)
				}
				d.SetId(id)
				return nil
			}
		}

		if withRetries {
			d.SetId(id)
			return diag.Errorf(notFoundAfterCreateOrUpdateError)
		}
		return nil
	}
}

func resourceTaikunOrganizationBillingRuleAttachmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)

	organizationId, billingRuleId, err := parseOrganizationBillingRuleAttachmentId(d.Id())
	if err != nil {
		return diag.Errorf("Error while deleting taikun_organization_billing_rule_attachment : %s", err)
	}

	organizationsListParams := organizations.NewOrganizationsListParams().WithV(ApiVersion).WithID(&organizationId)
	organizationsListResponse, err := apiClient.Client.Organizations.OrganizationsList(organizationsListParams, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(organizationsListResponse.Payload.Data) != 1 {
		d.SetId("")
		return nil
	}

	billingRulesListParams := prometheus.NewPrometheusListOfRulesParams().WithV(ApiVersion).WithID(&billingRuleId)
	billingRulesListResponse, err := apiClient.Client.Prometheus.PrometheusListOfRules(billingRulesListParams, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(billingRulesListResponse.Payload.Data) != 1 {
		d.SetId("")
		return nil
	}

	body := &models.BindPrometheusOrganizationsCommand{
		Organizations: []*models.BindOrganizationsToRuleDto{
			{
				IsBound:        false,
				OrganizationID: organizationId,
			},
		},
		PrometheusRuleID: billingRuleId,
	}
	params := prometheus.NewPrometheusBindOrganizationsParams().WithV(ApiVersion).WithBody(body)
	_, err = apiClient.Client.Prometheus.PrometheusBindOrganizations(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func parseOrganizationBillingRuleAttachmentId(id string) (int32, int32, error) {
	list := strings.Split(id, "/")
	if len(list) != 2 {
		return 0, 0, fmt.Errorf("unable to determine taikun_organization_billing_rule_attachment ID")
	}

	organizationId, err := atoi32(list[0])
	billingRuleId, err2 := atoi32(list[1])
	if err != nil || err2 != nil {
		return 0, 0, fmt.Errorf("unable to determine taikun_organization_billing_rule_attachment ID")
	}

	return organizationId, billingRuleId, nil
}
