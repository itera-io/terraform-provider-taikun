package taikun

import (
	"context"
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
						},
						"name": {
							Description: "Flavor name.",
							Type:        schema.TypeString,
						},
						"ram": {
							Description: "RAM size.",
							Type:        schema.TypeInt,
						},
					},
				},
			},
		},
	}
}

func dataSourceTaikunFlavorsRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cloudType := data.Get("cloud_type").(string)
	params := dataSourceTaikunFlavorsMakeParams(cloudType, data)
	flavorDTOs, err := dataSourceTaikunFlavorsGetDTOs(cloudType, params, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	flavors := flattenDataSourceTaikunFlavors(cloudType, flavorDTOs)
	if err := data.Set("flavors", flavors); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(i32toa(data.Get("cloud_credential_id").(int32)))
	return nil
}

func dataSourceTaikunFlavorsMakeParams(cloudType string, data *schema.ResourceData) interface{} {
	var paramsAddr interface{}
	startCPU := data.Get("min_cpu").(int32)
	endCPU := data.Get("max_cpu").(int32)
	startRAM := data.Get("min_ram").(int32)
	endRAM := data.Get("max_ram").(int32)
	switch cloudType {
	case "AWS":
		params := flavors.NewFlavorsAwsFlavorsParams().WithV(ApiVersion).WithCloudID(data.Get("cloud_credential_id").(int32))
		params = params.WithStartCPU(&startCPU).WithEndCPU(&endCPU).WithStartRAM(&startRAM).WithEndRAM(&endRAM)
		paramsAddr = &params
	case "Azure":
		params := flavors.NewFlavorsAzureFlavorsParams().WithV(ApiVersion).WithCloudID(data.Get("cloud_credential_id").(int32))
		params = params.WithStartCPU(&startCPU).WithEndCPU(&endCPU).WithStartRAM(&startRAM).WithEndRAM(&endRAM)
		paramsAddr = &params
	case "OpenStack":
		params := flavors.NewFlavorsOpenstackFlavorsParams().WithV(ApiVersion).WithCloudID(data.Get("cloud_credential_id").(int32))
		params = params.WithStartCPU(&startCPU).WithEndCPU(&endCPU).WithStartRAM(&startRAM).WithEndRAM(&endRAM)
		paramsAddr = &params
	}
	return paramsAddr
}

func dataSourceTaikunFlavorsGetDTOs(cloudType string, params interface{}, meta interface{}) (interface{}, error) {
	switch cloudType {
	case "AWS":
		flavorDTOs, err := dataSourceTaikunFlavorsAWSGetDTOs(params.(*flavors.FlavorsAwsFlavorsParams), meta)
		return &flavorDTOs, err
	case "Azure":
		flavorDTOs, err := dataSourceTaikunFlavorsAzureGetDTOs(params.(*flavors.FlavorsAzureFlavorsParams), meta)
		return &flavorDTOs, err
	case "OpenStack":
		flavorDTOs, err := dataSourceTaikunFlavorsOpenStackGetDTOs(params.(*flavors.FlavorsOpenstackFlavorsParams), meta)
		return &flavorDTOs, err
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
	flavorDTOs := rawFlavorDTOs.([]interface{})
	flavors := make([]map[string]interface{}, len(flavorDTOs))
	flattenFunc := getFlattenDataSourceTaikunFlavorItemFunc(cloudType)
	for i, flavorDTO := range flavorDTOs {
		flavors[i] = flattenFunc(flavorDTO)
	}
	return flavors
}

type flattenDataSourceTaikunFlavorItemFunc func(interface{}) map[string]interface{}

func getFlattenDataSourceTaikunFlavorItemFunc(cloudType string) flattenDataSourceTaikunFlavorItemFunc {
	switch cloudType {
	case "AWS":
		return flattenDataSourceTaikunFlavorsAWSItem
	case "Azure":
		return flattenDataSourceTaikunFlavorsAzureItem
	case "OpenStack":
		return flattenDataSourceTaikunFlavorsOpenStackItem
	default:
		return nil
	}
}

func flattenDataSourceTaikunFlavorsAWSItem(rawFlavorDTO interface{}) map[string]interface{} {
	flavorDTO := rawFlavorDTO.(*models.AwsFlavorListDto)
	return map[string]interface{}{
		"cpu":  flavorDTO.CPU.(int32),
		"name": flavorDTO.Name.(string),
		"ram":  flavorDTO.RAM.(int32),
	}
}

func flattenDataSourceTaikunFlavorsAzureItem(rawFlavorDTO interface{}) map[string]interface{} {
	flavorDTO := rawFlavorDTO.(*models.AzureFlavorListDto)
	return map[string]interface{}{
		"cpu":  flavorDTO.CPU.(int32),
		"name": flavorDTO.Name.(string),
		"ram":  flavorDTO.RAM.(int32),
	}
}

func flattenDataSourceTaikunFlavorsOpenStackItem(rawFlavorDTO interface{}) map[string]interface{} {
	flavorDTO := rawFlavorDTO.(*models.OpenstackFlavorListDto)
	return map[string]interface{}{
		"cpu":  int32(flavorDTO.CPU),
		"name": flavorDTO.Name,
		"ram":  int32(flavorDTO.RAM),
	}
}
