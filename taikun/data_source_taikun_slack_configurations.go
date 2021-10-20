package taikun

import (
	"context"

	"github.com/itera-io/taikungoclient/client/slack"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/taikungoclient/models"
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

func dataSourceTaikunSlackConfigurationsRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	dataSourceID := "all"

	params := slack.NewSlackListParams().WithV(ApiVersion)

	organizationIDData, organizationIDProvided := data.GetOk("organization_id")
	if organizationIDProvided {
		dataSourceID = organizationIDData.(string)
		organizationID, err := atoi32(dataSourceID)
		if err != nil {
			return diag.FromErr(err)
		}
		params = params.WithOrganizationID(&organizationID)
	}

	var slackConfigurationsList []*models.SlackConfigurationDto
	for {
		response, err := apiClient.client.Slack.SlackList(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		slackConfigurationsList = append(slackConfigurationsList, response.Payload.Data...)
		if len(slackConfigurationsList) == int(response.Payload.TotalCount) {
			break
		}
		offset := int32(len(slackConfigurationsList))
		params = params.WithOffset(&offset)
	}

	slackConfigurations := make([]map[string]interface{}, len(slackConfigurationsList))
	for i, rawSlackConfiguration := range slackConfigurationsList {
		slackConfigurations[i] = flattenTaikunSlackConfiguration(rawSlackConfiguration)
	}

	if err := data.Set("slack_configurations", slackConfigurations); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(dataSourceID)

	return nil
}
