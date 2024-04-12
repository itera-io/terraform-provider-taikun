package utils

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// DataSourceSchemaFromResourceSchema is a recursive func that
// converts an existing Resource schema to a Datasource schema.
// All schema elements are copied, but certain attributes are ignored or changed:
// - all attributes have Computed = true
// - all attributes have ForceNew, Required = false
// - Validation funcs and attributes (e.g. MaxItems) are not copied
func DataSourceSchemaFromResourceSchema(rs map[string]*schema.Schema) map[string]*schema.Schema {
	ds := make(map[string]*schema.Schema, len(rs))
	for k, v := range rs {
		dv := &schema.Schema{
			Computed:    true,
			Sensitive:   v.Sensitive,
			Description: v.Description,
			Type:        v.Type,
		}

		switch v.Type {
		case schema.TypeSet:
			dv.Set = v.Set
			fallthrough
		case schema.TypeList:
			// List & Set types are generally used for 2 cases:
			// - a list/set of simple primitive values (e.g. list of strings)
			// - a sub resource
			if elem, ok := v.Elem.(*schema.Resource); ok {
				// handle the case where the Element is a sub-resource
				dv.Elem = &schema.Resource{
					Schema: DataSourceSchemaFromResourceSchema(elem.Schema),
				}
			} else {
				// handle simple primitive case
				dv.Elem = v.Elem
			}

		default:
			// Elem of all other types are copied as-is
			dv.Elem = v.Elem

		}
		ds[k] = dv

	}
	return ds
}

// fixDatasourceSchemaFlags is a convenience func that toggles the Computed,
// Optional + Required flags on a schema element. This is useful when the schema
// has been generated (using `datasourceSchemaFromResourceSchema` above for
// example) and therefore the attribute flags were not set appropriately when
// first added to the schema definition. Currently only supports top-level
// schema elements.
func fixDatasourceSchemaFlags(schema map[string]*schema.Schema, required bool, keys ...string) {
	for _, v := range keys {
		schema[v].Computed = false
		schema[v].Optional = !required
		schema[v].Required = required
	}
}

func RemoveForceNewsFromSchema(schema map[string]*schema.Schema) {
	for _, v := range schema {
		v.ForceNew = false
	}
}

func SetFieldInSchema(schema map[string]*schema.Schema, key string, value *schema.Schema) {
	schema[key] = value
}

func AddRequiredFieldsToSchema(schema map[string]*schema.Schema, keys ...string) {
	fixDatasourceSchemaFlags(schema, true, keys...)
}

func AddOptionalFieldsToSchema(schema map[string]*schema.Schema, keys ...string) {
	fixDatasourceSchemaFlags(schema, false, keys...)
}

func SetValidateDiagFuncToSchema(schema map[string]*schema.Schema, key string, f schema.SchemaValidateDiagFunc) {
	schema[key].ValidateDiagFunc = f
}

func DeleteFieldsFromSchema(schema map[string]*schema.Schema, keys ...string) {
	for _, v := range keys {
		delete(schema, v)
	}
}
