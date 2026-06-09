package organization

import (
	"context"
	"regexp"

	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTaikunOrganizationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cloud_credentials": {
			Description: "Number of associated cloud credentials.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"created_at": {
			Description: "Time and date of creation.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"email": {
			Description: "Email.",
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
		},
		"full_name": {
			Description:  "Full name.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"id": {
			Description: "Organization's ID.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description: "Organization's name.",
			Type:        schema.TypeString,
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 30),
				validation.StringMatch(
					regexp.MustCompile("^[a-z0-9-_.]+$"),
					"expected only lowercase alpha numeric characters or non alpha numeric (_-.)",
				),
			),
		},
		"projects": {
			Description: "Number of associated projects.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"servers": {
			Description: "Number of associated servers.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
	}
}

func ResourceTaikunOrganization() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Organization",
		CreateContext: resourceTaikunOrganizationCreate,
		ReadContext:   generateResourceTaikunOrganizationReadWithoutRetries(),
		UpdateContext: resourceTaikunOrganizationUpdate,
		DeleteContext: resourceTaikunOrganizationDelete,
		Schema:        resourceTaikunOrganizationSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunOrganizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	body := tkcore.OrganizationCreateCommand{}
	body.SetName(d.Get("name").(string))
	body.SetFullName(d.Get("full_name").(string))
	if email, ok := d.GetOk("email"); ok {
		body.SetEmail(email.(string))
	}

	createResult, res, err := apiClient.Client.OrganizationsAPI.OrganizationsCreate(ctx).OrganizationCreateCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId(createResult.GetId())

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunOrganizationReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunOrganizationReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunOrganizationRead(true)
}
func generateResourceTaikunOrganizationReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunOrganizationRead(false)
}
func generateResourceTaikunOrganizationRead(withRetries bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)
		id := d.Id()
		id32, _ := utils.Atoi32(d.Id())
		d.SetId("")

		response, res, err := apiClient.Client.OrganizationsAPI.OrganizationsList(ctx).Id(id32).Execute()

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

		rawOrganization := response.Data[0]

		err = utils.SetResourceDataFromMap(d, flattenTaikunOrganization(&rawOrganization))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(id)

		return nil
	}
}

func resourceTaikunOrganizationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	id, err := utils.Atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	body := tkcore.UpdateOrganizationCommand{}
	body.SetId(id)
	body.SetName(d.Get("name").(string))
	body.SetFullName(d.Get("full_name").(string))

	res, err := apiClient.Client.OrganizationsAPI.OrganizationsUpdate(ctx).UpdateOrganizationCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	return utils.ReadAfterUpdateWithRetries(generateResourceTaikunOrganizationReadWithRetries(), ctx, d, meta)
}

func resourceTaikunOrganizationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	id, err := utils.Atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := apiClient.Client.OrganizationsAPI.OrganizationsDelete(ctx, id).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId("")
	return nil
}

func flattenTaikunOrganization(rawOrganization *tkcore.OrganizationDetailsDto) map[string]interface{} {
	return map[string]interface{}{
		"id":                utils.I32toa(rawOrganization.GetId()),
		"name":             rawOrganization.GetName(),
		"full_name":        rawOrganization.GetFullName(),
		"email":            rawOrganization.GetEmail(),
		"created_at":       rawOrganization.GetCreatedAt(),
		"cloud_credentials": rawOrganization.GetCloudCredentials(),
		"projects":         rawOrganization.GetProjects(),
		"servers":          rawOrganization.GetServers(),
	}
}
