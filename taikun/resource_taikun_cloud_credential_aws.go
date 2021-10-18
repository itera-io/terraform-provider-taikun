package taikun

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/aws"
	"github.com/itera-io/taikungoclient/client/cloud_credentials"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunCloudCredentialAWSSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "The id of the AWS cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description: "The name of the AWS cloud credential.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"access_key_id": {
			Description:  "The AWS Access Key ID.",
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("AWS_ACCESS_KEY_ID", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"secret_access_key": {
			Description:  "The AWS Secret Access Key.",
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("AWS_SECRET_ACCESS_KEY", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"availability_zone": {
			Description:  "The AWS availability zone for the region.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"region": {
			Description:  "The AWS region.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			DefaultFunc:  schema.EnvDefaultFunc("AWS_DEFAULT_REGION", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"organization_id": {
			Description:      "The id of the organization which owns the AWS cloud credential.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ForceNew:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"organization_name": {
			Description: "The name of the organization which owns the AWS cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"is_locked": {
			Description: "Indicates whether the AWS cloud credential is locked or not.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"is_default": {
			Description: "Indicates whether the AWS cloud credential is the default one.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"created_by": {
			Description: "The creator of the AWS cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified": {
			Description: "Time of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user who modified the AWS cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func resourceTaikunCloudCredentialAWS() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun AWS Cloud Credential",
		CreateContext: resourceTaikunCloudCredentialAWSCreate,
		ReadContext:   resourceTaikunCloudCredentialAWSRead,
		UpdateContext: resourceTaikunCloudCredentialAWSUpdate,
		DeleteContext: resourceTaikunCloudCredentialDelete,
		Schema:        resourceTaikunCloudCredentialAWSSchema(),
	}
}

func resourceTaikunCloudCredentialAWSCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := &models.CreateAwsCloudCommand{
		Name:                data.Get("name").(string),
		AwsAccessKeyID:      data.Get("access_key_id").(string),
		AwsSecretAccessKey:  data.Get("secret_access_key").(string),
		AwsAvailabilityZone: data.Get("availability_zone").(string),
		AwsRegion:           getAWSRegion(data.Get("region").(string)),
	}

	organizationIDData, organizationIDIsSet := data.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", data.Get("organization_id").(string))
		}
		body.OrganizationID = organizationId
	}

	params := aws.NewAwsCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.client.Aws.AwsCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(createResult.Payload.ID)

	locked := data.Get("is_locked").(bool)
	if locked {
		id, err := atoi32(createResult.Payload.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		lockBody := models.CloudLockManagerCommand{
			ID:   id,
			Mode: getLockMode(locked),
		}
		lockParams := cloud_credentials.NewCloudCredentialsLockManagerParams().WithV(ApiVersion).WithBody(&lockBody)
		_, err = apiClient.client.CloudCredentials.CloudCredentialsLockManager(lockParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaikunCloudCredentialAWSRead(ctx, data, meta)
}

func resourceTaikunCloudCredentialAWSRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	data.SetId("")
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := apiClient.client.CloudCredentials.CloudCredentialsDashboardList(cloud_credentials.NewCloudCredentialsDashboardListParams().WithV(ApiVersion).WithID(&id), apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(response.Payload.Amazon) == 1 {
		rawCloudCredentialAWS := response.GetPayload().Amazon[0]

		if err := data.Set("created_by", rawCloudCredentialAWS.CreatedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("id", i32toa(rawCloudCredentialAWS.ID)); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("is_locked", rawCloudCredentialAWS.IsLocked); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("is_default", rawCloudCredentialAWS.IsDefault); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("last_modified", rawCloudCredentialAWS.LastModified); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("last_modified_by", rawCloudCredentialAWS.LastModifiedBy); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("name", rawCloudCredentialAWS.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("availability_zone", rawCloudCredentialAWS.AvailabilityZone); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("region", rawCloudCredentialAWS.Region); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_id", i32toa(rawCloudCredentialAWS.OrganizationID)); err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("organization_name", rawCloudCredentialAWS.OrganizationName); err != nil {
			return diag.FromErr(err)
		}

		data.SetId(i32toa(id))
	}

	return nil
}

func resourceTaikunCloudCredentialAWSUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if data.HasChange("access_key_id") || data.HasChange("secret_access_key") || data.HasChange("name") {
		updateBody := &models.UpdateAwsCommand{
			ID:                 id,
			Name:               data.Get("name").(string),
			AwsAccessKeyID:     data.Get("access_key_id").(string),
			AwsSecretAccessKey: data.Get("secret_access_key").(string),
		}
		updateParams := aws.NewAwsUpdateParams().WithV(ApiVersion).WithBody(updateBody)
		_, err := apiClient.client.Aws.AwsUpdate(updateParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if data.HasChange("is_locked") {
		lockBody := models.CloudLockManagerCommand{
			ID:   id,
			Mode: getLockMode(data.Get("is_locked").(bool)),
		}
		lockParams := cloud_credentials.NewCloudCredentialsLockManagerParams().WithV(ApiVersion).WithBody(&lockBody)
		_, err = apiClient.client.CloudCredentials.CloudCredentialsLockManager(lockParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaikunCloudCredentialAWSRead(ctx, data, meta)
}
