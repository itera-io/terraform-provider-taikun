package taikun

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// DEPRECATED: this data source is deprecated in favour of `taikun_cloud_credential_aws`, `taikun_cloud_credential_azure`, `taikun_cloud_credential_gcp` and `taikun_cloud_credential_openstack`...

func resourceTaikunCloudCredentialSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"az_count": {
			Description:  "The number of availability zone expected for the region/location. Required for AWS, Azure and GCP.",
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 3),
			Default:      1,
		},
		"access_key_id": {
			Description:  "The AWS access key ID. Required for AWS.",
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("AWS_ACCESS_KEY_ID", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"availability_zone": {
			Description: "The availability zone of the cloud credential. Optional for Openstack.",
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
		},
		"availability_zones": {
			Description: "The given AWS/Azure/GCP availability zones for the region/location.",
			Type:        schema.TypeSet,
			Computed:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"billing_account_id": {
			Description:   "The ID of the GCP credential's billing account.",
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"import_project"},
		},
		"billing_account_name": {
			Description: "The name of the GCP credential's billing account.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"config_file": {
			Description:      "The path of the GCP credential's configuration file. Required for GCP.",
			Type:             schema.TypeString,
			Optional:         true,
			ForceNew:         true,
			ValidateDiagFunc: stringIsFilePath,
		},
		"client_id": {
			Description:  "The Azure client ID. Required for Azure.",
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("AZURE_CLIENT_ID", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"client_secret": {
			Description:  "The Azure client secret. Required for Azure.",
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("AZURE_SECRET", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"continent": {
			Description: "The OpenStack continent (`Asia`, `Europe` or `America`).",
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
			DefaultFunc: schema.EnvDefaultFunc("OS_CONTINENT", nil),
			ValidateFunc: validation.StringInSlice([]string{
				"Asia",
				"Europe",
				"America",
			}, false),
		},
		"created_by": {
			Description: "The creator of the cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"domain": {
			Description:  "The OpenStack domain. Required for Openstack.",
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_USER_DOMAIN_NAME", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"folder_id": {
			Description:   "The folder ID of the GCP credential.",
			Optional:      true,
			Type:          schema.TypeString,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"import_project"},
		},
		"id": {
			Description: "The ID of the cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"imported_network_subnet_id": {
			Description: "The OpenStack network subnet ID to import a network.",
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
		},
		"is_default": {
			Description: "Indicates whether the cloud credential is the default one.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"import_project": {
			Description:   "Whether to import a project or not",
			Type:          schema.TypeBool,
			Default:       false,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"billing_account_id", "folder_id"},
		},
		"last_modified": {
			Description: "Time and date of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user to have modified the cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"location": {
			Description:  "The Azure location. Required for Azure.",
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"lock": {
			Description: "Indicates whether to lock the cloud credential.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"name": {
			Description: "The name of the cloud credential.",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 30),
				validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z0-9-]+$"),
					"expected only alpha numeric characters or '-'",
				),
			),
		},
		"organization_id": {
			Description:      "The ID of the organization which owns the cloud credential.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"organization_name": {
			Description: "The name of the organization which owns the cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"password": {
			Description:  "The OpenStack password. Required for Openstack.",
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_PASSWORD", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"project_id": {
			Description: "The OpenStack project ID. Required for Openstack.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"project_name": {
			Description:  "The OpenStack project name. Required for Openstack.",
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_PROJECT_NAME", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"public_network_name": {
			Description:  "The name of the public OpenStack network to use. Required for Openstack.",
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_INTERFACE", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"region": {
			Description:  "The region of the cloud credential. Required for Openstack, AWS and GCP.",
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"secret_access_key": {
			Description:  "The AWS secret access key. Required for AWS.",
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("AWS_SECRET_ACCESS_KEY", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"subscription_id": {
			Description:  "The Azure subscription ID. Required for Azure.",
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("AZURE_SUBSCRIPTION", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"tenant_id": {
			Description:  "The Azure tenant ID. Required for Azure.",
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("AZURE_TENANT", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"type": {
			Description: "The type of the cloud credential.",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"aws",
				"azure",
				"gcp",
				"openstack",
			}, false),
		},
		"url": {
			Description:  "The OpenStack authentication URL. Required for Openstack.",
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_AUTH_URL", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"user": {
			Description:  "The OpenStack user. Required for Openstack.",
			Type:         schema.TypeString,
			Optional:     true,
			DefaultFunc:  schema.EnvDefaultFunc("OS_USERNAME", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"volume_type_name": {
			Description: "The OpenStack type of volume.",
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
		},
		"zones": {
			Description: "The given zones of the GCP credential.",
			Type:        schema.TypeSet,
			Computed:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func resourceTaikunCloudCredential() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Cloud Credential",
		CreateContext: resourceTaikunCloudCredentialCreate,
		ReadContext:   generateResourceTaikunCloudCredentialReadWithoutRetries(),
		UpdateContext: resourceTaikunCloudCredentialUpdate,
		DeleteContext: resourceTaikunCloudCredentialDelete,
		Schema:        resourceTaikunCloudCredentialSchema(),
	}
}

func resourceTaikunCloudCredentialCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	type_cc := d.Get("type").(string)

	if strings.Compare(type_cc, "aws") == 0 {
		return resourceTaikunCloudCredentialAWSCreate(ctx, d, meta)
	} else if strings.Compare(type_cc, "azure") == 0 {
		return resourceTaikunCloudCredentialAzureCreate(ctx, d, meta)
	} else if strings.Compare(type_cc, "gcp") == 0 {
		return resourceTaikunCloudCredentialGCPCreate(ctx, d, meta)
	} else if strings.Compare(type_cc, "openstack") == 0 {
		return resourceTaikunCloudCredentialOpenStackCreate(ctx, d, meta)
	} else {
		return diag.Errorf("Invalid kind of cloud provider")
	}
}

func generateResourceTaikunCloudCredentialReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialRead(false)
}

func generateResourceTaikunCloudCredentialRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

		apiClient := meta.(*tk.Client)
		id, err := atoi32(d.Id())
		type_cc := d.Get("type").(string)
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, res, err := apiClient.Client.CloudCredentialAPI.CloudcredentialsDashboardList(context.TODO()).Id(id).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}

		if strings.Compare(type_cc, "aws") == 0 {
			if len(response.GetAmazon()) != 1 {
				if withRetries {
					d.SetId(i32toa(id))
					return diag.Errorf(notFoundAfterCreateOrUpdateError)
				}
				return nil
			}
			res := response.GetAmazon()[0]
			err = setResourceDataFromMap(d, flattenTaikunCloudCredentialAWS(&res))
			if err != nil {
				return diag.FromErr(err)
			}

		} else if strings.Compare(type_cc, "azure") == 0 {
			if len(response.GetAzure()) != 1 {
				if withRetries {
					d.SetId(i32toa(id))
					return diag.Errorf(notFoundAfterCreateOrUpdateError)
				}
				return nil
			}
			res := response.GetAzure()[0]
			err = setResourceDataFromMap(d, flattenTaikunCloudCredentialAzure(&res))
			if err != nil {
				return diag.FromErr(err)
			}

		} else if strings.Compare(type_cc, "gcp") == 0 {
			if len(response.GetGoogle()) != 1 {
				if withRetries {
					d.SetId(i32toa(id))
					return diag.Errorf(notFoundAfterCreateOrUpdateError)
				}
				return nil
			}
			res := response.GetGoogle()[0]
			err = setResourceDataFromMap(d, flattenTaikunCloudCredentialGCP(&res))
			if err != nil {
				return diag.FromErr(err)
			}

		} else {
			if len(response.GetOpenstack()) != 1 {
				if withRetries {
					d.SetId(i32toa(id))
					return diag.Errorf(notFoundAfterCreateOrUpdateError)
				}
				return nil
			}
			res := response.GetOpenstack()[0]
			err = setResourceDataFromMap(d, flattenTaikunCloudCredentialOpenStack(&res))
			if err != nil {
				return diag.FromErr(err)
			}
		}

		d.SetId(i32toa(id))

		return nil

	}

}

func resourceTaikunCloudCredentialUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	type_cc := d.Get("type").(string)

	if strings.Compare(type_cc, "aws") == 0 {
		return resourceTaikunCloudCredentialAWSUpdate(ctx, d, meta)
	} else if strings.Compare(type_cc, "azure") == 0 {
		return resourceTaikunCloudCredentialAzureUpdate(ctx, d, meta)
	} else if strings.Compare(type_cc, "gcp") == 0 {
		return resourceTaikunCloudCredentialGCPUpdate(ctx, d, meta)
	} else if strings.Compare(type_cc, "openstack") == 0 {
		return resourceTaikunCloudCredentialOpenStackUpdate(ctx, d, meta)
	} else {
		return diag.Errorf("Invalid kind of cloud credential")
	}
}
