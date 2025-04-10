package app_instance

import (
	"context"
	b64 "encoding/base64"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
	"regexp"
	"time"
)

func resourceTaikunAppInstanceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "The ID of the application instance.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description: "The name of the application instance.",
			Type:        schema.TypeString,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 30),
				validation.StringMatch(
					regexp.MustCompile("^[a-z0-9-]+$"),
					"Application Instance name must contain only alpha numeric characters or non alpha numeric (-)",
				),
			),
			Required: true,
			ForceNew: true,
		},
		"namespace": {
			Description: "Namespace where the application will be deployed.",
			Type:        schema.TypeString,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 30),
				validation.StringMatch(
					regexp.MustCompile("^[a-z0-9-]+$"),
					"Application instance name must contain only alpha numeric characters or non alpha numeric (-)",
				),
			),
			Required: true,
			ForceNew: true,
		},
		"project_id": {
			Description:      "The ID of the project where the application should be deployed.",
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: utils.StringIsInt,
		},
		"catalog_app_id": {
			Description:      "The ID of the catalog app from which we deploy the application.",
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: utils.StringIsInt,
		},
		"parameters_yaml": {
			Description:      "A path to a valid yaml file that includes the parameters for the application.",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: utils.StringIsFilePath,
			StateFunc: func(filePath interface{}) string {
				// Read file contents, encode in base64 and save to state
				paramsEncoded, err := utils.FilePathToBase64String(filePath.(string))
				if err != nil {
					panic(fmt.Errorf("Error reading file %s\nError: %s\n", filePath, err.Error()))
				}
				return paramsEncoded
			},
			ConflictsWith: []string{"parameters_base64"},
		},
		"parameters_base64": {
			Description:   "A base64 encoded file containing parameters for the application.",
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"parameters_yaml"},
		},
		"autosync": {
			Description: "Indicates whether enable or disable autosyc.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"status": {
			Description: "Do not set. Used for tracking application instance failures.",
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     "",
		},
	}
}

func ResourceTaikunAppInstance() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Application Instance Configuration.",
		CreateContext: resourceTaikunAppInstanceCreate,
		ReadContext:   generateResourceTaikunAppInstanceReadWithoutRetries(),
		UpdateContext: resourceTaikunAppInstanceUpdate,
		DeleteContext: resourceTaikunAppInstanceDelete,
		Schema:        resourceTaikunAppInstanceSchema(),
	}
}

func resourceTaikunAppInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	// Prepare arguments
	projectId, err := utils.Atoi32(d.Get("project_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	catalogAppId, err := utils.Atoi32(d.Get("catalog_app_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	extraValues := ""
	paramsInFile := paramsSpecifiedAsFile(d)
	if paramsInFile {
		extraValues, err = utils.FilePathToBase64String(d.Get("parameters_yaml").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		extraValues = d.Get("parameters_base64").(string)
	}

	// Send install
	body := &tkcore.CreateProjectAppCommand{}
	body.SetName(d.Get("name").(string))
	body.SetProjectId(projectId)
	body.SetCatalogAppId(catalogAppId)
	body.SetNamespace(d.Get("namespace").(string))
	body.SetExtraValues(extraValues)
	body.SetAutoSync(d.Get("autosync").(bool))
	data, response, err := apiClient.Client.ProjectAppsAPI.ProjectappInstall(context.TODO()).CreateProjectAppCommand(*body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(response, err))
	}

	// Wait for install to finish
	d.SetId(data.GetId())
	err = resourceTaikunAppInstanceWaitForReady(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunAppInstanceReadWithRetries(), ctx, d, meta)
}

func resourceTaikunAppInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	appInstanceId, err := utils.Atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	response, err := apiClient.Client.ProjectAppsAPI.ProjectappDelete(context.TODO(), appInstanceId).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(response, err))
	}

	// Wait for uninstall
	err = resourceTaikunAppInstanceWaitForDelete(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func generateResourceTaikunAppInstanceReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunAppInstanceRead(true)
}
func generateResourceTaikunAppInstanceReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunAppInstanceRead(false)
}

func generateResourceTaikunAppInstanceRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)
		appId, err := utils.Atoi32(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
		data, _, err := apiClient.Client.ProjectAppsAPI.ProjectappDetails(context.TODO(), appId).Execute()
		if err != nil {
			// Already destroyed/create again
			d.SetId("")
			return nil
		}

		// Application was found in Failed state.
		// Delete application and create it again.
		if data.GetStatus() == tkcore.EINSTANCESTATUS_FAILURE {
			err = d.Set("status", "Failed")
			if err != nil {
				return diag.FromErr(err)
			}
			return nil
		}

		// Load all the found data to the local object
		paramsInFile := paramsSpecifiedAsFile(d)
		if paramsInFile {
			// File parameters were used
			err = utils.SetResourceDataFromMap(d, flattenTaikunAppInstance(true, data))
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			// Base64 parameters were used
			err = utils.SetResourceDataFromMap(d, flattenTaikunAppInstance(false, data))
			if err != nil {
				return diag.FromErr(err)
			}
		}

		// We need to tell provider that object was created
		d.SetId(d.Get("id").(string))
		return nil
	}
}

func flattenTaikunAppInstance(paramsSpecifiedAsFile bool, rawAppInstance *tkcore.ProjectAppDetailsDto) map[string]interface{} {
	paramsKey := "parameters_base64"
	if paramsSpecifiedAsFile {
		paramsKey = "parameters_yaml"
	}
	return map[string]interface{}{
		"id":             utils.I32toa(rawAppInstance.GetId()),
		"name":           rawAppInstance.GetName(),
		"namespace":      rawAppInstance.GetNamespace(),
		"project_id":     utils.I32toa(rawAppInstance.GetProjectId()),
		"catalog_app_id": utils.I32toa(rawAppInstance.GetCatalogAppId()),
		paramsKey:        b64.URLEncoding.EncodeToString([]byte(rawAppInstance.GetValues())),
		"autosync":       rawAppInstance.GetAutoSync(),
	}
}

func resourceTaikunAppInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	appId, err := utils.Atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Autosync
	autosyncOld, autosyncNew := d.GetChange("autosync")
	if autosyncOld != autosyncNew {
		body := tkcore.AutoSyncManagementCommand{}
		body.SetId(appId)
		body.SetMode(autosyncNew.(string))
		response, errSync := apiClient.Client.ProjectAppsAPI.ProjectappAutosync(context.TODO()).AutoSyncManagementCommand(body).Execute()
		if errSync != nil {
			return diag.FromErr(tk.CreateError(response, errSync))
		}
	}

	err = updateParams(appId, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	return utils.ReadAfterUpdateWithRetries(generateResourceTaikunAppInstanceReadWithRetries(), ctx, d, meta)
}

// Wait until app is Ready
func resourceTaikunAppInstanceWaitForReady(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*tk.Client)
	appId, err := utils.Atoi32(d.Id())
	if err != nil {
		return err
	}

	pendingStates := []string{string(tkcore.EINSTANCESTATUS_NONE), string(tkcore.EINSTANCESTATUS_NOT_READY), string(tkcore.EINSTANCESTATUS_INSTALLING), string(tkcore.EINSTANCESTATUS_UNINSTALLING)}
	targetStates := []string{string(tkcore.EINSTANCESTATUS_READY)}

	// Try to get the instance until timeout - If apps are listable, repository is ready
	createStateConf := &retry.StateChangeConf{
		Pending: pendingStates,
		Target:  targetStates,
		Refresh: func() (interface{}, string, error) {
			data, response, err := apiClient.Client.ProjectAppsAPI.ProjectappDetails(context.TODO(), appId).Execute()
			if err != nil {
				return nil, "", tk.CreateError(response, err)
			}

			return data, string(data.GetStatus()), nil
		},
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = createStateConf.WaitForStateContext(context.TODO())
	if err != nil {
		return fmt.Errorf("error waiting for application (%d) to be ready: %s", appId, err)
	}

	return nil
}

