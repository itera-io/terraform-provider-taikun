package cc_aws

import (
	"context"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTaikunCloudCredentialAWSSchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunCloudCredentialAWSSchema())
	utils.AddRequiredFieldsToSchema(dsSchema, "id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "id", utils.StringIsInt)
	utils.DeleteFieldsFromSchema(dsSchema, "secret_access_key", "access_key_id")
	return dsSchema
}

func DataSourceTaikunCloudCredentialAWS() *schema.Resource {
	return &schema.Resource{
		Description: "Get an AWS cloud credential by its ID.",
		ReadContext: dataSourceTaikunCloudCredentialAWSRead,
		Schema:      dataSourceTaikunCloudCredentialAWSSchema(),
	}
}

func dataSourceTaikunCloudCredentialAWSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("id").(string))

	return generateResourceTaikunCloudCredentialAWSReadWithoutRetries()(ctx, d, meta)
}
