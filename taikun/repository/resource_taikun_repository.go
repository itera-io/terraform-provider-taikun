package repository

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"time"
)

func resourceTaikunRepositorySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// repositoryId - user does not set
		"id": {
			Description: "The ID of the repository.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"id_apprepo": {
			Description: "The ID of the application repository.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description:  "The name of the repository.",
			Type:         schema.TypeString,
			ValidateFunc: validation.StringLenBetween(3, 30),
			Required:     true,
			ForceNew:     true,
		},
		// Taikun does not allow picking into which organization we should create the private repository
		// Before creating private repository, we need to check if it matches with the logged in organization
		"organization_name": {
			Description: "The name of the organization which owns the public repository.",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		// isPrivate
		"private": {
			Description: "Indicates whether the repository is private or public.",
			Type:        schema.TypeBool,
			Required:    true,
			ForceNew:    true,
		},
		// url - Required when is_private is enabled, otherwise it gets filled from server
		"url": {
			Description:  "The URL of the repository.",
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			ForceNew:     true,
		},
		"enabled": {
			Description: "Indicates whether the repository is enabled.",
			Type:        schema.TypeBool,
			Required:    true,
		},
		// username - Optional, for private repos, does not get downloaded
		"username": {
			Description:  "The registry username. (Can be set with env REGISTRY_USERNAME)",
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			Default:      "",
			DefaultFunc:  schema.EnvDefaultFunc("REGISTRY_USERNAME", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		// password - Optional, for private repos, does not get downloaded
		"password": {
			Description:  "The registry password. (Can be set with env REGISTRY_PASSWORD)",
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("REGISTRY_PASSWORD", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func ResourceTaikunRepository() *schema.Resource {
	return &schema.Resource{
		Description:   "Public Repository for Taikun Applications Configuration.",
		CreateContext: resourceTaikunRepositoryCreate,
		ReadContext:   generateResourceTaikunRepositoryReadWithoutRetries(),
		UpdateContext: resourceTaikunRepositoryUpdate,
		DeleteContext: resourceTaikunRepositoryDelete, // Skip if public, we cannot delete public.
		Schema:        resourceTaikunRepositorySchema(),
	}
}

func resourceTaikunRepositoryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.Get("private") == true {
		apiClient := meta.(*tk.Client)
		err := ensureDesiredState(false, d.Get("enabled").(bool), d, meta)
		if err != nil {
			return diag.FromErr(err)
		}
		deleteCommand := tkcore.DeleteRepositoryCommand{}
		apprepoId, err := utils.Atoi32(d.Get("id_apprepo").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		deleteCommand.SetAppRepoId(apprepoId)
		response, err2 := apiClient.Client.AppRepositoriesAPI.RepositoryDelete(context.TODO()).DeleteRepositoryCommand(deleteCommand).Execute()
		if err2 != nil {
			return diag.FromErr(tk.CreateError(response, err2))
		}
	}
	d.SetId("")
	return nil
}

func resourceTaikunRepositoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	oldEnabled, newEnabled := d.GetChange("enabled")
	err := ensureDesiredState(newEnabled.(bool), oldEnabled.(bool), d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return utils.ReadAfterUpdateWithRetries(generateResourceTaikunRepositoryReadWithRetries(), ctx, d, meta)
}

func resourceTaikunRepositoryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	// Is this private or public repo we are talking about?
	should_be_private := bool(d.Get("private").(bool))
	should_be_enabled := d.Get("enabled").(bool)
	name := d.Get("name").(string)

	if should_be_private {
		// Check if organization specified matches organization declared
		orgIsValid, err := checkOrganization(d, meta)
		if !orgIsValid || err != nil {
			return diag.FromErr(err)
		}
		// Send create query
		body_private := &tkcore.ImportRepoCommand{}
		url_private := d.Get("url").(string)
		username_private := d.Get("username").(string)
		password_private := d.Get("password").(string)
		body_private.SetName(name)
		body_private.SetUrl(url_private)
		body_private.SetUsername(username_private)
		body_private.SetPassword(password_private)

		response, err := apiClient.Client.AppRepositoriesAPI.RepositoryImport(context.TODO()).ImportRepoCommand(*body_private).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(response, err))
		}
	}

	// Download the details of the repository into our resource
	errDiag := utils.ReadAfterCreateWithRetries(generateResourceTaikunRepositoryReadWithRetries(), ctx, d, meta)
	if errDiag != nil {
		return errDiag
	}

	// Ensure the desired state - Use function (desired state, current_state , id, name, org_name)
	bindUnbindError := ensureDesiredState(should_be_enabled, d.Get("enabled").(bool), d, meta)
	if bindUnbindError != nil {
		return diag.FromErr(bindUnbindError)
	}

	// Wait until apps from private repo are accessible - Artifactory is a slow sloth
	if should_be_private {
		errErr := resourceTaikunPrivateRepositoryWaitForApps(name, ctx, meta)
		if errErr != nil {
			return diag.FromErr(errErr)
		}
	}

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunRepositoryReadWithRetries(), ctx, d, meta)
}

func generateResourceTaikunRepositoryReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunRepositoryRead(true)
}
func generateResourceTaikunRepositoryReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunRepositoryRead(false)
}
func generateResourceTaikunRepositoryRead(withRetries bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)
		repositoryName := d.Get("name").(string)
		organizationName := d.Get("organization_name").(string)
		private := d.Get("private").(bool)
		data, response, err := apiClient.Client.AppRepositoriesAPI.RepositoryAvailableList(context.TODO()).IsPrivate(private).Search(repositoryName).IsPrivate(private).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(response, err))
		}

		// Iterate through data to find the correct Repository
		foundMatch := false
		var rawRepository tkcore.ArtifactRepositoryDto
		for _, repo := range data.GetData() {
			if (repo.GetName() == repositoryName) && (repo.GetOrganizationName() == organizationName) {
				foundMatch = true
				rawRepository = repo
				break
			}
		}

		if !foundMatch {
			if withRetries {
				d.SetId(d.Get("id").(string)) // We need to tell provider that object was created
				return diag.FromErr(fmt.Errorf("could not find the specified repository (name: %s, organization: %s, private: %t)", repositoryName, organizationName, private))
			}
			return nil
		}

		// Load all the found data to the local object
		err = utils.SetResourceDataFromMap(d, flattenTaikunRepository(&rawRepository, private))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(d.Get("id").(string)) // We need to tell provider that object was created

		return nil
	}
}