// Wait until app is uninstalled, removed, not found.
func resourceTaikunAppInstanceWaitForDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*tk.Client)
	appId, err := utils.Atoi32(d.Id())
	if err != nil {
		return err
	}

	pendingStates := []string{"present"}
	targetStates := []string{"gone"}

	// Try to get the instance until timeout - If app is not present, it was deleted.
	// If uninstall fails during deletion, use your second chance to send uninstall again - usually it can get us unstuck.
	secondChance := true
	createStateConf := &retry.StateChangeConf{
		Pending: pendingStates,
		Target:  targetStates,
		Refresh: func() (interface{}, string, error) {
			data, response, err := apiClient.Client.ProjectAppsAPI.ProjectappList(context.TODO()).Id(appId).Execute()
			if err != nil {
				return nil, "", tk.CreateError(response, err)
			}

			foundMatch := "present"
			if data.GetTotalCount() == 0 {
				foundMatch = "gone"
			}

			// If it failed, we try to delete again.
			if data.GetTotalCount() == 1 {
				if (data.GetData()[0].GetStatus() == tkcore.EINSTANCESTATUS_FAILURE) && secondChance {
					secondChance = false
					response, err = apiClient.Client.ProjectAppsAPI.ProjectappDelete(context.TODO(), appId).Execute()
					if err != nil {
						return nil, "", tk.CreateError(response, err)
					}
				}
			}

			return data, foundMatch, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = createStateConf.WaitForStateContext(context.TODO())
	if err != nil {
		return fmt.Errorf("error waiting for application (%d) to be ready: %s", appId, err)
	}

	return nil
}

// Set new parameters, sync app, wait until ready
func setParamsAndSyncTaikunAppInstance(appId int32, extraValues string, d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*tk.Client)

	body := tkcore.EditProjectAppExtraValuesCommand{}
	body.SetProjectAppId(appId)
	body.SetExtraValues(extraValues)
	_, response, errParams := apiClient.Client.ProjectAppsAPI.ProjectappUpdateExtraValues(context.TODO()).EditProjectAppExtraValuesCommand(body).Execute()
	if errParams != nil {
		return tk.CreateError(response, errParams)
	}

	bodySync := tkcore.SyncProjectAppCommand{}
	bodySync.SetProjectAppId(appId)
	response, errSync := apiClient.Client.ProjectAppsAPI.ProjectappSync(context.TODO()).SyncProjectAppCommand(bodySync).Execute()
	if errSync != nil {
		return tk.CreateError(response, errSync)
	}
	err := resourceTaikunAppInstanceWaitForReady(d, meta)
	if err != nil {
		return err
	}
	return nil
}

// Update parameters of this app in correct order
// Check if user specified file of base64 string. Then modify App instance in correct order.
func updateParams(appId int32, d *schema.ResourceData, meta interface{}) (err error) {
	var extraValues string

	paramsInFile := paramsSpecifiedAsFile(d)

	oldBase64Parameters, newBase64Parameters := d.GetChange("parameters_base64")
	oldYamlParameters, newYamlParameters := d.GetChange("parameters_yaml")

	if paramsInFile {
		//  Fist delete base64 params, then create params from file
		if oldBase64Parameters != newBase64Parameters {
			err = setParamsAndSyncTaikunAppInstance(appId, newBase64Parameters.(string), d, meta)
			if err != nil {
				return err
			}
		}
		if oldYamlParameters != newYamlParameters {
			extraValues, err = utils.FilePathToBase64String(d.Get("parameters_yaml").(string))
			if err != nil {
				return err
			}
			err = setParamsAndSyncTaikunAppInstance(appId, extraValues, d, meta)
			if err != nil {
				return err
			}
		}
	} else {
		//  Fist delete params from file, then create params from base64
		if oldYamlParameters != newYamlParameters {
			extraValues, err = utils.FilePathToBase64String(d.Get("parameters_yaml").(string))
			if err != nil {
				return err
			}
			err = setParamsAndSyncTaikunAppInstance(appId, extraValues, d, meta)
			if err != nil {
				return err
			}
		}
		if oldBase64Parameters != newBase64Parameters {
			err = setParamsAndSyncTaikunAppInstance(appId, newBase64Parameters.(string), d, meta)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Single function to answer in what format did the user specify the params
func paramsSpecifiedAsFile(d *schema.ResourceData) bool {
	parameters_yaml := d.Get("parameters_yaml")
	return parameters_yaml != ""
}
