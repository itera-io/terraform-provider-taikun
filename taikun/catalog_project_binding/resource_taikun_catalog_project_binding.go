package catalog_project_binding

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"log"
)

func resourceTaikunCatalogProjectBindingSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "The ID of the catalog binding is catalogName+projectId.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"catalog_name": {
			Description: "The name of the catalog you wish to bind.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"project_id": {
			Description:      "The ID of the READY project you wish to bind to catalog.",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: utils.StringIsInt,
		},
		"is_bound": {
			Description: "Specifies whether the catalog will be bound to the project or not.",
			Type:        schema.TypeBool,
			Required:    true,
		},
		"organization_id": {
			Description: "The ID of the organization which owns both the catalog and the project.",
			Type:        schema.TypeString,
			Optional:    true,
			//Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: utils.StringIsInt,
		},
	}
}

func ResourceTaikunCatalogProjectBinding() *schema.Resource {
	return &schema.Resource{
		Description:   "Catalog for Taikun Applications Configuration.",
		CreateContext: resourceTaikunCatalogProjectBindingCreate,
		ReadContext:   generateResourceTaikunCatalogProjectBindingReadWithoutRetries(),
		UpdateContext: resourceTaikunCatalogProjectBindingUpdate,
		DeleteContext: resourceTaikunCatalogProjectBindingDelete,
		Schema:        resourceTaikunCatalogProjectBindingSchema(),
	}
}

func resourceTaikunCatalogProjectBindingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	var orgId int32
	if organizationIDData, organizationIDIsSet := d.GetOk("organization_id"); organizationIDIsSet {
		organizationIDDataConverted, err := utils.Atoi32(organizationIDData.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		orgId = organizationIDDataConverted
	}
	log.Println("CREATE organization_id:", orgId)

	projectId, err := utils.Atoi32(d.Get("project_id").(string))
	log.Println("CREATE project_id:", projectId)
	if err != nil {
		return diag.FromErr(err)
	}
	catalogName := d.Get("catalog_name").(string)
	shouldBeBound := d.Get("is_bound").(bool)
	err = reconcileBinding(apiClient, orgId, catalogName, projectId, shouldBeBound)
	if err != nil {
		return diag.FromErr(err)
	}

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunCatalogProjectBindingReadWithRetries(), ctx, d, meta)
}

func resourceTaikunCatalogProjectBindingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	//organizationId, err := utils.Atoi32(d.Get("organization_id").(string))
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	var orgId int32
	if organizationIDData, organizationIDIsSet := d.GetOk("organization_id"); organizationIDIsSet {
		organizationIDDataConverted, err := utils.Atoi32(organizationIDData.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		orgId = organizationIDDataConverted
	}

	projectId, err := utils.Atoi32(d.Get("project_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	catalogName := d.Get("catalog_name").(string)

	// Does catalog exist and get what projects it has bound?
	foundCatalog, err := findCatalogByName(apiClient, orgId, catalogName)
	if err != nil {
		return diag.FromErr(err)
	}

	catalogHasProjectBound := false
	for _, boundProject := range foundCatalog.BoundProjects {
		if boundProject.GetId() == projectId {
			catalogHasProjectBound = true
		}
	}

	// If not already, unbind the project from the catalog
	if catalogHasProjectBound {
		body := []int32{projectId}
		response, err := apiClient.Client.CatalogAPI.CatalogDeleteProject(context.TODO(), foundCatalog.GetId()).RequestBody(body).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(response, err))
		}
	}

	return nil
}

func generateResourceTaikunCatalogProjectBindingReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunCatalogProjectBindingRead(true)
}
func generateResourceTaikunCatalogProjectBindingReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunCatalogProjectBindingRead(false)
}

func generateResourceTaikunCatalogProjectBindingRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)

		var orgId int32
		if organizationIDData, organizationIDIsSet := d.GetOk("organization_id"); organizationIDIsSet {
			organizationIDDataConverted, err := utils.Atoi32(organizationIDData.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			orgId = organizationIDDataConverted
		}

		projectId, err := utils.Atoi32(d.Get("project_id").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		catalogName := d.Get("catalog_name").(string)

		// Does catalog exist and get what projects it has bound?
		foundCatalog, err := findCatalogByName(apiClient, orgId, catalogName)
		if err != nil {
			return diag.FromErr(err)
		}

		catalogHasProjectBound := false
		for _, boundProject := range foundCatalog.BoundProjects {
			if boundProject.GetId() == projectId {
				catalogHasProjectBound = true
			}
		}

		// Load all the found data to the local object
		err = utils.SetResourceDataFromMap(d, flattenTaikunCatalogProjectBinding(catalogHasProjectBound, d.Get("project_id").(string), catalogName))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(d.Get("catalog_name").(string) + d.Get("project_id").(string)) // We need to tell provider that object was created

		return nil
	}
}

