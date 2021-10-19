package taikun

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/flavors"
	"github.com/itera-io/taikungoclient/models"
)

func dataSourceTaikunFlavors() *schema.Resource {
	return &schema.Resource{
		Description: "Get the list of flavors for a given cloud credential",
		ReadContext: dataSourceTaikunFlavorsRead,
		Schema: map[string]*schema.Schema{
			"cloud_type": {
				Description: "Cloud type.",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"AWS",
					"Azure",
					"OpenStack",
				}, false),
			},
			"cloud_credential_id": {
				Description:      "Cloud credential ID.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: stringIsInt,
			},
			"min_cpu": {
				Description:  "Minimal CPU count",
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validation.IntBetween(2, 36),
			},
			"max_cpu": {
				Description:  "Maximal CPU count",
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      36,
				ValidateFunc: validation.IntBetween(2, 36),
			},
			"min_ram": {
				Description:  "Minimal RAM size in GB",
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validation.IntBetween(2, 500),
			},
			"max_ram": {
				Description:  "Maximal RAM size in GB",
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      500,
				ValidateFunc: validation.IntBetween(2, 500),
			},
			"flavors": {
				Description: "List of flavors.",
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
		},
	}
}

func dataSourceTaikunFlavorsRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	flavorDTOs, err := dataSourceTaikunFlavorsGetDTOs(data, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	cloudType := data.Get("cloud_type").(string)
	flavors := flattenDataSourceTaikunFlavors(cloudType, flavorDTOs)
	if err := data.Set("flavors", flavors); err != nil {
		return diag.FromErr(err)
	}
	cloudCredentialID, _ := atoi32(data.Get("cloud_credential_id").(string))
	data.SetId(i32toa(cloudCredentialID))
	return nil
}

func dataSourceTaikunFlavorsGetDTOs(data *schema.ResourceData, meta interface{}) (interface{}, error) {
	startCPU := int32(data.Get("min_cpu").(int))
	endCPU := int32(data.Get("max_cpu").(int))
	startRAM := gibiByteToMebiByte(int32(data.Get("min_ram").(int)))
	endRAM := gibiByteToMebiByte(int32(data.Get("max_ram").(int)))
	cloudType := data.Get("cloud_type").(string)
	cloudCredentialID, err := atoi32(data.Get("cloud_credential_id").(string))
	if err != nil {
		return nil, err
	}
	switch cloudType {
	case "AWS":
		params := flavors.NewFlavorsAwsFlavorsParams().WithV(ApiVersion).WithCloudID(cloudCredentialID)
		params = params.WithStartCPU(&startCPU).WithEndCPU(&endCPU).WithStartRAM(&startRAM).WithEndRAM(&endRAM)
		flavorDTOs, err := dataSourceTaikunFlavorsAWSGetDTOs(params, meta)
		return flavorDTOs, err
	case "Azure":
		params := flavors.NewFlavorsAzureFlavorsParams().WithV(ApiVersion).WithCloudID(cloudCredentialID)
		params = params.WithStartCPU(&startCPU).WithEndCPU(&endCPU).WithStartRAM(&startRAM).WithEndRAM(&endRAM)
		flavorDTOs, err := dataSourceTaikunFlavorsAzureGetDTOs(params, meta)
		return flavorDTOs, err
	case "OpenStack":
		params := flavors.NewFlavorsOpenstackFlavorsParams().WithV(ApiVersion).WithCloudID(cloudCredentialID)
		params = params.WithStartCPU(&startCPU).WithEndCPU(&endCPU).WithStartRAM(&startRAM).WithEndRAM(&endRAM)
		flavorDTOs, err := dataSourceTaikunFlavorsOpenStackGetDTOs(params, meta)
		return flavorDTOs, err
	default:
		return nil, fmt.Errorf("%s is not a valid cloud type", cloudType)
	}
}