func flattenTaikunRepository(rawRepository *tkcore.ArtifactRepositoryDto, isPrivate bool) map[string]interface{} {
	// Ignore changes for URL for public repository
	private_url_or_empty := ""
	if isPrivate {
		private_url_or_empty = rawRepository.GetUrl()
	}

	return map[string]interface{}{
		"id":                rawRepository.GetRepositoryId(),
		"id_apprepo":        fmt.Sprint(rawRepository.GetAppRepoId()),
		"name":              rawRepository.GetName(),
		"organization_name": rawRepository.GetOrganizationName(),
		"private":           isPrivate,
		"url":               private_url_or_empty, // Ignore changes for URL for public repository
		"enabled":           rawRepository.GetIsBound(),
	}
}

// Taikun cannot create private repositories outside the default organization of the logged-in user.
func checkOrganization(d *schema.ResourceData, meta interface{}) (bool, error) {
	apiClient := meta.(*tk.Client)
	data, response, err := apiClient.Client.UsersAPI.UsersUserInfo(context.TODO()).Execute()
	if err != nil {
		return false, tk.CreateError(response, err)
	}
	orgnameDefault := data.Data.GetOrganizationName()
	orgnameDeclared := d.Get("organization_name").(string)
	if data.Data.GetOrganizationName() == d.Get("organization_name").(string) {
		return true, nil
	}
	return false, fmt.Errorf("specified organization (%s) does not match user's organization (%s). You cannot create private repositories outside of your default organization", orgnameDeclared, orgnameDefault)
}

// Ensure the state of this repository matches the desired state provided
func ensureDesiredState(enabledNew bool, enabledCurrent bool, d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*tk.Client)
	// Enabled -> Disabled
	if enabledCurrent && !enabledNew {
		body := &tkcore.UnbindAppRepositoryCommand{}
		body.SetIds([]string{d.Get("id").(string)})
		response, err := apiClient.Client.AppRepositoriesAPI.RepositoryUnbind(context.TODO()).UnbindAppRepositoryCommand(*body).Execute()
		if err != nil {
			return tk.CreateError(response, err)
		}
	}
	// Disabled -> Enabled
	if !enabledCurrent && enabledNew {
		body := &tkcore.BindAppRepositoryCommand{FilteringElements: make([]tkcore.FilteringElementDto, 1)}
		unbind_filter := make([]tkcore.FilteringElementDto, 1)
		unbind_filter[0].SetName(d.Get("name").(string))
		unbind_filter[0].SetOrganizationName(d.Get("organization_name").(string))
		body.SetFilteringElements(unbind_filter)
		response, err := apiClient.Client.AppRepositoriesAPI.RepositoryBind(context.TODO()).BindAppRepositoryCommand(*body).Execute()
		if err != nil {
			return tk.CreateError(response, err)
		}
	}

	// Enabled -> Enabled
	// Disabled -> Disabled

	// Update to latest changes - Download the details of the repository into our resource
	err := utils.ReadAfterCreateWithRetries(generateResourceTaikunRepositoryReadWithRetries(), context.TODO(), d, meta)
	if err != nil {
		return fmt.Errorf("update after enable/disable failed")
	}
	return nil
}

// After we create a private repository, the apps inside it are not available immediately. It takes some time (~60s).
func resourceTaikunPrivateRepositoryWaitForApps(repositoryName string, ctx context.Context, meta interface{}) error {
	apiClient := meta.(*tk.Client)

	pendingStates := []string{"pending"}
	targetStates := []string{"finished"}

	// Try to get the apps until timeout - If apps are listable, repository is ready
	createStateConf := &retry.StateChangeConf{
		Pending: pendingStates,
		Target:  targetStates,
		Refresh: func() (interface{}, string, error) {
			data, response, err := apiClient.Client.PackageAPI.PackageList(context.TODO()).IsPrivate(true).FilterBy(repositoryName).Execute()
			if err != nil {
				return nil, "", tk.CreateError(response, err)
			}

			foundMatch := "pending"
			if data.GetTotalCount() > 0 {
				foundMatch = "finished"
			}

			return data, foundMatch, nil
		},
		Timeout:    5 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err := createStateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for repository (%s) to be read: %s", repositoryName, err)
	}

	return nil
}
