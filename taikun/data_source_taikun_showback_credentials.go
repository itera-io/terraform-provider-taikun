package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/client/showback"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunShowbackCredentials() *schema.Resource {
	return &schema.Resource{
		Description: "Get the list of showback credentials, optionally filtered by organization.",
		ReadContext: dataSourceTaikunShowbackCredentialsRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:  "Organization id filter.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: stringIsInt,
			},
			"showback_credentials": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The id of the showback credential.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "The name of the showback credential.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"prometheus_username": {
							Description: "The prometheus username.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"prometheus_password": {
							Description: "The prometheus password.",
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
						},
						"prometheus_url": {
							Description: "The prometheus url.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"organization_id": {
							Description: "The id of the organization which owns the showback credential.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"organization_name": {
							Description: "The name of the organization which owns the showback credential.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"is_locked": {
							Description: "Indicates whether the showback credential is locked or not.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"created_by": {
							Description: "The creator of the showback credential.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"last_modified": {
							Description: "Time of last modification.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"last_modified_by": {
							Description: "The last user who modified the showback credential.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTaikunShowbackCredentialsRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	params := showback.NewShowbackCredentialsListParams().WithV(ApiVersion)

	organizationIDData, organizationIDProvided := data.GetOk("organization_id")
	var organizationID int32 = -1
	if organizationIDProvided {
		organizationID, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.WithOrganizationID(&organizationID)
	}

	var showbackCredentialsList []*models.ShowbackCredentialsListDto
	for {
		response, err := apiClient.client.Showback.ShowbackCredentialsList(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		showbackCredentialsList = append(showbackCredentialsList, response.GetPayload().Data...)
		if len(showbackCredentialsList) == int(response.GetPayload().TotalCount) {
			break
		}
		offset := int32(len(showbackCredentialsList))
		params = params.WithOffset(&offset)
	}

	showbackCredentials := make([]map[string]interface{}, len(showbackCredentialsList), len(showbackCredentialsList))
	for i, rawShowbackCredential := range showbackCredentialsList {
		showbackCredentials[i] = flattenDatasourceTaikunShowbackCredentialItem(rawShowbackCredential)
	}
	if err := data.Set("showback_credentials", showbackCredentials); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(i32toa(organizationID))

	return nil
}

func flattenDatasourceTaikunShowbackCredentialItem(rawShowbackCredential *models.ShowbackCredentialsListDto) map[string]interface{} {

	return map[string]interface{}{
		"created_by":          rawShowbackCredential.CreatedBy,
		"id":                  i32toa(rawShowbackCredential.ID),
		"is_locked":           rawShowbackCredential.IsLocked,
		"last_modified":       rawShowbackCredential.LastModified,
		"last_modified_by":    rawShowbackCredential.LastModifiedBy,
		"name":                rawShowbackCredential.Name,
		"organization_id":     i32toa(rawShowbackCredential.OrganizationID),
		"organization_name":   rawShowbackCredential.OrganizationName,
		"prometheus_password": rawShowbackCredential.Password,
		"prometheus_url":      rawShowbackCredential.URL,
		"prometheus_username": rawShowbackCredential.Username,
	}
}
