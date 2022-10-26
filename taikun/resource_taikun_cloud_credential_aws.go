package taikun

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient"
	"github.com/itera-io/taikungoclient/client/aws"
	"github.com/itera-io/taikungoclient/client/cloud_credentials"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunCloudCredentialAWSSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_key_id": {
			Description:  "The AWS access key ID.",
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("AWS_ACCESS_KEY_ID", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"availability_zones": {
			Description: "The given AWS availability zones for the region.",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"az_count": {
			Description:  "The number of AWS availability zone expected for the region.",
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 3),
			Default:      1,
		},
		"created_by": {
			Description: "The creator of the AWS cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"id": {
			Description: "The ID of the AWS cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"is_default": {
			Description: "Indicates whether the AWS cloud credential is the default one.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"last_modified": {
			Description: "Time and date of last modification.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"last_modified_by": {
			Description: "The last user to have modified the AWS cloud credential.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"lock": {
			Description: "Indicates whether to lock the AWS cloud credential.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"name": {
			Description: "The name of the AWS cloud credential.",
			Type:        schema.TypeString,
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 30),
				validation.StringMatch(
					regexp.MustCompile("^[a-zA-Z0-9-]+$"),
					"expected only alpha numeric characters or non alpha numeric (-)",
				),
			),
		},
		"organization_id": {
			Description:      "The ID of the organization which owns the AWS cloud credential.",
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
		"region": {
			Description: "The AWS region.",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			DefaultFunc: schema.EnvDefaultFunc("AWS_DEFAULT_REGION", nil),
			ValidateFunc: validation.StringInSlice(
				[]string{
					"af-south-1",
					"ap-east-1",
					"ap-northeast-1",
					"ap-northeast-2",
					"ap-northeast-3",
					"ap-south-1",
					"ap-southeast-1",
					"ap-southeast-2",
					"ca-central-1",
					"eu-central-1",
					"eu-north-1",
					"eu-south-1",
					"eu-west-1",
					"eu-west-2",
					"eu-west-3",
					"me-south-1",
					"sa-east-1",
					"us-east-1",
					"us-east-2",
					"us-west-1",
					"us-west-2",
					"cn-north-1",
					"cn-northwest-1",
					"us-gov-east-1",
					"us-gov-west-1",
					"us-iso-east-1",
					"us-iso-west-1",
					"us-isob-east-1",
				},
				false,
			),
		},
		"secret_access_key": {
			Description:  "The AWS secret access key.",
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			DefaultFunc:  schema.EnvDefaultFunc("AWS_SECRET_ACCESS_KEY", nil),
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func resourceTaikunCloudCredentialAWS() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun AWS Cloud Credential",
		CreateContext: resourceTaikunCloudCredentialAWSCreate,
		ReadContext:   generateResourceTaikunCloudCredentialAWSReadWithoutRetries(),
		UpdateContext: resourceTaikunCloudCredentialAWSUpdate,
		DeleteContext: resourceTaikunCloudCredentialDelete,
		Schema:        resourceTaikunCloudCredentialAWSSchema(),
	}
}

func resourceTaikunCloudCredentialAWSCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)

	body := &models.CreateAwsCloudCommand{
		Name:               d.Get("name").(string),
		AwsAccessKeyID:     d.Get("access_key_id").(string),
		AwsSecretAccessKey: d.Get("secret_access_key").(string),
		AwsRegion:          d.Get("region").(string),
	}

	/*
		azCount, err := atoi32(d.Get("az_count").(string))
		if err != nil {
			return diag.FromErr(err)
		} else if azCount < 1 || azCount > 3 {
			return diag.Errorf("The az_count value must be between 1 and 3 inclusive.")
		}
	*/
	body.AzCount = int32(d.Get("az_count").(int))

	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
		body.OrganizationID = organizationId
	}

	params := aws.NewAwsCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.Client.Aws.AwsCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	id, err := atoi32(createResult.Payload.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.Payload.ID)

	if d.Get("lock").(bool) {
		if err := resourceTaikunCloudCredentialAWSLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterCreateWithRetries(generateResourceTaikunCloudCredentialAWSReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunCloudCredentialAWSReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialAWSRead(true)
}
func generateResourceTaikunCloudCredentialAWSReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunCloudCredentialAWSRead(false)
}
func generateResourceTaikunCloudCredentialAWSRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*taikungoclient.Client)
		id, err := atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, err := apiClient.Client.CloudCredentials.CloudCredentialsDashboardList(cloud_credentials.NewCloudCredentialsDashboardListParams().WithV(ApiVersion).WithID(&id), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(response.Payload.Amazon) != 1 {
			if withRetries {
				d.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawCloudCredentialAWS := response.GetPayload().Amazon[0]

		err = setResourceDataFromMap(d, flattenTaikunCloudCredentialAWS(rawCloudCredentialAWS))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunCloudCredentialAWSUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*taikungoclient.Client)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if locked, _ := d.GetChange("lock"); locked.(bool) {
		if err := resourceTaikunCloudCredentialAWSLock(id, false, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("access_key_id", "secret_access_key", "name") {
		updateBody := &models.UpdateAwsCommand{
			ID:                 id,
			Name:               d.Get("name").(string),
			AwsAccessKeyID:     d.Get("access_key_id").(string),
			AwsSecretAccessKey: d.Get("secret_access_key").(string),
		}
		updateParams := aws.NewAwsUpdateParams().WithV(ApiVersion).WithBody(updateBody)
		_, err := apiClient.Client.Aws.AwsUpdate(updateParams, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.Get("lock").(bool) {
		if err := resourceTaikunCloudCredentialAWSLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterUpdateWithRetries(generateResourceTaikunCloudCredentialAWSReadWithRetries(), ctx, d, meta)
}

func flattenTaikunCloudCredentialAWS(rawAWSCredential *models.AmazonCredentialsListDto) map[string]interface{} {
	return map[string]interface{}{
		"created_by":         rawAWSCredential.CreatedBy,
		"id":                 i32toa(rawAWSCredential.ID),
		"lock":               rawAWSCredential.IsLocked,
		"is_default":         rawAWSCredential.IsDefault,
		"last_modified":      rawAWSCredential.LastModified,
		"last_modified_by":   rawAWSCredential.LastModifiedBy,
		"name":               rawAWSCredential.Name,
		"organization_id":    i32toa(rawAWSCredential.OrganizationID),
		"organization_name":  rawAWSCredential.OrganizationName,
		"availability_zones": rawAWSCredential.AvailabilityZones,
		"region":             rawAWSCredential.Region,
	}
}

func resourceTaikunCloudCredentialAWSLock(id int32, lock bool, apiClient *taikungoclient.Client) error {
	body := models.CloudLockManagerCommand{
		ID:   id,
		Mode: getLockMode(lock),
	}
	params := cloud_credentials.NewCloudCredentialsLockManagerParams().WithV(ApiVersion).WithBody(&body)
	_, err := apiClient.Client.CloudCredentials.CloudCredentialsLockManager(params, apiClient)
	return err
}
