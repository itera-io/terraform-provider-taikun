package group

import (
	"context"

	tk "github.com/itera-io/taikungoclient"
	tkcore "github.com/itera-io/taikungoclient/client"
	"github.com/itera-io/terraform-provider-taikun/taikun/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTaikunGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description: "Group's name.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"account_id": {
			Description: "Account ID the group belongs to.",
			Type:        schema.TypeInt,
			Required:    true,
			ForceNew:    true,
		},
		"claim_value": {
			Description: "Claim value for the group.",
			Type:        schema.TypeString,
			Optional:    true,
		},
	}
}

func ResourceTaikunGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Taikun Group",
		CreateContext: resourceTaikunGroupCreate,
		ReadContext:   generateResourceTaikunGroupReadWithoutRetries(),
		UpdateContext: resourceTaikunGroupUpdate,
		DeleteContext: resourceTaikunGroupDelete,
		Schema:        resourceTaikunGroupSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceTaikunGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	body := tkcore.CreateGroupCommand{}
	body.SetName(d.Get("name").(string))
	body.SetAccountId(int32(d.Get("account_id").(int)))
	if v, ok := d.GetOk("claim_value"); ok {
		body.SetClaimValue(v.(string))
	}

	id, res, err := apiClient.Client.GroupsAPI.GroupsCreate(ctx).CreateGroupCommand(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId(utils.I32toa(id))

	return utils.ReadAfterCreateWithRetries(generateResourceTaikunGroupReadWithRetries(), ctx, d, meta)
}

func generateResourceTaikunGroupReadWithRetries() schema.ReadContextFunc {
	return generateResourceTaikunGroupRead(true)
}

func generateResourceTaikunGroupReadWithoutRetries() schema.ReadContextFunc {
	return generateResourceTaikunGroupRead(false)
}

func generateResourceTaikunGroupRead(withRetries bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		apiClient := meta.(*tk.Client)
		id, err := utils.Atoi32(d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		accountId := int32(d.Get("account_id").(int))

		response, res, err := apiClient.Client.GroupsAPI.GroupsList(ctx).AccountId(accountId).Execute()
		if err != nil {
			return diag.FromErr(tk.CreateError(res, err))
		}

		var found *tkcore.GroupListItem
		for i := range response.GetData() {
			item := response.GetData()[i]
			if item.GetId() == id {
				found = &item
				break
			}
		}

		if found == nil {
			if withRetries {
				return diag.Errorf(utils.NotFoundAfterCreateOrUpdateError)
			}
			d.SetId("")
			return nil
		}

		err = utils.SetResourceDataFromMap(d, flattenTaikunGroup(found, accountId))
		if err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
}

func resourceTaikunGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)

	id, err := utils.Atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	body := tkcore.UpdateGroupDto{}
	body.SetName(d.Get("name").(string))
	if v, ok := d.GetOk("claim_value"); ok {
		body.SetClaimValue(v.(string))
	} else {
		body.SetClaimValueNil()
	}

	res, err := apiClient.Client.GroupsAPI.GroupsUpdate(ctx, id).UpdateGroupDto(body).Execute()
	if err != nil {
		return diag.FromErr(tk.CreateError(res, err))
	}

	return utils.ReadAfterUpdateWithRetries(generateResourceTaikunGroupReadWithRetries(), ctx, d, meta)
}

func resourceTaikunGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	apiClient := meta.(*tk.Client)
	id, err := utils.Atoi32(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := apiClient.Client.GroupsAPI.GroupsDelete(ctx, id).Execute()
	if err != nil {
		if res != nil && res.StatusCode == 404 {
			return nil
		}
		return diag.FromErr(tk.CreateError(res, err))
	}

	d.SetId("")
	return nil
}

func flattenTaikunGroup(group *tkcore.GroupListItem, accountId int32) map[string]interface{} {
	return map[string]interface{}{
		"name":        group.GetName(),
		"account_id":  accountId,
		"claim_value": group.GetClaimValue(),
	}
}
