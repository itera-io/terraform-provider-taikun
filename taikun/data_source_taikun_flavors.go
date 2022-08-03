package taikun

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/cloud_credentials"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunFlavors() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve flavors for a given cloud credential.",
		ReadContext: dataSourceTaikunFlavorsRead,
		Schema: map[string]*schema.Schema{
			"cloud_credential_id": {
				Description:      "Cloud credential ID.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: stringIsInt,
			},
			"flavors": {
				Description: "List of retrieved flavors.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu": {
							Description: "CPU count.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"name": {
							Description: "Flavor name.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"ram": {
							Description: "RAM size.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
			},
			"max_cpu": {
				Description:  "Maximal CPU count.",
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      36,
				ValidateFunc: validation.IntBetween(2, 36),
			},
			"max_ram": {
				Description:  "Maximal RAM size in GB.",
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      500,
				ValidateFunc: validation.IntBetween(2, 500),
			},
			"min_cpu": {
				Description:  "Minimal CPU count.",
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validation.IntBetween(2, 36),
			},
			"min_ram": {
				Description:  "Minimal RAM size in GB.",
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validation.IntBetween(2, 500),
			},
		},
	}
}

func dataSourceTaikunFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	cloudCredentialID, err := atoi32(d.Get("cloud_credential_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	startCPU := int32(d.Get("min_cpu").(int))
	endCPU := int32(d.Get("max_cpu").(int))
	// startRAM := gibiByteToMebiByte(int32(d.Get("min_ram").(int)))
	startRAM := float64(gibiByteToMebiByte(int32(d.Get("min_ram").(int)))) // TODO check validity of conversion
	// endRAM := gibiByteToMebiByte(int32(d.Get("max_ram").(int)))
	endRAM := float64(gibiByteToMebiByte(int32(d.Get("max_ram").(int)))) // TODO check validity of conversion
	sortBy := "name"
	sortDir := "asc"

	params := cloud_credentials.NewCloudCredentialsAllFlavorsParams().WithV(ApiVersion).WithCloudID(cloudCredentialID)
	params = params.WithStartCPU(&startCPU).WithEndCPU(&endCPU).WithStartRAM(&startRAM).WithEndRAM(&endRAM)
	params = params.WithSortBy(&sortBy).WithSortDirection(&sortDir)

	apiClient := meta.(*apiClient)
	var cloudType string
	var flavorDTOs []*models.FlavorsListDto
	for {
		response, err := apiClient.client.CloudCredentials.CloudCredentialsAllFlavors(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		cloudType = response.Payload.CloudType
		flavorDTOs = append(flavorDTOs, response.Payload.Data...)
		if len(flavorDTOs) == int(response.Payload.TotalCount) {
			break
		}
		offset := int32(len(flavorDTOs))
		params = params.WithOffset(&offset)
	}

	flavors := flattenDataSourceTaikunFlavors(cloudType, flavorDTOs)
	if err := d.Set("flavors", flavors); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(i32toa(cloudCredentialID))
	return nil
}

func flattenDataSourceTaikunFlavors(cloudType string, flavorDTOs []*models.FlavorsListDto) []map[string]interface{} {
	flavors := make([]map[string]interface{}, len(flavorDTOs))
	flattenFunc := getFlattenDataSourceTaikunFlavorsItemFunc(cloudType)
	for i, flavorDTO := range flavorDTOs {
		flavors[i] = flattenFunc(flavorDTO)
	}
	return flavors
}

type flattenDataSourceTaikunFlavorsItemFunc func(flavorDTO *models.FlavorsListDto) map[string]interface{}

func getFlattenDataSourceTaikunFlavorsItemFunc(cloudType string) flattenDataSourceTaikunFlavorsItemFunc {
	switch strings.ToLower(cloudType) {
	case "aws":
		return flattenDataSourceTaikunFlavorsAWSItem
	case "azure":
		return flattenDataSourceTaikunFlavorsAzureItem
	case "openstack":
		return flattenDataSourceTaikunFlavorsOpenStackItem
	default:
		return nil
	}
}

func flattenDataSourceTaikunFlavorsAWSItem(flavorDTO *models.FlavorsListDto) map[string]interface{} {
	return map[string]interface{}{
		"cpu":  flavorDTO.CPU,
		"name": flavorDTO.Name,
		"ram":  mebiByteToGibiByte(flavorDTO.RAM),
	}
}

func flattenDataSourceTaikunFlavorsAzureItem(flavorDTO *models.FlavorsListDto) map[string]interface{} {
	return map[string]interface{}{
		"cpu":  flavorDTO.CPU,
		"name": flavorDTO.Name,
		"ram":  flavorDTO.RAM,
	}
}

func flattenDataSourceTaikunFlavorsOpenStackItem(flavorDTO *models.FlavorsListDto) map[string]interface{} {
	return map[string]interface{}{
		"cpu":  flavorDTO.CPU,
		"name": flavorDTO.Name,
		"ram":  mebiByteToGibiByte(flavorDTO.RAM),
	}
}
