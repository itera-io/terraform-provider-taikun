package repository

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
)

func dataSourceTaikunRepositorySchema() map[string]*schema.Schema {
	dsSchema := utils.DataSourceSchemaFromResourceSchema(resourceTaikunRepositorySchema())
	utils.AddRequiredFieldsToSchema(dsSchema, "name")
	utils.SetValidateDiagFuncToSchema(dsSchema, "name", utils.StringLenBetween(3, 30))
	utils.AddRequiredFieldsToSchema(dsSchema, "organization_name")
	utils.SetValidateDiagFuncToSchema(dsSchema, "organization_name", utils.StringLenBetween(3, 30))
	utils.AddOptionalFieldsToSchema(dsSchema, "organization_id")
	utils.SetValidateDiagFuncToSchema(dsSchema, "organization_id", utils.StringIsInt)
	utils.AddRequiredFieldsToSchema(dsSchema, "private")
	dsSchema["private"].Type = schema.TypeBool
	return dsSchema
}

func DataSourceTaikunRespository() *schema.Resource {
	return &schema.Resource{
		Description: "Get an Repository by its ID.",
		ReadContext: dataSourceTaikunRepositoryRead,
		Schema:      dataSourceTaikunRepositorySchema(),
	}
}

func dataSourceTaikunRepositoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := d.Set("name", d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("organization_name", d.Get("organization_name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("private", d.Get("private").(bool))
	if err != nil {
		return diag.FromErr(err)
	}

	return generateResourceTaikunRepositoryReadWithoutRetries()(ctx, d, meta)
}
