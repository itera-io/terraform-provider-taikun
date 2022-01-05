package taikun

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/itera-io/taikungoclient/client/security_group"
	"github.com/itera-io/taikungoclient/client/stand_alone_profile"
	"github.com/itera-io/taikungoclient/models"
)

func resourceTaikunStandaloneProfileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "The ID of the Standalone profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"lock": {
			Description: "Indicates whether to lock the Standalone profile.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"name": {
			Description:      "The name of the Standalone profile.",
			Type:             schema.TypeString,
			Required:         true,
			ValidateFunc:     validation.StringLenBetween(3, 30),
			ValidateDiagFunc: stringIsLowercase,
		},
		"organization_id": {
			Description:      "The ID of the organization which owns the Standalone profile.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ForceNew:         true,
			ValidateDiagFunc: stringIsInt,
		},
		"organization_name": {
			Description: "The name of the organization which owns the Standalone profile.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"public_key": {
			Description:  "The public key of the Standalone profile..",
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"security_group": {
			Description: "List of Security groups.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"cidr": {
						Description:  "Remote IP Prefix.",
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.IsCIDR,
					},
					"from_port": {
						Description:  "Min Range Port",
						Type:         schema.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(0, 65535),
					},
					"id": {
						Description: "ID of the Security group.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"ip_protocol": {
						Description:  "Protocol",
						Type:         schema.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice([]string{"tcp", "udp", "icmp"}, true),
					},
					"name": {
						Description:      "Name of the Security group.",
						Type:             schema.TypeString,
						Required:         true,
						ValidateFunc:     validation.StringLenBetween(3, 30),
						ValidateDiagFunc: stringIsLowercase,
					},
					"to_port": {
						Description:  "Max Range Port",
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

func resourceTaikunStandaloneProfileCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	body := &models.StandAloneProfileCreateCommand{
		Name:      data.Get("name").(string),
		PublicKey: data.Get("public_key").(string),
	}

	if securityGroups, isSecurityGroupsSet := data.GetOk("security_group"); isSecurityGroupsSet {
		rawSecurityGroupList := securityGroups.([]interface{})
		securityGroupList := make([]*models.StandAloneProfileSecurityGroupDto, len(rawSecurityGroupList))
		for i, e := range rawSecurityGroupList {
			rawSecurityGroup := e.(map[string]interface{})
			securityGroupList[i] = &models.StandAloneProfileSecurityGroupDto{
				Name:           rawSecurityGroup["name"].(string),
				PortMaxRange:   int32(rawSecurityGroup["to_port"].(int)),
				PortMinRange:   int32(rawSecurityGroup["from_port"].(int)),
				Protocol:       getSecurityGroupProtocol(rawSecurityGroup["ip_protocol"].(string)),
				RemoteIPPrefix: rawSecurityGroup["cidr"].(string),
			}
		}
		body.SecurityGroups = securityGroupList
	}

	organizationIDData, organizationIDIsSet := data.GetOk("organization_id")
	if organizationIDIsSet {
		organizationId, err := atoi32(organizationIDData.(string))
		if err != nil {
			return diag.Errorf("organization_id isn't valid: %s", data.Get("organization_id").(string))
		}
		body.OrganizationID = organizationId
	}

	params := stand_alone_profile.NewStandAloneProfileCreateParams().WithV(ApiVersion).WithBody(body)
	createResult, err := apiClient.client.StandAloneProfile.StandAloneProfileCreate(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}
	id, err := atoi32(createResult.Payload.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(createResult.Payload.ID)

	if data.Get("lock").(bool) {
		if err := resourceTaikunStandaloneProfileLock(id, true, apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	return readAfterCreateWithRetries(generateResourceTaikunStandaloneProfileReadWithRetries(), ctx, data, meta)
}
func generateResourceTaikunStandaloneProfileReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunStandaloneProfileRead(true)
}
func generateResourceTaikunStandaloneProfileReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunStandaloneProfileRead(false)
}
func generateResourceTaikunStandaloneProfileRead(withRetries bool) schema.ReadContextFunc {
	return func(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*apiClient)
		id, err := atoi32(data.Id())
		data.SetId("")
		if err != nil {
			return diag.FromErr(err)
		}

		response, err := apiClient.client.StandAloneProfile.StandAloneProfileList(stand_alone_profile.NewStandAloneProfileListParams().WithV(ApiVersion).WithID(&id), apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(response.Payload.Data) != 1 {
			if withRetries {
				data.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return nil
		}

		rawStandaloneProfile := response.GetPayload().Data[0]

		securityGroupResponse, err := apiClient.client.SecurityGroup.SecurityGroupList(security_group.NewSecurityGroupListParams().WithV(ApiVersion).WithStandAloneProfileID(id), apiClient)
		if err != nil {
			if _, ok := err.(*security_group.SecurityGroupListNotFound); ok && withRetries {
				data.SetId(i32toa(id))
				return diag.Errorf(notFoundAfterCreateOrUpdateError)
			}
			return diag.FromErr(err)
		}

		err = setResourceDataFromMap(data, flattenTaikunStandaloneProfile(rawStandaloneProfile, securityGroupResponse.GetPayload()))
		if err != nil {
			return diag.FromErr(err)
		}

		data.SetId(i32toa(id))

		return nil
	}
}

func resourceTaikunStandaloneProfileUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)

	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if data.HasChange("name") {
		body := models.StandAloneProfileUpdateCommand{
			ID:   id,
			Name: data.Get("name").(string),
		}
		params := stand_alone_profile.NewStandAloneProfileEditParams().WithV(ApiVersion).WithBody(&body)
		_, err := apiClient.client.StandAloneProfile.StandAloneProfileEdit(params, apiClient)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if data.HasChange("lock") {
		if err := resourceTaikunStandaloneProfileLock(id, data.Get("lock").(bool), apiClient); err != nil {
			return diag.FromErr(err)
		}
	}

	if data.HasChange("security_group") {
		old, new := data.GetChange("security_group")

		// Delete
		oldSecurityGroupList := old.([]interface{})
		for _, e := range oldSecurityGroupList {
			rawSecurityGroup := e.(map[string]interface{})
			id, err := atoi32(rawSecurityGroup["id"].(string))
			if err != nil {
				return diag.FromErr(err)
			}
			body := &models.DeleteSecurityGroupCommand{ID: id}
			params := security_group.NewSecurityGroupDeleteParams().WithV(ApiVersion).WithBody(body)
			_, err = apiClient.client.SecurityGroup.SecurityGroupDelete(params, apiClient)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		// Add
		newSecurityGroupList := new.([]interface{})
		for _, e := range newSecurityGroupList {
			rawSecurityGroup := e.(map[string]interface{})
			body := &models.CreateSecurityGroupCommand{
				Name:                rawSecurityGroup["name"].(string),
				PortMaxRange:        int32(rawSecurityGroup["to_port"].(int)),
				PortMinRange:        int32(rawSecurityGroup["from_port"].(int)),
				Protocol:            getSecurityGroupProtocol(rawSecurityGroup["ip_protocol"].(string)),
				RemoteIPPrefix:      rawSecurityGroup["cidr"].(string),
				StandAloneProfileID: id,
			}
			params := security_group.NewSecurityGroupCreateParams().WithV(ApiVersion).WithBody(body)
			_, err = apiClient.client.SecurityGroup.SecurityGroupCreate(params, apiClient)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return readAfterUpdateWithRetries(generateResourceTaikunStandaloneProfileReadWithRetries(), ctx, data, meta)
}

func resourceTaikunStandaloneProfileDelete(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*apiClient)
	id, err := atoi32(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	body := &models.DeleteStandAloneProfileCommand{
		ID: id,
	}
	params := stand_alone_profile.NewStandAloneProfileDeleteParams().WithV(ApiVersion).WithBody(body)
	_, err = apiClient.client.StandAloneProfile.StandAloneProfileDelete(params, apiClient)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

func flattenTaikunStandaloneProfile(rawStandaloneProfile *models.StandAloneProfilesListDto, sg []*models.SecurityGroupListDto) map[string]interface{} {

	securityGroups := make([]map[string]interface{}, len(sg))
	for i, rawSecurityGroup := range sg {
		securityGroups[i] = map[string]interface{}{
			"id":          i32toa(rawSecurityGroup.ID),
			"name":        rawSecurityGroup.Name,
			"cidr":        rawSecurityGroup.RemoteIPPrefix,
			"ip_protocol": rawSecurityGroup.Protocol,
		}
		if rawSecurityGroup.PortMinRange != -1 {
			securityGroups[i]["from_port"] = rawSecurityGroup.PortMinRange
		}
		if rawSecurityGroup.PortMaxRange != -1 {
			securityGroups[i]["to_port"] = rawSecurityGroup.PortMaxRange
		}
	}

	return map[string]interface{}{
		"id":                i32toa(rawStandaloneProfile.ID),
		"lock":              rawStandaloneProfile.IsLocked,
		"name":              rawStandaloneProfile.Name,
		"organization_id":   i32toa(rawStandaloneProfile.OrganizationID),
		"organization_name": rawStandaloneProfile.OrganizationName,
		"public_key":        rawStandaloneProfile.PublicKey,
		"security_group":    securityGroups,
	}
}

func resourceTaikunStandaloneProfileLock(id int32, lock bool, apiClient *apiClient) error {
	body := models.StandAloneProfileLockManagementCommand{
		ID:   id,
		Mode: getLockMode(lock),
	}
	params := stand_alone_profile.NewStandAloneProfileLockManagementParams().WithV(ApiVersion).WithBody(&body)
	_, err := apiClient.client.StandAloneProfile.StandAloneProfileLockManagement(params, apiClient)
	return err
}
