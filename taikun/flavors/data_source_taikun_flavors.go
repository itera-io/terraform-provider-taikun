package flavors

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"
)

func DataSourceTaikunFlavors() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve flavors for a given cloud credential.",
		ReadContext: dataSourceTaikunFlavorsRead,
		Schema: map[string]*schema.Schema{
			"cloud_credential_id": {
				Description:      "Cloud credential ID.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: utils.StringIsInt,
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
				ValidateFunc: validation.IntBetween(2, 1000),
			},
			"max_ram": {
				Description:  "Maximal RAM size in GB.",
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      500,
				ValidateFunc: validation.IntBetween(2, 2199),
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

	cloudCredentialID, err := utils.Atoi32(d.Get("cloud_credential_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	apiClient := meta.(*tk.Client)

	startCPU := int32(d.Get("min_cpu").(int))
	endCPU := int32(d.Get("max_cpu").(int))

	var startRAM float64
	var endRAM float64

	startRAM = utils.GibiByteToByte(d.Get("min_ram").(int)) - 100000 // 0.1 promile of a GibiByte. Correcting rounding errors from type conversion.
	endRAM = utils.GibiByteToByte(d.Get("max_ram").(int)) + 100000   // 0.1 promile of a GibiByte. Correcting rounding errors from type conversion.

	sortBy := "name"
	sortDir := "asc"

	prepare := apiClient.Client.CloudCredentialAPI.CloudcredentialsAllFlavors(context.TODO(), cloudCredentialID)
	prepare = prepare.StartCpu(startCPU).EndCpu(endCPU).StartRam(startRAM).EndRam(endRAM).SortBy(sortBy).SortDirection(sortDir)
	var offset int32 = 0

	var flavorDTOs []tkcore.FlavorsListDto
	for {
		response, res, err := prepare.Offset(offset).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}
		flavorDTOs = append(flavorDTOs, response.GetData()...)
		if len(flavorDTOs) == int(response.GetTotalCount()) {
			break
		}
		offset = int32(len(flavorDTOs))
	}

	flavors := flattenDataSourceTaikunFlavors(flavorDTOs)
	if err := d.Set("flavors", flavors); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.I32toa(cloudCredentialID))
	return nil
}

func flattenDataSourceTaikunFlavors(flavorDTOs []tkcore.FlavorsListDto) []map[string]interface{} {
	flavors := make([]map[string]interface{}, len(flavorDTOs))
	flattenFunc := flattenDataSourceTaikunFlavorsGenericItem
	for i, flavorDTO := range flavorDTOs {
		flavors[i] = flattenFunc(&flavorDTO)
	}
	return flavors
}

func flattenDataSourceTaikunFlavorsGenericItem(flavorDTO *tkcore.FlavorsListDto) map[string]interface{} {
	return map[string]interface{}{
		"cpu":  flavorDTO.GetCpu(),
		"name": flavorDTO.GetName(),
		"ram":  utils.ByteToGibiByte(flavorDTO.GetRam()),
	}
}
