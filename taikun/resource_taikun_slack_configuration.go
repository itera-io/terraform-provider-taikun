package taikun

import (
	"context"
	tk "github.com/chnyda/taikungoclient"
	tkcore "github.com/chnyda/taikungoclient/client"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTaikunSlackConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"channel": {
			Description: "Slack channel for notifications.",
			Type:        schema.TypeString,
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 40),
			),
		},
		"id": {
			Description: "The Slack configuration's ID.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description: "The Slack configuration's name.",
			Type:        schema.TypeString,
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 40),
				validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z0-9-_.]+$"),
					"expected only alpha numeric characters or non alpha numeric (_-.)",
				),
			),
		},
		"organization_id": {
			Description:      "The ID of the organization which owns the Slack configuration.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"organization_name": {
			Description: "The name of the organization which owns the Slack configuration.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"type": {
			Description:  "The type of notifications to receive: `Alert` (only alert-type notifications) or `General` (all notifications).",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"Alert", "General"}, false),
		},
		"url": {
			Description:  "Webhook URL from Slack app.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
		},
	}
}

func resourceTaikunSlackConfiguration() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Slack Configuration",
		CreateContext: resourceTaikunSlackConfigurationCreate,
		ReadContext:   generateResourceTaikunSlackConfigurationReadWithoutRetries(),
		UpdateContext: resourceTaikunSlackConfigurationUpdate,
		DeleteContext: resourceTaikunSlackConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: resourceTaikunSlackConfigurationSchema(),
	}
}

func resourceTaikunSlackConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	body := tkcore.CreateSlackConfigurationCommand{}
	body.SetName(d.Get("name").(string))
	body.SetUrl(d.Get("url").(string))
	body.SetChannel(d.Get("channel").(string))
	body.SetSlackType(tkcore.SlackType(d.Get("type").(string)))

	if organizationIDData, organizationIDIsSet := d.GetOk("organization_id"); organizationIDIsSet {
		organizationID, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		body.SetOrganizationId(organizationID)
	}

	response, res, err := apiClient.Client.SlackApi.SlackCreate(context.TODO()).CreateSlackConfigurationCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId(response.GetId())

	return readAfterCreateWithRetries(generateResourceTaikunSlackConfigurationReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunSlackConfigurationReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunSlackConfigurationRead(true)
}
func generateResourceTaikunSlackConfigurationReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunSlackConfigurationRead(false)
}
func generateResourceTaikunSlackConfigurationRead(withRetries bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)

		id, err := atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, res, err := apiClient.Client.SlackApi.SlackList(context.TODO()).Id(id).Execute()
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

		rawSlackConfiguration := response.Data[0]

		err = setResourceDataFromMap(d, flattenTaikunSlackConfiguration(&rawSlackConfiguration))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunSlackConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	body := tkcore.UpdateSlackConfigurationDto{}
	body.SetName(d.Get("name").(string))
	body.SetUrl(d.Get("url").(string))
	body.SetChannel(d.Get("channel").(string))
	body.SetSlackType(tkcore.SlackType(d.Get("type").(string)))

	if organizationIDData, organizationIDIsSet := d.GetOk("organization_id"); organizationIDIsSet {
		organizationID, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		body.SetOrganizationId(organizationID)
	}

	if res, err := apiClient.Client.SlackApi.SlackUpdate(context.TODO(), id).UpdateSlackConfigurationDto(body).Execute(); err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	return readAfterUpdateWithRetries(generateResourceTaikunSlackConfigurationReadWithRetries(), ctx, d, meta)
}

func resourceTaikunSlackConfigurationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	body := tkcore.DeleteSlackConfigCommand{Ids: []int32{id}}
	res, err := apiClient.Client.SlackApi.SlackDeleteMultiple(context.TODO()).DeleteSlackConfigCommand(body).Execute()

	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId("")

	return nil
}

func flattenTaikunSlackConfiguration(rawSlackConfiguration *tkcore.SlackConfigurationDto) map[string]interface{} {
	return map[string]interface{}{
		"channel":           rawSlackConfiguration.GetChannel(),
		"id":                i32toa(rawSlackConfiguration.GetId()),
		"name":              rawSlackConfiguration.GetName(),
		"organization_id":   i32toa(rawSlackConfiguration.GetOrganizationId()),
		"organization_name": rawSlackConfiguration.GetOrganizationName(),
		"type":              rawSlackConfiguration.GetSlackType(),
		"url":               rawSlackConfiguration.GetUrl(),
	}
}