func flattenTaikunCatalogProjectBinding(catalogHasProjectBound bool, projectId, catalogName string) map[string]interface{} {
	return map[string]interface{}{
		"id":           catalogName + projectId,
		"catalog_name": catalogName,
		"project_id":   projectId,
		"is_bound":     catalogHasProjectBound,
	}
}

func resourceTaikunCatalogProjectBindingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	//organizationId, err := utils.Atoi32(d.Get("organization_id").(string))
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	var orgId int32
	if organizationIDData, organizationIDIsSet := d.GetOk("organization_id"); organizationIDIsSet {
		organizationIDDataConverted, err := utils.Atoi32(organizationIDData.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		orgId = organizationIDDataConverted
	}

	projectId, err := utils.Atoi32(d.Get("project_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	catalogName := d.Get("catalog_name").(string)
	shouldBeBound := d.Get("is_bound").(bool)

	err = reconcileBinding(apiClient, orgId, catalogName, projectId, shouldBeBound)
	if err != nil {
		return diag.FromErr(err)
	}

	return utils.ReadAfterUpdateWithRetries(generateResourceTaikunCatalogProjectBindingReadWithRetries(), ctx, d, meta)
}

func reconcileBinding(apiClient *tk.Client, organizationId int32, catalogName string, projectId int32, shouldBeBound bool) error {
	log.Printf("Start reconcile binding")
	query := apiClient.Client.ProjectsAPI.ProjectsList(context.TODO()).Id(projectId)

	if organizationId != 0 {
		query = query.OrganizationId(organizationId)
	}

	// Verify project exists and is ready
	data, response, err := query.Execute()
	if err != nil {
		return tk.CreateError(response, err)
	}
	if data.Data[0].Status != tkcore.PROJECTSTATUS_READY {
		return fmt.Errorf("project '%d' is not ready", projectId)
	}

	// Verify catalog exists and get what projects it has bound
	foundCatalog, err := findCatalogByName(apiClient, organizationId, catalogName)
	if err != nil {
		return err
	}

	catalogHasProjectBound := false
	for _, boundProject := range foundCatalog.BoundProjects {
		if boundProject.GetId() == projectId {
			catalogHasProjectBound = true
		}
	}

	// If you need, unbind the project from the catalog
	if catalogHasProjectBound && !shouldBeBound {
		body := []int32{projectId}
		response, err := apiClient.Client.CatalogAPI.CatalogDeleteProject(context.TODO(), foundCatalog.GetId()).RequestBody(body).Execute()
		if err != nil {
			return tk.CreateError(response, err)
		}
	}

	// If you need, bind the project to the catalog
	if !catalogHasProjectBound && shouldBeBound {
		body := []int32{projectId}
		response, err := apiClient.Client.CatalogAPI.CatalogAddProject(context.TODO(), foundCatalog.GetId()).RequestBody(body).Execute()
		if err != nil {
			return tk.CreateError(response, err)
		}
	}

	return nil
}

func findCatalogByName(apiClient *tk.Client, organizationId int32, catalogName string) (rawCatalog tkcore.CatalogListDto, err error) {
	log.Printf("Start find catalog by name")

	query := apiClient.Client.CatalogAPI.CatalogList(context.TODO()).Search(catalogName)
	if organizationId != 0 {
		log.Printf("Setting org find catalog by name")
		query = query.OrganizationId(organizationId)
	}
	log.Printf("Executing find catalog by name")
	data, response, err := query.Execute()
	if err != nil {
		log.Printf("Failed find catalog by name")
		return rawCatalog, tk.CreateError(response, err)
	}
	log.Printf("Success find catalog by name")

	// Iterate through data to find the correct Catalog
	foundMatch := false
	for _, catalog := range data.GetData() {
		if catalog.GetName() == catalogName {
			foundMatch = true
			rawCatalog = catalog
			break
		}
	}

	if !foundMatch {
		if organizationId != 0 {
			return rawCatalog, fmt.Errorf("catalog '%s' not found in organization %d", catalogName, organizationId)
		}
		return rawCatalog, fmt.Errorf("catalog '%s' not found in default organization", catalogName)
	}

	log.Printf("Finished find catalog by name")
	return rawCatalog, nil
}
