package catalog

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"regexp"
)

func taikunApplicationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "The ID of the application.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description: "The name of the application.",
			Type:        schema.TypeString,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 30),
				validation.StringMatch(
					regexp.MustCompile("^[a-z0-9-]+$"),
					"Application name must contain only alpha numeric characters or non alpha numeric (-)",
				),
			),
			Required: true,
		},
		"repository": {
			Description: "The name of the repository.",
			Type:        schema.TypeString,
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 30),
				validation.StringMatch(
					regexp.MustCompile("^[a-z0-9-]+$"),
					"Application name must contain only alpha numeric characters or non alpha numeric (-)",
				),
			),
		},
	}
}

func resourceTaikunCatalogSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "The ID of the catalog.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description: "The name of the catalog.",
			Type:        schema.TypeString,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 30),
				validation.StringMatch(
					regexp.MustCompile("^[a-z0-9-]+$"),
					"Catalog name must contain only alpha numeric characters or non alpha numeric (-)",
				),
			),
			Required: true,
		},
		"description": {
			Description: "The description of the catalog.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"lock": {
			Description: "Indicates whether to lock the catalog.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"default": {
			Description: "Indicates whether the catalog is the default catalog.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"organization_id": {
			Description:      "The ID of the organization which owns the catalog.",
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			ValidateDiagFunc: utils.StringIsInt,
		},
		"application": {
			Description: "Bound Applications.",
			Type:        schema.TypeSet,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: taikunApplicationSchema(),
			},
		},
		"projects": {
			Description: "DEPRECATED: List of projects bound to the catalog.",
			Type:        schema.TypeSet,
			Deprecated:  "Please use the resource taikun_catalog_project_binding to bind projects to the catalog.",
			Optional:    true,
			//DefaultFunc: func() (interface{}, error) {
			//	return []interface{}{}, nil
			//},
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func ResourceTaikunCatalog() *schema.Resource {
	return &schema.Resource{
		Description:   "Catalog for Taikun Applications Configuration.",
		CreateContext: resourceTaikunCatalogCreate,
		ReadContext:   generateResourceTaikunCatalogReadWithoutRetries(),
		UpdateContext: resourceTaikunCatalogUpdate,
		DeleteContext: resourceTaikunCatalogDelete,
		Schema:        resourceTaikunCatalogSchema(),
	}
}

func resourceTaikunCatalogCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	applicationsWhichShouldBeBound := d.Get("application").(*schema.Set)
	// TODO - Legacy bind projects
	projects := d.Get("projects")
	if projects == nil {
		projects = schema.NewSet(schema.HashString, []interface{}{})
	}
	legacyProjectsWhichShouldBeBound := projects.(*schema.Set)

	// Create catalog
	body := &tkcore.CreateCatalogCommand{}
	body.SetName(d.Get("name").(string))
	body.SetDescription(d.Get("description").(string))
	if organizationIDData, organizationIDIsSet := d.GetOk("organization_id"); organizationIDIsSet {
		orgId, err := utils.Atoi32(organizationIDData.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		body.SetOrganizationId(orgId)
	}

	response, err := apiClient.Client.CatalogAPI.CatalogCreate(context.TODO()).CreateCatalogCommand(*body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(response, err))
	}

	// Get catalogId
	errDiag := utils.ReadAfterCreateWithRetries(generateResourceTaikunCatalogReadWithRetries(), ctx, d, meta) // Get the ID
	if errDiag != nil {
		return errDiag
	}
	catalogId, err := utils.Atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Bind Legacy projects
	emptyBoundProjects := schema.Set{}
	errDiag = reconcileProjectsBound(&emptyBoundProjects, legacyProjectsWhichShouldBeBound, catalogId, meta)
	if errDiag != nil {
		return errDiag
	}

	// Bind applications
	emptyBoundApplications := schema.Set{}
	errDiag = reconcileApplicationsBound(&emptyBoundApplications, applicationsWhichShouldBeBound, catalogId, meta)
	if errDiag != nil {
		return errDiag
	}

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunCatalogReadWithRetries(), ctx, d, meta)
}

func resourceTaikunCatalogDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	catalogId, err := utils.Atoi32(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	// Unbind all projects
	emptyBoundProjects := schema.Set{}
	projects := d.Get("projects")
	if projects == nil {
		projects = schema.NewSet(schema.HashString, []interface{}{})
	}
	legacyProjectsWhichShouldBeUnbound := projects.(*schema.Set)
	errDiag := reconcileProjectsBound(legacyProjectsWhichShouldBeUnbound, &emptyBoundProjects, catalogId, meta)
	if errDiag != nil {
		return errDiag
	}

	// Unbind all applications
	emptyBoundApplications := schema.Set{}
	errDiag = reconcileApplicationsBound(d.Get("application"), &emptyBoundApplications, catalogId, meta)
	if errDiag != nil {
		return errDiag
	}

	// Delete catalog
	response, err := apiClient.Client.CatalogAPI.CatalogDelete(context.TODO(), catalogId).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(response, err))
	}
	d.SetId("")
	return nil
}

func generateResourceTaikunCatalogReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunCatalogRead(true)
}
func generateResourceTaikunCatalogReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunCatalogRead(false)
}

func generateResourceTaikunCatalogRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)
		catalogName := d.Get("name").(string)
		listQuery := apiClient.Client.CatalogAPI.CatalogList(context.TODO()).Search(catalogName)

		if organizationIDData, organizationIDIsSet := d.GetOk("organization_id"); organizationIDIsSet {
			orgId, err := utils.Atoi32(organizationIDData.(string))
			if err != nil {
				return diag.FromErr(err)
			}
			listQuery = listQuery.OrganizationId(orgId)
		}

		data, response, err := listQuery.Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(response, err))
		}

		// Iterate through data to find the correct Catalog
		foundMatch := false
		var rawCatalog tkcore.CatalogListDto
		for _, catalog := range data.GetData() {
			if catalog.GetName() == catalogName {
				foundMatch = true
				rawCatalog = catalog
				break
			}
		}

		if !foundMatch {
			if withRetries {
				d.SetId(d.Get("id").(string)) // We need to tell provider that object was created
				return diag.FromErr(fmt.Errorf("could not find the specified catalog (name: %s)", catalogName))
			}
			return nil
		}

		// Load all the found data to the local object
		err = utils.SetResourceDataFromMap(d, flattenTaikunCatalog(&rawCatalog))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(d.Get("id").(string)) // We need to tell provider that object was created

		return nil
	}
}

func flattenTaikunCatalog(rawCatalog *tkcore.CatalogListDto) map[string]interface{} {

	// Flatten bound applications
	applicationsBound := rawCatalog.GetBoundApplications()
	applications := make([]map[string]interface{}, 0)
	for _, app := range applicationsBound {
		appMap := map[string]interface{}{
			"id":         utils.I32toa(app.GetCatalogAppId()),
			"name":       app.GetName(),
			"repository": app.Repository.GetName(),
		}
		applications = append(applications, appMap)
	}

	// Flatten bound projects
	boundProjects := rawCatalog.GetBoundProjects()
	projects := make([]string, len(boundProjects))
	for i, proj := range boundProjects {
		projects[i] = utils.I32toa(proj.GetId())
	}

	return map[string]interface{}{
		"id":          utils.I32toa(rawCatalog.GetId()),
		"name":        rawCatalog.GetName(),
		"description": rawCatalog.GetDescription(),
		"lock":        rawCatalog.GetIsLocked(),
		"default":     rawCatalog.GetIsDefault(),
		"projects":    projects,
		"application": applications,
	}
}

func resourceTaikunCatalogUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	// Name, description
	catalogId, err := utils.Atoi32(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	_, newName := d.GetChange("name")
	_, newDescription := d.GetChange("description")
	updatedCatalog := tkcore.EditCatalogCommand{}
	updatedCatalog.SetId(catalogId)
	updatedCatalog.SetName(newName.(string))
	updatedCatalog.SetDescription(newDescription.(string))
	response, err := apiClient.Client.CatalogAPI.CatalogEdit(context.TODO()).EditCatalogCommand(updatedCatalog).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(response, err))
	}

	// lock
	oldCatalogLocked, newCatalogLocked := d.GetChange("lock")
	if oldCatalogLocked != newCatalogLocked {
		updateLock := tkcore.CatalogLockManagementCommand{}
		updateLock.SetId(catalogId)
		updateLock.SetMode(utils.GetLockMode(newCatalogLocked.(bool)))
		response, err = apiClient.Client.CatalogAPI.CatalogLock(context.TODO()).CatalogLockManagementCommand(updateLock).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(response, err))
		}
	}

	// default # If we set multiple default, the last one stays
	oldDefault, newDefault := d.GetChange("default")
	if oldDefault != newDefault {
		if newDefault.(bool) {
			updateDefault := tkcore.CatalogMakeDefaultCommand{}
			updateDefault.SetId(catalogId)
			response, err = apiClient.Client.CatalogAPI.CatalogMakeDefault(context.TODO()).CatalogMakeDefaultCommand(updateDefault).Execute()
			if err != nil {
				return diag.FromErr(tk.CreateError(response, err))
			}
		} else {
			// This will get changed by other catalog taking the default status
			_ = d.Set("default", false)
		}
	}

	// Binding applications
	oldCatalogApplicationsBound, newCatalogApplicationsBound := d.GetChange("application")
	errReconcile := reconcileApplicationsBound(oldCatalogApplicationsBound, newCatalogApplicationsBound, catalogId, meta)
	if errReconcile != nil {
		return errReconcile
	}

	return utils.ReadAfterUpdateWithRetries(generateResourceTaikunCatalogReadWithRetries(), ctx, d, meta)
}

// Unbind apps that should be unbound, bind apps that should be bound
func reconcileApplicationsBound(oldCatalogApplicationsBound interface{}, newCatalogApplicationsBound interface{}, catalogId int32, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	oldApplications := oldCatalogApplicationsBound.(*schema.Set)
	newApplications := newCatalogApplicationsBound.(*schema.Set)

	// Old applications that we should delete
	toRemove := oldApplications.Difference(newApplications)
	for _, app := range toRemove.List() {
		catalogAppId, err := utils.Atoi32(app.(map[string]interface{})["id"].(string))
		if err != nil {
			diag.FromErr(err)
		}
		response, err := apiClient.Client.CatalogAppAPI.CatalogAppDelete(context.TODO(), catalogAppId).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(response, err))
		}
	}

	// New applications that we should create
	toAdd := newApplications.Difference(oldApplications)
	for _, app := range toAdd.List() {
		catalogAppToCreate := tkcore.CreateCatalogAppCommand{}
		catalogAppToCreate.SetCatalogId(catalogId)
		catalogAppToCreate.SetRepoName(app.(map[string]interface{})["repository"].(string))
		catalogAppToCreate.SetPackageName(app.(map[string]interface{})["name"].(string))
		catalogAppToCreate.SetParameters([]tkcore.CatalogAppParamsDto{})
		_, response, err := apiClient.Client.CatalogAppAPI.CatalogAppCreate(context.TODO()).CreateCatalogAppCommand(catalogAppToCreate).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(response, err))
		}
	}

	return nil
}

// Ensure the catalog has only all the projects in newCatalogProjectsBound bound
func reconcileProjectsBound(oldCatalogProjectsBound interface{}, newCatalogProjectsBound interface{}, catalogId int32, meta interface{}) diag.Diagnostics {
	oldProjects := oldCatalogProjectsBound.(*schema.Set)
	newProjects := newCatalogProjectsBound.(*schema.Set)
	apiClient := meta.(*tk.Client)

	// Old stuff that we should unbind - What was in old catalog, but is not new catalog.
	toRemove := oldProjects.Difference(newProjects).List()
	if len(toRemove) > 0 {
		body, err := utils.SliceOfSTringsToSliceOfInt32(toRemove)
		if err != nil {
			return diag.FromErr(err)
		}
		response, err := apiClient.Client.CatalogAPI.CatalogDeleteProject(context.TODO(), catalogId).RequestBody(body).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(response, err))
		}
	}

	// New stuff that we should bind - What was is in new catalog and was not in old catalog.
	toAdd := newProjects.Difference(oldProjects).List()
	if len(toAdd) > 0 {
		body, err := utils.SliceOfSTringsToSliceOfInt32(toAdd)
		if err != nil {
			return diag.FromErr(err)
		}
		response, err := apiClient.Client.CatalogAPI.CatalogAddProject(context.TODO(), catalogId).RequestBody(body).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(response, err))
		}
	}

	return nil
}
