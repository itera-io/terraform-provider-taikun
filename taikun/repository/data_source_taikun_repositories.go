package repository

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
)

func DataSourceTaikunRepositories() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve all repositories.",
		ReadContext: dataSourceTaikunRepositoriesRead,
		Schema: map[string]*schema.Schema{
			"repositories": {
				Description: "List of retrieved repositories.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: dataSourceTaikunRepositorySchema(),
				},
			},
			"organization_id": {
				Description:      "Organization ID filter.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: utils.StringIsInt,
			},
		},
	}
}

func dataSourceTaikunRepositoriesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"
	orgId, err := getSpecifiedOrDefaultOrganizationId(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	//  Public
	var offsetPublic int32 = 0
	params := apiClient.Client.AppRepositoriesAPI.RepositoryAvailableList(context.TODO()).OrganizationId(orgId).IsPrivate(false)
	var repositoriesList []tkcore.ArtifactRepositoryDto
	for {
		response, res, err := params.Offset(offsetPublic).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		repositoriesList = append(repositoriesList, response.Data...)
		if len(repositoriesList) == int(response.GetTotalCount()) {
			break
		}
		offsetPublic = int32(len(repositoriesList))
	}
	var publicTotalGot = len(repositoriesList)

	// Private
	var offsetPrivate int32 = 0
	params = apiClient.Client.AppRepositoriesAPI.RepositoryAvailableList(context.TODO()).OrganizationId(orgId).IsPrivate(true)
	for {
		response, res, err := params.Offset(offsetPrivate).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		repositoriesList = append(repositoriesList, response.Data...)
		if len(repositoriesList) == int(response.GetTotalCount())+publicTotalGot {
			break
		}
		offsetPrivate = int32(len(repositoriesList) - publicTotalGot)
	}

	// Flatten together
	repositories := make([]map[string]interface{}, len(repositoriesList))
	for i, rawRepository := range repositoriesList {
		if i < publicTotalGot {
			repositories[i] = flattenTaikunRepository(orgId, &rawRepository, false)
		} else {
			repositories[i] = flattenTaikunRepository(orgId, &rawRepository, true)
		}
	}

	if err := d.Set("repositories", repositories); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)
	return nil
}
