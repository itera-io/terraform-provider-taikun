package taikun

import (
	"context"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/cloud_credentials"
)

func resourceTaikunCloudCredentialSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_key_id": {
			Description:  "The AWS access key ID. Required for AWS.",
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("AWS_ACCESS_KEY_ID", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"availability_zone": {
			Description: "The availability zone of the cloud credential. Optional for Openstack. Required for AWS and Azure. See `zone` for GCP.",
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
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
			DefaultFunc:  schema.EnvDefaultFunc("ARM_CLIENT_ID", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"client_secret": {
			Description:  "The Azure client secret. Required for Azure.",
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("ARM_CLIENT_SECRET", nil),
			ValidateFunc: validation.StringIsNotEmpty,
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
			DefaultFunc:  schema.EnvDefaultFunc("ARM_SUBSCRIPTION_ID", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"tenant_id": {
			Description:  "The Azure tenant ID. Required for Azure.",
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("ARM_TENANT_ID", nil),
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
		"zone": {
			Description:  "The zone of the GCP credential.",
			Type:         schema.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
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

		apiClient := meta.(*taikungoclient.Client)
		id, err := atoi32(d.Id())
		type_cc := d.Get("type").(string)
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, err := apiClient.Client.CloudCredentials.CloudCredentialsDashboardList(cloud_credentials.NewCloudCredentialsDashboardListParams().WithV(ApiVersion).WithID(&id), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}

		if strings.Compare(type_cc, "aws") == 0 {
			if len(response.Payload.Amazon) != 1 {
				if withRetries {
					d.SetId(i32toa(id))
					return diag.Errorf(notFoundAfterCreateOrUpdateError)
				}
				return nil
			}
			res := response.GetPayload().Amazon[0]
			err = setResourceDataFromMap(d, flattenTaikunCloudCredentialAWS(res))
			if err != nil {
				return diag.FromErr(err)
			}

		} else if strings.Compare(type_cc, "azure") == 0 {
			if len(response.Payload.Azure) != 1 {
				if withRetries {
					d.SetId(i32toa(id))
					return diag.Errorf(notFoundAfterCreateOrUpdateError)
				}
				return nil
			}
			res := response.GetPayload().Azure[0]
			err = setResourceDataFromMap(d, flattenTaikunCloudCredentialAzure(res))
			if err != nil {
				return diag.FromErr(err)
			}

		} else if strings.Compare(type_cc, "gcp") == 0 {
			if len(response.Payload.Google) != 1 {
				if withRetries {
					d.SetId(i32toa(id))
					return diag.Errorf(notFoundAfterCreateOrUpdateError)
				}
				return nil
			}
			res := response.GetPayload().Google[0]
			err = setResourceDataFromMap(d, flattenTaikunCloudCredentialGCP(res))
			if err != nil {
				return diag.FromErr(err)
			}

		} else {
			if len(response.Payload.Openstack) != 1 {
				if withRetries {
					d.SetId(i32toa(id))
					return diag.Errorf(notFoundAfterCreateOrUpdateError)
				}
				return nil
			}
			res := response.GetPayload().Openstack[0]
			err = setResourceDataFromMap(d, flattenTaikunCloudCredentialOpenStack(res))
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
