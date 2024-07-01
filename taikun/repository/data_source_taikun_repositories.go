package repository

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
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
		},
	}
}

func dataSourceTaikunRepositoriesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	dataSourceID := "all"

	//  Public
	var offsetPublic int32 = 0
	params := apiClient.Client.AppRepositoriesAPI.RepositoryAvailableList(context.TODO()).IsPrivate(false)
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
	var publicTotalGot int = len(repositoriesList)

	// Private
	var offsetPrivate int32 = 0
	params = apiClient.Client.AppRepositoriesAPI.RepositoryAvailableList(context.TODO()).IsPrivate(true)
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
		repositories[i] = flattenTaikunRepository(&rawRepository)
	}
	if err := d.Set("repositories", repositories); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dataSourceID)
	return nil
}
