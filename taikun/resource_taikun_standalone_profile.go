package taikun

import (
	"context"
	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTaikunStandaloneProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "The ID of the standalone profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"lock": {
			Description: "Indicates whether to lock the standalone profile.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"name": {
			Description: "The name of the standalone profile.",
			Type:        schema.TypeString,
			Required:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(3, 30),
				validation.StringMatch(
					regexp.MustCompile("^[a-z0-9_-]*$"),
					"expected a lowercase string",
				),
			),
		},
		"organization_id": {
			Description:      "The ID of the organization which owns the standalone profile.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ForceNew:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"organization_name": {
			Description: "The name of the organization which owns the standalone profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"public_key": {
			Description:  "The public key of the standalone profile.",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"security_group": {
			Description: "List of security groups.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"cidr": {
						Description:  "Remote IP prefix.",
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.IsCIDR,
					},
					"from_port": {
						Description:  "Min range port.",
						Type:         schema.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"id": {
						Description: "ID of the security group.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"ip_protocol": {
						Description:  "IP Protocol: `TCP`, `UDP` or `ICMP`.",
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP", "ICMP"}, false),
					},
					"name": {
						Description: "Name of the security group.",
						Type:        schema.TypeString,
						Required:    true,
						ValidateFunc: validation.All(
							validation.StringLenBetween(3, 30),
							validation.StringMatch(
								regexp.MustCompile("^[a-z0-9_-]*$"),
								"expected a lowercase string",
							),
						),
					},
					"to_port": {
						Description:  "Max range port.",
						Type:         schema.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
				},
			},
		},
	}
}

