package taikun

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunSlackConfigurations() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all Slack configurations.",
		ReadContext: dataSourceTaikunSlackConfigurationsRead,
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: stringIsInt,
			},
			"slack_configurations": {
				Description: "List of retrieved Slack configurations.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunSlackConfigurationSchema(),
				},
			},
		},
	}
}

func dataSourceTaikunSlackConfigurationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"
	var offset int32 = 0

	params := apiClient.Client.SlackAPI.SlackList(context.TODO())

	organizationIDData, organizationIDProvided := d.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.OrganizationId(organizationID)
	}

	var slackConfigurationsList []tkcore.SlackConfigurationDto
	for {
		response, res, err := params.Offset(offset).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		slackConfigurationsList = append(slackConfigurationsList, response.Data...)
		if len(slackConfigurationsList) == int(response.GetTotalCount()) {
			break
		}
		offset = int32(len(slackConfigurationsList))
	}

	slackConfigurations := make([]map[string]interface{}, len(slackConfigurationsList))
	for i, rawSlackConfiguration := range slackConfigurationsList {
		slackConfigurations[i] = flattenTaikunSlackConfiguration(&rawSlackConfiguration)
	}

	if err := d.Set("slack_configurations", slackConfigurations); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)

	return nil
}
