package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunBillingRule() *schema.Resource {
	return &schema.Resource{
		Description: "Get a billing rule by its id.",
		ReadContext: dataSourceTaikunBillingRuleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description:  "The id of the billing rule.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: stringIsInt,
			},
			"name": {
				Description: "The name of the billing rule.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"metric_name": {
				Description: "The name of the metric from Prometheus you want to bill.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"label": {
				Description: "Labels linked to this billing rule.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label": {
							Description: "Key of the label.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"value": {
							Description: "Value of the label.",
							Type:        schema.TypeString,
							Computed:    true,
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
				Description: "Type of the billing rule. `Count` (calculate package as unit) or `Sum` (calculate per quantity)",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"price": {
				Description: "The price in CZK per selected unit.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"created_by": {
				Description: "The creator of the billing credential.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"billing_credential_id": {
				Description: "Id of the billing credential.",
				Type:        schema.TypeString,
				Computed:    true,
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
	}
}

func dataSourceTaikunBillingRuleRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	data.SetId(data.Get("id").(string))

	return resourceTaikunBillingRuleRead(ctx, data, meta)
}