func resourceTaikunStandaloneProfile() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Standalone Profile",
		CreateContext: resourceTaikunStandaloneProfileCreate,
		ReadContext:   generateResourceTaikunStandaloneProfileReadWithoutRetries(),
		UpdateContext: resourceTaikunStandaloneProfileUpdate,
		DeleteContext: resourceTaikunStandaloneProfileDelete,
		Schema:        resourceTaikunStandaloneProfileSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunStandaloneProfileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	body := tkcore.StandAloneProfileCreateCommand{}
	body.SetName(d.Get("name").(string))
	body.SetPublicKey(d.Get("public_key").(string))

	if securityGroups, isSecurityGroupsSet := d.GetOk("security_group"); isSecurityGroupsSet {
		rawSecurityGroupList := securityGroups.([]interface{})
		securityGroupList := make([]tkcore.StandAloneProfileSecurityGroupDto, len(rawSecurityGroupList))
		for i, e := range rawSecurityGroupList {
			rawSecurityGroup := e.(map[string]interface{})
			securityGroupList[i] = tkcore.StandAloneProfileSecurityGroupDto{}
			securityGroupList[i].SetName(rawSecurityGroup["name"].(string))
			securityGroupList[i].SetPortMaxRange(int32(rawSecurityGroup["to_port"].(int)))
			securityGroupList[i].SetPortMinRange(int32(rawSecurityGroup["from_port"].(int)))
			securityGroupList[i].SetProtocol(getSecurityGroupProtocol(rawSecurityGroup["ip_protocol"].(string)))
			securityGroupList[i].SetRemoteIpPrefix(rawSecurityGroup["cidr"].(string))
		}
		body.SecurityGroups = securityGroupList
	}

	organizationIDData, organizationIDIsSet := d.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", d.Get("organization_id").(string))
		}
		body.SetOrganizationId(organizationId)
	}

	createResult, _, err := apiClient.Client.StandaloneProfileAPI.StandaloneprofileCreate(ctx).StandAloneProfileCreateCommand(body).Execute()
	if err != nil {
		return diag.FromErr(err)
	}
	id, err := atoi32(createResult.GetId())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.GetId())

	if d.Get("lock").(bool) {
		if err := resourceTaikunStandaloneProfileLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterCreateWithRetries(generateResourceTaikunStandaloneProfileReadWithRetries(), ctx, d, meta)
}
func generateResourceTaikunStandaloneProfileReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunStandaloneProfileRead(true)
}
func generateResourceTaikunStandaloneProfileReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunStandaloneProfileRead(false)
}
func generateResourceTaikunStandaloneProfileRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)
		id, err := atoi32(d.Id())
		d.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, _, err := apiClient.Client.StandaloneProfileAPI.StandaloneprofileList(context.TODO()).Id(id).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
		if len(response.GetData()) != 1 {
			if withRetries {
				d.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawStandaloneProfile := response.GetData()[0]

		securityGroupResponse, _, err := apiClient.Client.SecurityGroupAPI.SecuritygroupList(context.TODO(), id).Execute()
		if err != nil {

			/*
					if _, ok := err.(*security_group.SecurityGroupListNotFound); ok && withRetries {
					d.SetId(i32toa(id))
					return diag.Errorf(notFoundAfterCreateOrUpdateError)
				}

			*/
			return diag.FromErr(err)
		}

		err = setResourceDataFromMap(d, flattenTaikunStandaloneProfile(&rawStandaloneProfile, securityGroupResponse))
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunStandaloneProfileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("name") {
		body := tkcore.StandAloneProfileUpdateCommand{}
		body.SetId(id)
		body.SetName(d.Get("name").(string))

		_, err := apiClient.Client.StandaloneProfileAPI.StandaloneprofileEdit(ctx).StandAloneProfileUpdateCommand(body).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("lock") {
		if err := resourceTaikunStandaloneProfileLock(id, d.Get("lock").(bool), apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("security_group") {
		old, new := d.GetChange("security_group")

		// Delete
		oldSecurityGroupList := old.([]interface{})
		for _, e := range oldSecurityGroupList {
			rawSecurityGroup := e.(map[string]interface{})
			secId, err := atoi32(rawSecurityGroup["id"].(string))
			if err != nil {
				return diag.FromErr(err)
			}
			_, err = apiClient.Client.SecurityGroupAPI.SecuritygroupDelete(ctx, secId).Execute()
			if err != nil {
				return diag.FromErr(err)
			}
		}

		// Add
		newSecurityGroupList := new.([]interface{})
		for _, e := range newSecurityGroupList {
			rawSecurityGroup := e.(map[string]interface{})
			body := tkcore.CreateSecurityGroupCommand{}
			body.SetName(rawSecurityGroup["name"].(string))
			body.SetPortMaxRange(int32(rawSecurityGroup["to_port"].(int)))
			body.SetPortMinRange(int32(rawSecurityGroup["from_port"].(int)))
			body.SetProtocol(getSecurityGroupProtocol(rawSecurityGroup["ip_protocol"].(string)))
			body.SetRemoteIpPrefix(rawSecurityGroup["cidr"].(string))
			body.SetStandAloneProfileId(id)

			_, res, err := apiClient.Client.SecurityGroupAPI.SecuritygroupCreate(ctx).CreateSecurityGroupCommand(body).Execute()
			if err != nil {
				return diag.FromErr(tk.CreateError(res, err))
			}
		}
	}

	return readAfterUpdateWithRetries(generateResourceTaikunStandaloneProfileReadWithRetries(), ctx, d, meta)
}

func resourceTaikunStandaloneProfileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	id, err := atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	body := tkcore.DeleteStandAloneProfileCommand{}
	body.SetId(id)

	res, err := apiClient.Client.StandaloneProfileAPI.StandaloneprofileDelete(ctx).DeleteStandAloneProfileCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId("")
	return nil
}

func flattenTaikunStandaloneProfile(rawStandaloneProfile *tkcore.StandAloneProfilesListDto, sg []tkcore.SecurityGroupListDto) map[string]interface{} {

	securityGroups := make([]map[string]interface{}, len(sg))
	for i, rawSecurityGroup := range sg {
		securityGroups[i] = map[string]interface{}{
			"id":          i32toa(rawSecurityGroup.GetId()),
			"name":        rawSecurityGroup.GetName(),
			"cidr":        rawSecurityGroup.GetRemoteIpPrefix(),
			"ip_protocol": strings.ToUpper(rawSecurityGroup.GetProtocol()),
		}
		if rawSecurityGroup.GetPortMinRange() != -1 {
			securityGroups[i]["from_port"] = rawSecurityGroup.GetPortMinRange()
		}
		if rawSecurityGroup.GetPortMaxRange() != -1 {
			securityGroups[i]["to_port"] = rawSecurityGroup.GetPortMaxRange()
		}
	}

	return map[string]interface{}{
		"id":                i32toa(rawStandaloneProfile.GetId()),
		"lock":              rawStandaloneProfile.GetIsLocked(),
		"name":              rawStandaloneProfile.GetName(),
		"organization_id":   i32toa(rawStandaloneProfile.GetOrganizationId()),
		"organization_name": rawStandaloneProfile.GetOrganizationName(),
		"public_key":        rawStandaloneProfile.GetPublicKey(),
		"security_group":    securityGroups,
	}
}

func resourceTaikunStandaloneProfileLock(id int32, lock bool, apiClient *tk.Client) error {
	body := tkcore.StandAloneProfileLockManagementCommand{}
	body.SetId(id)
	body.SetMode(getLockMode(lock))

	res, err := apiClient.Client.StandaloneProfileAPI.StandaloneprofileLockManagement(context.TODO()).StandAloneProfileLockManagementCommand(body).Execute()
	return tk.CreateError(res, err)
}