func dataSourceTaikunFlavorsAWSGetDTOs(params *flavors.FlavorsAwsFlavorsParams, meta interface{}) ([]*models.AwsFlavorListDto, error) {
	apiClient := meta.(*apiClient)
	var flavorDTOs []*models.AwsFlavorListDto
	for {
		response, err := apiClient.client.Flavors.FlavorsAwsFlavors(params, apiClient)
		if err != nil {
			return nil, err
		}
		flavorDTOs = append(flavorDTOs, response.GetPayload().Data...)
		if len(flavorDTOs) == int(response.GetPayload().TotalCount) {
			break
		}
		offset := int32(len(flavorDTOs))
		params = params.WithOffset(&offset)
	}
	return flavorDTOs, nil
}

func dataSourceTaikunFlavorsAzureGetDTOs(params *flavors.FlavorsAzureFlavorsParams, meta interface{}) ([]*models.AzureFlavorListDto, error) {
	apiClient := meta.(*apiClient)
	var flavorDTOs []*models.AzureFlavorListDto
	for {
		response, err := apiClient.client.Flavors.FlavorsAzureFlavors(params, apiClient)
		if err != nil {
			return nil, err
		}
		flavorDTOs = append(flavorDTOs, response.GetPayload().Data...)
		if len(flavorDTOs) == int(response.GetPayload().TotalCount) {
			break
		}
		offset := int32(len(flavorDTOs))
		params = params.WithOffset(&offset)
	}
	return flavorDTOs, nil
}

func dataSourceTaikunFlavorsOpenStackGetDTOs(params *flavors.FlavorsOpenstackFlavorsParams, meta interface{}) ([]*models.OpenstackFlavorListDto, error) {
	apiClient := meta.(*apiClient)
	var flavorDTOs []*models.OpenstackFlavorListDto
	for {
		response, err := apiClient.client.Flavors.FlavorsOpenstackFlavors(params, apiClient)
		if err != nil {
			return nil, err
		}
		flavorDTOs = append(flavorDTOs, response.GetPayload().Data...)
		if len(flavorDTOs) == int(response.GetPayload().TotalCount) {
			break
		}
		offset := int32(len(flavorDTOs))
		params = params.WithOffset(&offset)
	}
	return flavorDTOs, nil
}

func flattenDataSourceTaikunFlavors(cloudType string, rawFlavorDTOs interface{}) []map[string]interface{} {
	switch cloudType {
	case "AWS":
		return flattenDataSourceTaikunFlavorsAWS(rawFlavorDTOs.([]*models.AwsFlavorListDto))
	case "Azure":
		return flattenDataSourceTaikunFlavorsAzure(rawFlavorDTOs.([]*models.AzureFlavorListDto))
	case "OpenStack":
		return flattenDataSourceTaikunFlavorsOpenStack(rawFlavorDTOs.([]*models.OpenstackFlavorListDto))
	default:
		return nil
	}
}

func flattenDataSourceTaikunFlavorsAWS(flavorDTOs []*models.AwsFlavorListDto) []map[string]interface{} {
	flavors := make([]map[string]interface{}, len(flavorDTOs))
	for i, flavorDTO := range flavorDTOs {
		cpu, _ := atoi32(string(flavorDTO.CPU.(json.Number)))
		ram, _ := atoi32(string(flavorDTO.RAM.(json.Number)))
		flavors[i] = map[string]interface{}{
			"cpu":  cpu,
			"name": flavorDTO.Name.(string),
			"ram":  mebiByteToGibiByte(ram),
		}
	}
	return flavors
}

func flattenDataSourceTaikunFlavorsAzure(flavorDTOs []*models.AzureFlavorListDto) []map[string]interface{} {
	flavors := make([]map[string]interface{}, len(flavorDTOs))
	for i, flavorDTO := range flavorDTOs {
		flavors[i] = map[string]interface{}{
			"cpu":  flavorDTO.CPU.(int32),
			"name": flavorDTO.Name.(string),
			"ram":  flavorDTO.RAM.(int32),
		}
	}
	return flavors
}

func flattenDataSourceTaikunFlavorsOpenStack(flavorDTOs []*models.OpenstackFlavorListDto) []map[string]interface{} {
	flavors := make([]map[string]interface{}, len(flavorDTOs))
	for i, flavorDTO := range flavorDTOs {
		flavors[i] = map[string]interface{}{
			"cpu":  int32(flavorDTO.CPU),
			"name": flavorDTO.Name,
			"ram":  int32(flavorDTO.RAM),
		}
	}
	return flavors
}
